use can_socket::{CanFrame, CanId, tokio::CanSocket};
use can2mqtt_rs::{config::ToCanMap, types::MQTTMngEvent};
use can2mqtt_rs::{config::ToMqttMap, types::CANMngEvent};
use inotify::{Inotify, WatchMask};
use rumqttc::{AsyncClient, Event, EventLoop, MqttOptions, Packet, Publish};
//use rumqttc::Packet::Publish;
use std::{
    path::{self, PathBuf},
    time::Duration,
};
use tokio::sync::mpsc::{self, Receiver, Sender};
use tokio_stream::StreamExt;

use tokio;

#[tokio::main]
async fn main() {
    // --- MQTT ---
    // Client
    let mut mqttoptions = MqttOptions::new("can2mqtt", "localhost", 1883);
    mqttoptions.set_keep_alive(Duration::from_secs(5));
    let (client, eventloop) = AsyncClient::new(mqttoptions, 10);
    // Channel
    let (tx_mqtt, rx_mqtt) = mpsc::channel::<MQTTMngEvent>(10);

    // --- CAN ---
    // Client
    let cs = CanSocket::bind("vcan0").expect("Can't bind CAN socket");
    let _ = cs.set_receive_own_messages(false);
    // Channel
    let (tx_can, rx_can) = mpsc::channel::<CANMngEvent>(10);

    // --- "FUTURES" --- instances of things doing something and communicating via the channels created above
    // config
    let cfg = send_config(&tx_mqtt, &tx_can); // inotify
    // mqtt
    let mqtt_mng = mqtt_mng(rx_mqtt, client, &tx_can); // this will be where the mqtt client is
    let mqtt_rx = mqtt_rx(&tx_mqtt, eventloop); // this will be where the eventloop runs
    // can
    let can_mng = can_mng(rx_can, &cs, &tx_mqtt);
    let can_rx = can_rx(&tx_can, &cs);
    tokio::join!(cfg, mqtt_mng, mqtt_rx, can_mng, can_rx);
}

// CONFIG
async fn send_config(tx_mqtt: &Sender<MQTTMngEvent>, tx_can: &Sender<CANMngEvent>) {
    let path = "example.csv"; // received argument
    let abs_path = path::absolute(path).unwrap(); // to read the file
    let watch_path = abs_path.parent().unwrap(); // to watch the dir
    let filename = abs_path.file_name().unwrap().to_owned(); // to filter the watch 
    let inotify = Inotify::init().expect("Error while initializing inotify instance");

    // Watch for modify and close events.
    inotify
        .watches()
        .add(watch_path, WatchMask::CLOSE_WRITE)
        .expect("Failed to add file watch");

    let buffer = [0; 1024];
    let mut stream = inotify.into_event_stream(buffer).unwrap();
    while let Some(e) = stream.next().await {
        if e.unwrap().name.unwrap() == filename {
            if let Ok(c) = can2mqtt_rs::config::parse(abs_path.to_str().unwrap()) {
                let _ = tx_mqtt.send(MQTTMngEvent::Config(Box::new(c.to_can))).await;
                let _ = tx_can.send(CANMngEvent::Config(Box::new(c.to_mqtt))).await;
            }
        }
    }
}

// MQTT
// our receiver that deals with all 3 events
// it owns the client, to manage subscriptions and to send messages
async fn mqtt_mng(
    mut rx_mqtt: Receiver<MQTTMngEvent>,
    client: AsyncClient,
    tx_can: &Sender<CANMngEvent>,
) {
    let mut config: ToCanMap = ToCanMap::new();

    while let Some(ev) = rx_mqtt.recv().await {
        match ev {
            MQTTMngEvent::Config(new_config) => {
                println!("New config!");
                for topic in config.keys() {
                    match client.unsubscribe(topic).await {
                        Ok(_) => println!("Successfully unsubscribed {}", topic),
                        Err(e) => println!("Error unsubscribing {}: {}", topic, e),
                    }
                }
                for topic in new_config.keys() {
                    match client.subscribe(topic, rumqttc::QoS::AtLeastOnce).await {
                        Ok(_) => println!("Successfully subscribed {}", topic),
                        Err(e) => println!("Error subscribing {}: {}", topic, e),
                    }
                }
                config = *new_config;
            }
            MQTTMngEvent::RX(p) => {
                match config.get(&p.topic) {
                    Some(to_can_pair) => {
                        match to_can_pair.convertmode.towards_can(p.payload) {
                            Ok(cd) => {
                                tx_can.send(CANMngEvent::TX(CanFrame::new(to_can_pair.id, cd))).await.unwrap();
                            }
                            Err(e) => println!("Error while converting to CAN, Topic: {} convermode: {}: {}", p.topic, to_can_pair.convertmode, e).into()
                            
                        }
                    }
                    None => { /* should now happen, but if nothing needs to be done */}
                }
            }
            MQTTMngEvent::TX(p) => {
                // unsubscribe, publish, subscribe, we don't want to receive our own stuff :)
                // TODO find out if it is wise to wait with the awaits, or to chain them or something like that.
                client.unsubscribe(&p.topic).await.unwrap();
                client
                    .publish(&p.topic, p.qos, false, p.payload)
                    .await
                    .unwrap();
                client
                    .subscribe(&p.topic, rumqttc::QoS::AtLeastOnce)
                    .await
                    .unwrap();
            }
        }
    }
}

async fn mqtt_rx(tx: &Sender<MQTTMngEvent>, mut eventloop: EventLoop) {
    loop {
        match eventloop.poll().await {
            Ok(Event::Incoming(Packet::Publish(p))) => {
                let _ = tx.send(MQTTMngEvent::RX(p)).await;
            }
            _ => {}
        }
    }
}

// CAN
async fn can_mng(mut rx: Receiver<CANMngEvent>, cs: &CanSocket, tx_mqtt: &Sender<MQTTMngEvent>) {
    let mut config = ToMqttMap::new();
    while let Some(ev) = rx.recv().await {
        match ev {
            CANMngEvent::Config(new_config) => config = *new_config,
            CANMngEvent::RX(cf) => {
                match config.get(&cf.id()) {
                    Some(to_mqtt_pair) => {
                        match to_mqtt_pair.convertmode.towards_mqtt(cf.data().unwrap()) {
                            Ok(msg) => {
                                let _ = tx_mqtt
                                    .send(MQTTMngEvent::TX(Publish::new(
                                        to_mqtt_pair.topic.clone(),
                                        rumqttc::QoS::AtLeastOnce,
                                        msg,
                                    )))
                                    .await;
                            }
                            Err(e) => println!(
                                "Error converting {:?} to mqtt, convertmode was {}: {}",
                                cf.data().unwrap(),
                                to_mqtt_pair.convertmode,
                                e
                            ),
                        }
                    }
                    None => { /* nothing to do */ }
                }
            }
            CANMngEvent::TX(cf) => {
                cs.send(&cf).await.unwrap();
            }
        }
    }
}

async fn can_rx(tx: &Sender<CANMngEvent>, cs: &CanSocket) {
    loop {
        let cf = cs.recv().await.unwrap();
        tx.send(CANMngEvent::RX(cf)).await.unwrap();
    }
}
