use std::path;
use inotify::{Inotify, WatchMask};
use tokio;
use tokio_stream::StreamExt;


#[tokio::main]
async fn main() {
    let path = "example.csv"; // received argument
    let abs_path = path::absolute(path).unwrap(); // to read the file
    let watch_path = abs_path.parent().unwrap(); // to watch the dir
    let filename = abs_path.file_name().unwrap(); // to filter the watch 
    
    let inotify = Inotify::init().expect("Error while initializing inotify instance");

    // Watch for modify and close events.
    inotify
        .watches()
        .add(
            // we have to monitor the whole path that our config file is contained in
            // Vim for example does not modify a file but creates a new file and moves 
            // it into place. If we just inotify on the strict singular file we will miss 
            // the change
            // CLOSE_WRITE turned out to trigger the right events following some short
            // tests with code and vim
            // further, we have to filter for the correct filename of course
            watch_path,
            WatchMask::CLOSE_WRITE,
        )
        .expect("Failed to add file watch");

    let buffer = [0; 1024];
    let mut stream = inotify.into_event_stream(buffer).unwrap();
    while let Some(e) = stream.next().await {
        if e.unwrap().name.unwrap() == filename {
            println!("Config updated, please reread: {:?}", abs_path);
        }
    }
}
