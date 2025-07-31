use can2mqtt_rs::{config::ToCanMap, types::MQTTMngEvent};
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
    let mut mqttoptions = MqttOptions::new("can2mqtt", "localhost", 1883);
    mqttoptions.set_keep_alive(Duration::from_secs(5));

    let (client, eventloop) = AsyncClient::new(mqttoptions, 10);
    // There will be three senders (1 mqtt receive) (1 mqtt send (on behalf of CANMgr)) and 1 inotify
    let (tx, rx) = mpsc::channel::<MQTTMngEvent>(10);
    let a = mqtt_mng(rx, client); // this will be where the mqtt client is
    let b = send_config(&tx); // inotify
    let c = mqtt_tx(&tx); // this will be at some time something done by CAN
    let d = mqtt_rx(&tx, eventloop); // this will be where the eventloop runs
    tokio::join!(a, b, c, d);
}

// our receiver that deals with all 3 events
// it owns the client, to manage subscriptions and to send messages
async fn mqtt_mng(mut rx: Receiver<MQTTMngEvent>, client: AsyncClient) {
    let mut config: ToCanMap = ToCanMap::new();

    while let Some(ev) = rx.recv().await {
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
                // convert and send towards CAN
                println!("RX {:?}", p);
            }
            MQTTMngEvent::TX(p) => {
                // unsubscribe, publish, subscribe, we don't want to receive our own stuff :)
                // TODO find out if it is wise to wait with the awaits, or to chain them or something like that.
                client.unsubscribe(&p.topic).await.unwrap();
                client.publish(&p.topic, p.qos, false, p.payload).await.unwrap();
                client.subscribe(&p.topic, rumqttc::QoS::AtLeastOnce).await.unwrap();
            }
        }
    }
}

async fn mqtt_tx(tx: &Sender<MQTTMngEvent>) {
    loop {
        tokio::time::sleep(Duration::from_millis(1600)).await;
        let _ = tx
            .send(MQTTMngEvent::TX(Publish::new(
                "topic",
                rumqttc::QoS::AtLeastOnce,
                "test",
            )))
            .await;
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

async fn send_config(tx: &Sender<MQTTMngEvent>) {
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
                let _ = tx.send(MQTTMngEvent::Config(Box::new(c.to_can))).await;
            }
        }
    }
}
