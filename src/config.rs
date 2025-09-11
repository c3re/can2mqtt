use std::{collections::HashMap};
use std::sync::Arc;
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
    pub convertmode: Arc<dyn Converter>
}

#[derive(Debug)]
pub struct ToMqttPair {
    pub topic: String,
    pub convertmode: Arc<dyn Converter>
}

pub type ToMqttMap = HashMap<CanId, ToMqttPair>;
pub type ToCanMap = HashMap<String, ToCanPair>;

pub struct ConversionConfig {
    pub to_mqtt: ToMqttMap,
    pub to_can: ToCanMap,
}

pub fn parse(configfile: &str) -> Result<ConversionConfig, Arc<String>> {
    let convertmodes = all::get_convertmodes();
    let mut to_mqtt: ToMqttMap = HashMap::new();
    let mut to_can: ToCanMap = HashMap::new();

    // Build the CSV reader and iterate over each record.
    let csv_rdr_res = csv::ReaderBuilder::new()
    .has_headers(false)
    .quoting(false)
    .from_path(configfile);

    let mut csv_rdr = match csv_rdr_res {
        Ok(r) => r,
        Err(e) => return Err(Arc::new(format!("{e}")))
    };

    for (line, result) in csv_rdr.deserialize().enumerate() {
        let record: CSVTripel = match result {
            Ok(r) => r,
            Err(_) => {CSVTripel{id:0,convertmode: "t".into(), topic: "te".into() }}
        };
        let line = line + 1; // enumerate starts with 0, lines begin with 1

        let canid = match CanId::new(record.id) {
            Ok(c) => c,
            Err(e) => { return Err(Arc::new(format!("Line {}: Invalid CAN ID: {}: {}", line, record.id, e))); }
        };

        if !convertmodes.contains_key(&record.convertmode) {
            return Err(Arc::new(format!("Line {}: Invalid convertmode: {}", line, record.convertmode)));
        }

        let cv = convertmodes[&record.convertmode].clone();

        let to_mqtt_pair = ToMqttPair{topic: record.topic.clone(), convertmode: cv.clone()};
        if to_mqtt.insert(canid, to_mqtt_pair).is_some() {
            return Err(Arc::new(format!("Line {}: CAN ID already exists: {}", line, record.id)));
        }

        let to_can_pair = ToCanPair{id: canid, convertmode: cv};
        if to_can.insert(record.topic.clone(), to_can_pair).is_some() {
            return Err(Arc::new(format!("Line {}: Topic already exists: {}", line, record.topic)));
        }
    }
    Ok(ConversionConfig { to_mqtt, to_can })
}