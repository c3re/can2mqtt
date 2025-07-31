use crate::converter::types::{CANFrame, Converter, MQTTPayload};
use hex;
use std::fmt;

#[derive(Debug)]
pub struct PixelBin2Ascii {}

impl PixelBin2Ascii {
    pub fn new() -> PixelBin2Ascii {
        PixelBin2Ascii {}
    }
}

impl Converter for PixelBin2Ascii {
    fn towards_mqtt(self: &Self, cf: CANFrame) -> Result<MQTTPayload, String> {
        if cf.len() != 4 {
            return Err(format!(
                "Input does not contain exactly 4 bytes, got {} instead",
                cf.len()
            ));
        }
        Ok(MQTTPayload::from(format!(
            "{} #{:02x}{:02x}{:02x}",
            cf[0], cf[1], cf[2], cf[3]
        )))
    }

    fn towards_can(self: &Self, msg: MQTTPayload) -> Result<CANFrame, String> {
        // filter out total mismatches
        if !msg.is_ascii() {
            return Err(("input contains non-ASCII characters").to_string());
        }

        // Split and match
        let str = String::from(msg.escape_ascii().to_string());
        let fields: Vec<&str> = str.split(' ').collect();
        if fields.len() != 2 {
            return Err("input split at whitespace does not contain two fields".to_string());
        }

        // returning CANFrame
        let mut cf = CANFrame::new([0; 4]);

        // deal with first part (pixel number)
        match fields[0].parse::<u8>() {
            Ok(i) => cf[0] = i,
            Err(e) => return Err(e.to_string()),
        }

        // deal with second part (color)
        let str = fields[1].strip_prefix("#").unwrap_or(&fields[1]);
        if str.len() != 6 {
            return Err(format!(
                "color input does not contain exactly 6 nibbles each represented by one character, got {} instead",
                str.len()
            ));
        }
        match hex::decode(str) {
            Ok(v) => {
                assert_eq!(v.len(), 3);
                cf[1] = v[0];
                cf[2] = v[1];
                cf[3] = v[2];
            }
            Err(e) => return Err(format!("Error while converting: {}", e.to_string())),
        }
        return Ok(cf);
    }
}

impl fmt::Display for PixelBin2Ascii {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "pixelbin2ascii")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn towards_mqtt1() {
        let cv = PixelBin2Ascii {};
        let cf = CANFrame::new([12, 0, 255, 0]);
        assert_eq!(cv.towards_mqtt(cf).unwrap(), "12 #00ff00");
    }

    #[test]
    fn towards_mqtt2() {
        let cv = PixelBin2Ascii {};
        let cf = CANFrame::new([12, 0, 255, 0, 12]);
        cv.towards_mqtt(cf)
            .expect_err("should panic, too much data");
    }

    #[test]
    fn towards_mqtt3() {
        let cv = PixelBin2Ascii {};
        let cf = CANFrame::new([255, 0, 12]);
        cv.towards_mqtt(cf)
            .expect_err("should panic, too little data");
    }

    #[test]
    fn towards_can1() {
        let cv = PixelBin2Ascii {};
        let msg = MQTTPayload::copy_from_slice("12 #00ff00".as_bytes());
        let cf = cv.towards_can(msg).unwrap();
        assert_eq!(cf[0], 12);
        assert_eq!(cf[1], 0);
        assert_eq!(cf[2], 255);
        assert_eq!(cf[3], 0);
    }

    #[test]
    fn towards_can2() {
        let cv = PixelBin2Ascii {};
        let msg = MQTTPayload::copy_from_slice("12 00Ff20".as_bytes());
        let cf = cv.towards_can(msg).unwrap();
        assert_eq!(cf[0], 12);
        assert_eq!(cf[1], 0);
        assert_eq!(cf[2], 255);
        assert_eq!(cf[3], 32);
    }

    #[test]
    fn towards_can3() {
        let cv = PixelBin2Ascii {};
        let msg = MQTTPayload::copy_from_slice("12  00Ff20".as_bytes());
        cv.towards_can(msg).expect_err("should err, double space");
    }

    #[test]
    fn towards_can4() {
        let cv = PixelBin2Ascii {};
        let msg = MQTTPayload::copy_from_slice("12 00f20".as_bytes());
        cv.towards_can(msg).expect_err("should err, too short");
    }

    #[test]
    fn towards_can5() {
        let cv = PixelBin2Ascii {};
        let msg = MQTTPayload::copy_from_slice("12 00f20000".as_bytes());
        cv.towards_can(msg).expect_err("should err, too long");
    }

    #[test]
    fn towards_can6() {
        let cv = PixelBin2Ascii {};
        let msg = MQTTPayload::copy_from_slice("256 #00ff00".as_bytes());
        cv.towards_can(msg).expect_err("should err, overflow");
    }

    #[test]
    fn towards_can7() {
        let cv = PixelBin2Ascii {};
        let msg = MQTTPayload::copy_from_slice("-1 #00ff00".as_bytes());
        cv.towards_can(msg).expect_err("should err, signed");
    }

    #[test]
    fn test_name() {
        let cv = PixelBin2Ascii {};
        assert_eq!(cv.to_string(), "pixelbin2ascii".to_string());
    }
}
