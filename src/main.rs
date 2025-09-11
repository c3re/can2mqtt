// TODO the go version detects if the can interface is down (this only checks whether interface is available)
// TODO make log messages similar to the one of the go version
// TODO the go version detects if the mqtt server is down / unreachable
// TODO the go version detects whether the config exists...
// TODO wrap asyncs in their own tasks to make all parallel
// man lots of stuff todo...
use can_socket::{CanFrame, tokio::CanSocket};
use can2mqtt::types::C2MFlags;
use can2mqtt::{config::ToCanMap, types::MQTTMngEvent};
use can2mqtt::{config::ToMqttMap, types::CANMngEvent};
use ctflag::Flags;
use inotify::{Inotify, WatchMask};
use log::*;
use rumqttc::{AsyncClient, ConnectReturnCode, Event, EventLoop, MqttOptions, Packet, Publish};
use tokio::task::JoinHandle;
use std::path::{self, PathBuf};
use std::time::Duration;
use tokio::sync::mpsc::{self, Receiver, Sender};
use tokio_stream::StreamExt;
use url::Url;

#[tokio::main]
async fn main() {
    console_subscriber::init();
    let cs: CanSocket;
    let client: AsyncClient;
    let eventloop: EventLoop;
    let flags = get_flags();
    start_logging(&flags);

    info!("Starting can2mqtt version=3.0.0");
    debug!("Config: {flags:?}");

    // --- MQTT ---
    match get_mqtt_connection(flags.mqtt_connection.clone()) {
        Err(e) => {
            error!("Error initializing MQTT-Connection: {e}");
            std::process::exit(1);
        }
        Ok((cl, el)) => {
            client = cl;
            eventloop = el;
        }
    }

    // --- CAN ---
    match get_can_socket(flags.can_interface.clone()) {
        Err(e) => {
            error!("Error initializing CAN-Interface: {e}");
            std::process::exit(1);
        }
        Ok(c) => cs = c,
    }

    // --- CONFIG ---
    let path = "example.csv"; // received argument
    let abs_path = path::absolute(path).unwrap(); // to read the file

    run(cs, client, eventloop, abs_path).await;
}

/// This is the core of the program. We have three main event sources:
/// * New CAN Frames
/// * New MQTT Messages
/// * A change in our config file
///
/// These events are monitored by the async functions `can_rx`, `mqtt_rx` and ìnotify`
async fn run(cs: CanSocket, client: AsyncClient, eventloop: EventLoop, abs_path: PathBuf) {
    // Channels
    let (tx_can, rx_can) = mpsc::channel::<CANMngEvent>(2);
    let (tx_mqtt, rx_mqtt) = mpsc::channel::<MQTTMngEvent>(2);
    let (tx_inotify, rx_inotify) = mpsc::channel::<()>(1); // unit type signals reload request
    tx_inotify.send(()).await.unwrap(); // initial config load, perhaps there is a more elegant way to do it

    // Futures
    // config
    let inotify = tokio::spawn(inotify_listener(abs_path.clone(), tx_inotify));
    let cfg = tokio::spawn(send_config(rx_inotify, abs_path, tx_mqtt.clone(), tx_can.clone())); // parse config
    // mqtt
    let mqtt_mng = tokio::spawn(mqtt_mng(rx_mqtt, client, tx_can.clone())); // this will be where the mqtt client is
    let mqtt_rx = tokio::spawn(mqtt_rx(tx_mqtt.clone(), eventloop)); // this will be where the eventloop runs
    // can
    let can_mng = can_mng(rx_can, &cs, tx_mqtt);
    let can_rx = can_rx(tx_can, &cs);

    // I have the suspicion that this thing here causes the high CPU load, probably tasks are better
    // according to the docs, this is evaluated concurrently on the same task, so probably I should wrap
    // all of these in their own task in order to get some parallelization?, yep thats exactly what the docs say
    let r = tokio::try_join!(flatten(cfg), flatten(inotify), flatten(mqtt_mng), flatten(mqtt_rx), can_mng, can_rx);
    match r {
        Ok(_) => info!("All tasks finished successfully (should not happen, so not so successful as it sounds...)"),
        Err(s) => error!("{s}")
    };
}

// CONFIG
async fn inotify_listener(abs_path: PathBuf, tx_inotify: Sender<()>) -> Result<(), &'static str> {
    let watch_path = abs_path
        .parent()
        .ok_or("setting up watch: file has no surrounding dir")?; // to watch the dir
    let filename = abs_path
        .file_name()
        .ok_or("no trailing filename in path found")?;
    // poor mans flatten because we deal with an io::Result and not with a core::result::Result
    let inotify = match Inotify::init() {
        Ok(i) => i,
        Err(_) => return Err("error initializing inotify"),
    };

    // Watch for modify and close events.
    inotify
        .watches()
        .add(watch_path, WatchMask::CLOSE_WRITE)
        .expect("Failed to add file watch");

    let buffer = [0; 1024];
    let mut stream = inotify.into_event_stream(buffer).unwrap();
    while let Some(e) = stream.next().await {
        if e.unwrap().name.unwrap() == filename {
            tx_inotify.send(()).await.unwrap();
        }
    }
    Err("left inotify loop. (thats not good)")
}
async fn send_config(
    mut rx_inotify: Receiver<()>,
    abs_path: PathBuf,
    tx_mqtt: Sender<MQTTMngEvent>,
    tx_can: Sender<CANMngEvent>,
) -> Result<(), &'static str> {
    while rx_inotify.recv().await.is_some() {
        if let Ok(c) = can2mqtt::config::parse(abs_path.to_str().unwrap()) {
            let _ = tx_mqtt.send(MQTTMngEvent::Config(Box::new(c.to_can))).await;
            let _ = tx_can.send(CANMngEvent::Config(Box::new(c.to_mqtt))).await;
        }
    }
    Err("left config parsing loop. (thats not good)")
}

// MQTT
// our receiver that deals with all 3 events
// it owns the client, to manage subscriptions and to send messages
async fn mqtt_mng(
    mut rx_mqtt: Receiver<MQTTMngEvent>,
    client: AsyncClient,
    tx_can: Sender<CANMngEvent>,
) -> Result<(), &'static str> {
    let mut config: ToCanMap = ToCanMap::new();

    while let Some(ev) = rx_mqtt.recv().await {
        match ev {
            MQTTMngEvent::Config(new_config) => {
                info!("New config!");
                for topic in config.keys() {
                    match client.unsubscribe(topic).await {
                        Ok(_) => (),
                        Err(e) => error!("Error unsubscribing {topic}: {e}"),
                    }
                }
                for topic in new_config.keys() {
                    match client.subscribe(topic, rumqttc::QoS::AtLeastOnce).await {
                        Ok(_) => (),
                        Err(e) => error!("Error subscribing {topic}: {e}"),
                    }
                }
                config = *new_config;
            }
            MQTTMngEvent::RX(p) => {
                match config.get(&p.topic) {
                    Some(to_can_pair) => match to_can_pair.convertmode.towards_can(p.payload) {
                        Ok(cd) => {
                            tx_can
                                .send(CANMngEvent::TX(CanFrame::new(to_can_pair.id, cd)))
                                .await
                                .unwrap();
                        }
                        Err(e) => warn!(
                            "Error while converting to CAN, Topic: {} convermode: {}: {}",
                            p.topic, to_can_pair.convertmode, e
                        ),
                    },
                    None => { /* should not happen, but if it does nothing needs to be done */ }
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
    Err("left the mqtt mng loop. (thats not good)")
}

async fn mqtt_rx(tx: Sender<MQTTMngEvent>, mut eventloop: EventLoop) -> Result<(), &'static str> {
    loop {
        match eventloop.poll().await {
            Ok(e) => {
                if let Event::Incoming(i_event) = e {
                    match i_event {
                        Packet::ConnAck(ca) => {
                            if ca.code == ConnectReturnCode::Success {
                                info!("MQTT: Connected.")
                            }
                        },
                        Packet::Publish(p) => {
                            // Received message, forward to mqtt_mng
                            let _ = tx.send(MQTTMngEvent::RX(p)).await;
                        }
                        _ => {}
                    }
                }
            }
            Err(e) => {
                warn!("MQTT: {e}, retrying...");
                // try to reconnect, if fatal give up (and kill the whole program that way)
                use tokio::time;
                time::sleep(Duration::from_secs(3)).await;
            }
        }
    }
}

// CAN
async fn can_mng(mut rx: Receiver<CANMngEvent>, cs: &CanSocket, tx_mqtt: Sender<MQTTMngEvent>) -> Result<(), &'static str> {
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
                            Err(e) => warn!(
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
    Err("left the can_mng loop. thats not good")
}

async fn can_rx(tx: Sender<CANMngEvent>, cs: &CanSocket) -> Result<(), &'static str> {
    loop {
        let cf = match cs.recv().await {
            Ok(cf) => cf,
            Err(_)  => return Err("issue while waiting for can frame") 
        };
        match tx.send(CANMngEvent::RX(cf)).await {
            Ok(_) => (),
            Err(_) => return Err("issue while sending canframe to can_mng")
        }
    }
}

// Typesystem boilerplate...
async fn flatten<T>(handle: JoinHandle<Result<T, &'static str>>) -> Result<T, &'static str> {
    match handle.await {
        Ok(Ok(result)) => Ok(result),
        Ok(Err(err)) => Err(err),
        Err(_) => Err("handling failed"),
    }
}

fn get_flags() -> C2MFlags {
    match C2MFlags::from_args(std::env::args()) {
        Ok((config, _)) => config,
        Err(_) => {
            println!("{}", C2MFlags::description());
            std::process::exit(1);
        }
    }
}

fn start_logging(flags: &C2MFlags) {
    let loglevel = match flags.verbose_output {
        true => Level::Debug,
        false => Level::Info,
    };
    stderrlog::new()
        .verbosity(loglevel)
        .module(module_path!())
        .init()
        .unwrap();
}

fn get_can_socket(interface: String) -> Result<CanSocket, String> {
    match CanSocket::bind(interface) {
        Err(e) => Err(e.to_string()),
        Ok(cs) => match cs.set_receive_own_messages(false) {
            Err(e) => Err(e.to_string()),
            Ok(_) => Ok(cs),
        },
    }
}
fn get_mqtt_connection(settings: String) -> Result<(AsyncClient, EventLoop), String> {
    let url = match Url::parse(&settings) {
        Err(e) => return Err(e.to_string()),
        Ok(u) => {
            // TODO support TLS in the future too
            if u.scheme() != "tcp" && u.scheme() != "" {
                return Err(format!("invalid scheme: {}", u.scheme()));
            }
            if u.path() != "" {
                return Err("invalid path: {u.path}".to_string());
            }
            if u.cannot_be_a_base() {
                return Err(
                    "URL can not be cannot-be-a-base, enjoy the double negation...".to_string(),
                );
            }
            if u.fragment().is_some() {
                return Err("URL fragment has to be empty".to_string());
            }
            if u.query().is_some() {
                return Err("URL cannot have a query".to_string());
            }
            if u.host().is_none() {
                return Err("URL needs to contain a host".to_string());
            }
            u
        }
    };
    let mut mqttoptions = MqttOptions::new(
        "can2mqtt v3.0.0",
        url.host().unwrap().to_string(),
        url.port().unwrap_or(1883),
    );
    if url.password().is_some() {
        mqttoptions.set_credentials(url.username(), url.password().unwrap());
    }
    Ok(AsyncClient::new(mqttoptions, 10))
}
