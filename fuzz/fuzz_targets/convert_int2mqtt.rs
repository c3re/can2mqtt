#![no_main]

use can2mqtt_rs::converter::types::*;
use can2mqtt_rs::converter::int2ascii::IntConverter;
use libfuzzer_sys::fuzz_target;

fuzz_target!(|data: &[u8]| {
    if data.len() < 10 {
        return;
    }
    let bits = data[0].into();
    let instances = data[1].into();
    let bytes = (bits / 8) * instances;
    if data.len() != bytes {return;}
    match IntConverter::new(bits, instances) {
        Ok(cv) => {
            let cf = CANFrame::try_new(data).unwrap();
            let _ = cv.towards_mqtt(cf);
        }
        Err(_) => {
            return;
        }
    }
});
