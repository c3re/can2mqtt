use std::time::Duration;

use can_socket::{tokio::CanSocket, CanFrame, CanId};
use can2mqtt_rs::{config::ToMqttMap, types::CANMngEvent};
use tokio::sync::mpsc::{self, Receiver, Sender};

#[tokio::main]
async fn main() {
    let (tx, rx) = mpsc::channel::<CANMngEvent>(10);

    let cs = CanSocket::bind("vcan0").expect("Can't bind CAN socket");
    let _ = cs.set_receive_own_messages(false);

    let a = can_mng(rx, &cs);
    let b = send_config(&tx);
    let c = can_rx(&tx, &cs);
    let d = can_tx(&tx);
    tokio::join!(a, b, c, d);
}

async fn can_mng(mut rx: Receiver<CANMngEvent>, cs: &CanSocket) {
    let mut config = ToMqttMap::new();
    while let Some(ev) = rx.recv().await {
        match ev {
            CANMngEvent::Config(new_config) => config = *new_config,
            CANMngEvent::RX(cf) => println!("Received Frame {:?}", cf),
            CANMngEvent::TX(cf) => {
                cs.send(&cf).await.unwrap();
            }
        }
    }
}

async fn send_config(tx: &Sender<CANMngEvent>) {

}

async fn can_rx(tx: &Sender<CANMngEvent>, cs: &CanSocket) {
    loop {
        let cf = cs.recv().await.unwrap();
        tx.send(CANMngEvent::RX(cf)).await.unwrap();
    }
}

async fn can_tx(tx: &Sender<CANMngEvent>) {
    loop {
        tokio::time::sleep(Duration::from_millis(3400)).await;
        let ci = CanId::new(15).unwrap();
        let cf = CanFrame::new(ci, [0]);
        tx.send(CANMngEvent::TX(cf)).await.unwrap();
    }
}
