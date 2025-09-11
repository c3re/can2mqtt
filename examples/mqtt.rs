use rumqttc::AsyncClient;
use rumqttc::MqttOptions;
use rumqttc::EventLoop;

#[tokio::main]
async fn main() {
    let mqttoptions = MqttOptions::new("can2mqtt v3.0.0", "127.0.0.1", 1883);
    let (_client, el ) = AsyncClient::new(mqttoptions, 10);
    tokio::join!(handle_mqtt(el));
    println!("hi");
}

async fn handle_mqtt (mut el: EventLoop) {
    while let Ok(notification) = el.poll().await {
        println!("Received = {:?}", notification);
    }
}