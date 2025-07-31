#![no_main]

use libfuzzer_sys::fuzz_target;
use can2mqtt_rs::converter::pixelbin2ascii::PixelBin2Ascii;
use can2mqtt_rs::converter::types::*;

fuzz_target!(|data: &[u8]| {
    let cv = PixelBin2Ascii{};
    let msg = MQTTPayload::copy_from_slice(data);
    let _ = cv.towards_can(msg);
});