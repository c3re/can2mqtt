use std::fmt;

use crate::converter::types::{CANFrame, Converter, MQTTPayload};
use hex;

#[derive(Debug)]
pub struct ByteColor2ColorCodeConverter {}

impl ByteColor2ColorCodeConverter {
    pub fn new() -> ByteColor2ColorCodeConverter {
        ByteColor2ColorCodeConverter {}
    }
}

impl Converter for ByteColor2ColorCodeConverter {
    fn towards_mqtt(self: &Self, cf: CANFrame) -> Result<MQTTPayload, String> {
        if cf.len() != 3 {
            return Err(format!(
                "Input does not contain exactly 3 bytes, got {} instead",
                cf.len()
            ));
        }
        Ok(MQTTPayload::from(format!(
            "#{:02x}{:02x}{:02x}",
            cf[0], cf[1], cf[2]
        )))
    }

    fn towards_can(self: &Self, msg: MQTTPayload) -> Result<CANFrame, String> {
        if !msg.is_ascii() {
            return Err(("input contains non-ASCII characters").to_string());
        }
        let str = String::from(msg.escape_ascii().to_string());
        let str = str.strip_prefix("#").unwrap_or(&str);
        if str.len() != 6 {
            return Err(format!(
                "input does not contain exactly 6 nibbles each represented by one character, got {} instead",
                str.len()
            ));
        }
        match hex::decode(str) {
            Ok(v) => {
                assert_eq!(v.len(), 3);
                Ok(CANFrame::new([v[0], v[1], v[2]]))
            }
            Err(e) => Err(format!("Error while converting: {}", e.to_string())),
        }
    }
}

impl fmt::Display for ByteColor2ColorCodeConverter {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "bytecolor2colorcode")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn towards_can_test() {
        let cv = ByteColor2ColorCodeConverter {};
        let msg = MQTTPayload::from("#00ff00");
        match cv.towards_can(msg) {
            Ok(cf) => {
                assert_eq!(cf.len(), 3);
                assert_eq!(cf[0], 0);
                assert_eq!(cf[1], 255);
                assert_eq!(cf[2], 0);
            }
            Err(e) => {
                panic!("{}", e)
            }
        }
    }

    #[test]
    fn towards_can_test2() {
        let cv = ByteColor2ColorCodeConverter {};
        let msg = MQTTPayload::from("00ff00");
        match cv.towards_can(msg) {
            Ok(cf) => {
                assert_eq!(cf.len(), 3);
                assert_eq!(cf[0], 0);
                assert_eq!(cf[1], 255);
                assert_eq!(cf[2], 0);
            }
            Err(e) => {
                panic!("{}", e)
            }
        }
    }

    #[test]
    fn towards_can_test3() {
        let cv = ByteColor2ColorCodeConverter {};
        let msg = MQTTPayload::from("0ff00");
        cv.towards_can(msg)
            .expect_err("message too short, should fail");
    }

    #[test]
    fn towards_can_test4() {
        let cv = ByteColor2ColorCodeConverter {};
        let msg = MQTTPayload::from("00tf00");
        cv.towards_can(msg)
            .expect_err("message contains non-hex bytes, should fail");
    }

    #[test]
    fn towards_can_test5() {
        let cv = ByteColor2ColorCodeConverter {};
        let msg = MQTTPayload::from("00Ff00");
        match cv.towards_can(msg) {
            Ok(cf) => {
                assert_eq!(cf.len(), 3);
                assert_eq!(cf[0], 0);
                assert_eq!(cf[1], 255);
                assert_eq!(cf[2], 0);
            }
            Err(e) => {
                panic!("{}", e)
            }
        }
    }

    #[test]
    fn towards_mqtt_test() {
        let cv = ByteColor2ColorCodeConverter {};
        let cf = CANFrame::new([0, 255, 0]);
        match cv.towards_mqtt(cf) {
            Ok(msg) => {
                assert_eq!(msg, "#00ff00")
            }
            Err(e) => {
                panic!("{}", e)
            }
        }
    }

    #[test]
    fn towards_mqtt_test2() {
        let cv = ByteColor2ColorCodeConverter {};
        let cf = CANFrame::new([0, 255, 0, 0]);
        cv.towards_mqtt(cf)
            .expect_err("Should err, 4 bytes input are too long");
    }

    #[test]
    fn towards_mqtt_test3() {
        let cv = ByteColor2ColorCodeConverter {};
        let cf = CANFrame::new([0, 255]);
        cv.towards_mqtt(cf)
            .expect_err("Should err, 2 bytes input are too short");
    }

    #[test]
    fn back_and_forth() {
        let cv = ByteColor2ColorCodeConverter {};
        let cf = CANFrame::new([0, 255, 0]);
        let msg = cv.towards_mqtt(cf).unwrap();
        assert_eq!(cf, cv.towards_can(msg).unwrap());
    }

    #[test]
    fn test_name() {
        let cv = ByteColor2ColorCodeConverter {};
        assert_eq!(cv.to_string(), "bytecolor2colorcode");
    }
}
