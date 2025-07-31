use std::process;
use can2mqtt_rs::config;

fn main() {
    if let Err(err) = config::parse("example.csv") {
        println!("error parsing config example: {}", err);
        process::exit(1);
    }
}
