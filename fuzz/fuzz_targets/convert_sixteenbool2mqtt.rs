#![no_main]

use libfuzzer_sys::fuzz_target;
use can2mqtt_rs::converter::sixteenbool2ascii::SixteenBool2Ascii;
use can2mqtt_rs::converter::types::*;

fuzz_target!(|data: &[u8]| {
    if data.len() < 9 { // we a are interested in fuzzing the convertmode not the frame creation
        let cv = SixteenBool2Ascii{};
        let cf = CANFrame::try_new(data).expect("something went wrong");
        match data.len() {
        2 => {
            let _ = cv.towards_mqtt(cf);
        },
        _ => {
            let _ = cv.towards_mqtt(cf).expect_err("Wrong framesize, should err");
        }
    }
    }
});
