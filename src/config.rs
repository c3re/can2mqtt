use std::{collections::HashMap, error::Error, rc::Rc};
use can_socket::can_id;
use can_socket::CanId;

use crate::converter::all;
use crate::converter::types::*;

#[derive(serde::Deserialize, Debug, Clone)]
struct CSVTripel {
    id: u32,
    convertmode: String,
    topic: String,
}

#[derive(Debug)]
pub struct ToCanPair {
    pub id: CanId,
    pub convertmode: Rc<dyn Converter>
}

#[derive(Debug)]
pub struct ToMqttPair {
    pub topic: String,
    pub convertmode: Rc<dyn Converter>
}

pub type ToMqttMap = HashMap<CanId, ToMqttPair>;
pub type ToCanMap = HashMap<String, ToCanPair>;

pub struct ConversionConfig {
    pub to_mqtt: ToMqttMap,
    pub to_can: ToCanMap,
}

pub fn parse(configfile: &str) -> Result<ConversionConfig, Box<dyn Error>> {
    let convertmodes = all::get_convertmodes();
    let mut to_mqtt: ToMqttMap = HashMap::new();
    let mut to_can: ToCanMap = HashMap::new();

    // Build the CSV reader and iterate over each record.
    let mut csv_rdr = csv::ReaderBuilder::new()
    .has_headers(false)
    .quoting(false)
    .from_path(configfile)?;

    for (line, result) in csv_rdr.deserialize().enumerate() {
        let record: CSVTripel = result?;
        let line = line + 1; // enumerate starts with 0, lines begin with 1

        let canid = match CanId::new(record.id) {
            Ok(c) => c,
            Err(e) => { return Err(format!("Line {}: Invalid CAN ID: {}: {}", line, record.id, e).into()); }
        };

        if !convertmodes.contains_key(&record.convertmode) {
            return Err(format!("Line {}: Invalid convertmode: {}", line, record.convertmode).into());
        }

        let cv = convertmodes[&record.convertmode].clone();

        let to_mqtt_pair = ToMqttPair{topic: record.topic.clone(), convertmode: cv.clone()};
        if let Some(_) = to_mqtt.insert(canid, to_mqtt_pair) {
            return Err(format!("Line {}: CAN ID already exists: {}", line, record.id).into());
        }

        let to_can_pair = ToCanPair{id: canid, convertmode: cv};
        if let Some(_) = to_can.insert(record.topic.clone(), to_can_pair) {
            return Err(format!("Line {}: Topic already exists: {}", line, record.topic).into());
        }
    }
    Ok(ConversionConfig { to_mqtt, to_can })
}