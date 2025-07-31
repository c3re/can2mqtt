use std::path::{self, PathBuf};
use inotify::{Inotify, WatchMask};
use tokio;
use tokio_stream::StreamExt;
use tokio::sync::mpsc;
use tokio::sync::mpsc::Sender;


#[tokio::main]
async fn main() {
    let path = "example.csv"; // received argument
    let abs_path = path::absolute(path).unwrap(); // to read the file

    let (tx, mut rx) = mpsc::channel::<()>(10);

    let f1 = monitor_file(abs_path.to_owned(), &tx);
    let f2 = async move {
        while let Some(_) = rx.recv().await {
            print!("file updated ");
            match can2mqtt_rs::config::parse(abs_path.to_str().unwrap()) {
                Ok(_) => println!("Ok"),
                Err(e) => println!("Error: {}", e)
            }
        }
    };
    // initial read
    let _ = tx.send(()).await;
    tokio::join!(f1, f2);
}

async fn monitor_file(abs_path: PathBuf, tx: &Sender<()>) {
    let watch_path = abs_path.parent().unwrap(); // to watch the dir
    let filename = abs_path.file_name().unwrap().to_owned(); // to filter the watch 
    let inotify = Inotify::init().expect("Error while initializing inotify instance");

    // Watch for modify and close events.
    inotify
        .watches()
        .add(
            watch_path,
            WatchMask::CLOSE_WRITE,
        )
        .expect("Failed to add file watch");

    let buffer = [0; 1024];
    let mut stream = inotify.into_event_stream(buffer).unwrap();
    while let Some(e) = stream.next().await {
        if e.unwrap().name.unwrap() == filename {
            let _ = tx.send(()).await;
        }
    }
}