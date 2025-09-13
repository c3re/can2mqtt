use log::*;
fn main() {
    stderrlog::new().verbosity(Level::Info).module(module_path!()).init().unwrap();
    error!("some failure {}, {}", "key", 1); // stuff that will end can2mqtt
    info!("hihi"); // regular output
    warn!("moin"); // non-critical error
    debug!("verbose shit"); // verbose foo
}
