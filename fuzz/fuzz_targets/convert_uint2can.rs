#![no_main]

use libfuzzer_sys::fuzz_target;
use can2mqtt_rs::converter::types::*;
use can2mqtt_rs::converter::uint2ascii::UintConverter;

fuzz_target!(|data: &[u8]| {
    if data.len() < 10 {
        return;
    }
    let bits = data[0].into();
    let instances = data[1].into();
    let bytes = (bits / 8) * instances;
    if data.len() != bytes {return;}
    match UintConverter::new(bits, instances) {
        Ok(cv) => {
            let msg = MQTTPayload::copy_from_slice(data);
            let _ = cv.towards_can(msg);
        }
        Err(_) => {
            return;
        }
    }
});
