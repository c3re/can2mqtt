use notify::event;
use rumqttc::Packet::Publish;
use rumqttc::{AsyncClient, Client, Event, MqttOptions, QoS};
use std::thread;
use std::time::Duration;
use tokio;

// Things to explore / get done in example mode
// receive mqtt messages 
// send mqtt messages

// here some thoughts about what can / should run concurrently and what not
// events that can happen during runtime:

// - subscription / unsubscription, because of deaf publishing, or because of hot-reload, this is handled with the "client" object
// - incoming message, this is handled with the "connection" object
// - shutdown signal


#[tokio::main]
async fn main() {
    let mut mqttoptions = MqttOptions::new("can2mqtt", "localhost", 1883);
    mqttoptions.set_keep_alive(Duration::from_secs(5));

    let (client, mut eventloop) = AsyncClient::new(mqttoptions, 10);

    // Iterate to poll the eventloop for connection progress
    loop {
        match eventloop.poll().await {
            Ok(Event::Incoming(Publish(p))) => {
                println!("Received packet {:?}", p);
                p.payload;
            }
            _ => {}
        }
    }
    // we never end up here
}
