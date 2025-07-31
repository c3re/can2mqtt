use std::path::{self, Path};
fn main() {
    let paths = ["/example.csv", "example.csv", "./example.csv", "/home/mamu/vcs/git/can2mqtt-rs/example.csv", "test.csv", "../can2mqtt-rs/example.csv", "/etc/can2mqtt/can2mqtt.csv"];

    for path in paths {
        if let Ok(abs_path) = std::path::absolute(Path::new(path)) {
            // abs_path -> use to read file
            // watch_path -> use to watch
            // filename -> use to filter the watch
            let watch_path = abs_path.parent();
            let filename = abs_path.file_name();

            println!("path: {}\nabsolute_path: {:?}\nwatch_path: {:?}\nfilename: {:?}\n", path, abs_path, watch_path, filename);
        }
    }
}