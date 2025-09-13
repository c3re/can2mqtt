#![no_main]

use libfuzzer_sys::fuzz_target;
use can2mqtt_rs::converter::pixelbin2ascii::PixelBin2Ascii;
use can2mqtt_rs::converter::types::*;

fuzz_target!(|data: &[u8]| {
    if data.len() < 9 { // we a are interested in fuzzing the convertmode not the frame creation
        let cv = PixelBin2Ascii{};
        let cf = CANFrame::try_new(data).expect("something went wrong");
        // so far only checking if we create a panic
        let _ = cv.towards_mqtt(cf);
    }
});
