use crate::converter::types::{CANFrame, Converter, MQTTPayload};
use bitvec::{prelude::*, view::BitView};
use std::fmt;

#[derive(Debug, Default)]
pub struct SixteenBool2Ascii {}

impl Converter for SixteenBool2Ascii {
    fn towards_mqtt(self: &SixteenBool2Ascii, cf: CANFrame) -> Result<MQTTPayload, String> {
        if cf.len() != 2 {
            return Err(format!(
                "Input does not contain exactly 2 bytes, got {} instead",
                cf.len()
            ));
        }
        let bits = cf.view_bits::<Msb0>();
        let res: String = bits
            .iter()
            .by_vals()
            .map(|b| if b { "1" } else { "0" })
            .collect::<Vec<&str>>()
            .join(" ");

        Ok(MQTTPayload::from(res))
    }

    fn towards_can(self: &SixteenBool2Ascii, msg: MQTTPayload) -> Result<CANFrame, String> {
        if !msg.is_ascii() {
            return Err(("input contains non-ASCII characters").to_string());
        }

        if msg.len() != 31 {
            return Err(format!(
                "message has to be 31 bytes long, received {} bytes",
                msg.len()
            )
            .to_string());
        }

        let str = msg.escape_ascii().to_string();

        if !str.split(" ").all(|s| s == "1" || s == "0") {
            return Err("message contains illegal characters".to_string());
        }

        // I used split_whitespace here before, fuzzing found a case that you can
        // add enough whitespace to satisfy the 31 byte length requirement stated
        // above and then end up with an out of bounds error below
        let bits = str
            .split(" ")
            .map(|s| s == "1")
            .collect::<BitVec<u8, Msb0>>();

        let b0 = bits[0..8].load::<u8>();
        let b1 = bits[8..16].load::<u8>();
        Ok(CANFrame::new([b0, b1]))
    }
}

impl fmt::Display for SixteenBool2Ascii {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "sixteenbool2ascii")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_name() {
        let cv = SixteenBool2Ascii {};
        assert_eq!("sixteenbool2ascii", cv.to_string());
    }
    #[test]
    fn towards_can_test() {
        let cv = SixteenBool2Ascii {};
        let msg = MQTTPayload::from("0 0 0 0 1 0 0 0 0 0 0 0 0 1 0 1");
        match cv.towards_can(msg) {
            Ok(cf) => {
                assert_eq!(cf.len(), 2);
                assert_eq!(cf[0], 8);
                assert_eq!(cf[1], 5);
            }
            Err(e) => {
                panic!("{}", e)
            }
        }
    }

    #[test]
    fn towards_can_test2() {
        let cv = SixteenBool2Ascii {};
        let msg = MQTTPayload::from("0 0 0 1 0 0 0 0 0 0 0 0 1 0 1");
        cv.towards_can(msg).expect_err("should err, only 15 bits");
    }

    #[test]
    fn towards_can_test3() {
        let cv = SixteenBool2Ascii {};
        let msg = MQTTPayload::from("1 0 0 0 1 0 a 0 0 0 0 0 0 1 0 1");
        cv.towards_can(msg).expect_err("should err, illegal char");
    }

    #[test]
    fn towards_can_test4() {
        let cv = SixteenBool2Ascii {};
        let msg = MQTTPayload::from("110 0 0 1 0 1 0 0 0 0 0 0 1 0 1");
        cv.towards_can(msg).expect_err("should err, illegal str");
    }

    #[test]
    fn towards_mqtt_test() {
        let cv = SixteenBool2Ascii {};
        let cf = CANFrame::new([0, 255]);
        match cv.towards_mqtt(cf) {
            Ok(msg) => {
                assert_eq!(msg, "0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1");
            }
            Err(e) => {
                panic!("{}", e)
            }
        }
    }

    #[test]
    fn towards_mqtt_test2() {
        let cv = SixteenBool2Ascii {};
        let cf = CANFrame::new([0, 255, 1]);
        cv.towards_mqtt(cf).expect_err("should err, wrong length");
    }

    #[test]
    fn towards_mqtt_test3() {
        let cv = SixteenBool2Ascii {};
        let cf = CANFrame::new([0]);
        cv.towards_mqtt(cf).expect_err("should err, wrong length");
    }
}
