#![no_main]

use libfuzzer_sys::fuzz_target;
use can2mqtt_rs::converter::bytecolor2colorcode::ByteColor2ColorCodeConverter;
use can2mqtt_rs::converter::types::*;

fuzz_target!(|data: &[u8]| {
    if data.len() < 9 { // we a are interested in fuzzing the convertmode not the frame creation
        let cv = ByteColor2ColorCodeConverter{};
        let cf = CANFrame::try_new(data).expect("something went wrong");
        match data.len() {
        3 => {
            let _ = cv.towards_mqtt(cf);
        },
        _ => {
            let _ = cv.towards_mqtt(cf).expect_err("Wrong framesize, should err");
        }
    }
    }
});
