#![no_main]
use libfuzzer_sys::fuzz_target;
use can2mqtt_rs::converter::none::NoneConverter;
use can2mqtt_rs::converter::types::*;

fuzz_target!(|data: &[u8]| {
    let cv = NoneConverter{};
    let msg = MQTTPayload::copy_from_slice(data);
    let _ = cv.towards_can(msg).expect("panic");
});