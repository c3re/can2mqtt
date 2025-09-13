// The idea in this example is to explore command line argument parsing
// I am not necessarily going to use clap but I really like that acronym
use ctflag::*;
use log::{info, Level};

#[derive(Flags, Debug)]
struct C2MFlags {
    #[flag(desc = "which config file to use", short='f', placeholder="", default = "can2mqtt.csv")]
    file: String,
    #[flag(desc = "which CAN interface to use", short='c', placeholder="", default = "can0")]
    can_interface: String,
    #[flag(desc = "which mqtt-broker to use. Example: tcp://user:password@broker.hivemq.com:1883", short='m', placeholder="", default = "tcp://localhost:1883")]
    mqtt_connection: String,
    #[flag(desc = "show (very) verbose debug log", short='v', placeholder="", default = false)]
    verbose_output: bool,
    #[flag(desc = "direction mode: 0 - bidirectional, 1 - can2mqtt only, 2 - mqtt2can only", short='d', placeholder="", default = 0)]
    dir_mode: usize,
}

fn main()  {
    match C2MFlags::from_args(std::env::args()) {
        Ok((config, _)) => {
            let loglevel = match config.verbose_output {
                true => Level::Debug,
                false => Level::Info
            };
            stderrlog::new().verbosity(loglevel).module(module_path!()).init().unwrap();
            info!("Starting can2mqtt version=3.0.0 {:?}", config);
        }
        Err(_) => {
            println!("{}", C2MFlags::description());
            std::process::exit(1);
        }
    }
}
