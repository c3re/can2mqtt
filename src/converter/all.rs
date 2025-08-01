use crate::converter::bytecolor2colorcode::ByteColor2ColorCodeConverter;
use crate::converter::int2ascii::IntConverter;
use crate::converter::mymode::MyModeConverter;
use crate::converter::none::NoneConverter;
use crate::converter::pixelbin2ascii::PixelBin2Ascii;
use crate::converter::sixteenbool2ascii::SixteenBool2Ascii;
use crate::converter::types::Converter;
use crate::converter::uint2ascii::UintConverter;
use std::collections::HashMap;
use std::rc::Rc;

/// This functions purpose is to group all convertmodes into a HashMap and make them addressable by their name
// Explanation for the Box type in the Signature:
// Converter (Trait) is not possible in a Signature
// dyn Converter (Trait object) is not Sized
// references do not work because they are used after this function where the referenced thingy will be dropped
// so we end up with Box<dyn Converter>>
// no we don't a Box has the Issue that it can only have one owner, but we want to use our Converter in multiple
// places in the HashMap (there can be multiple topics for example that need to convert to u8's)
pub fn get_convertmodes() -> HashMap<String, Rc<dyn Converter>> {
    // create a HashMap with all convertmodes:
    let mut convertmodes: HashMap<String, Rc<dyn Converter>> = HashMap::new();
    let nonecv = Rc::new(NoneConverter::default());
    let bytecolorcv = Rc::new(ByteColor2ColorCodeConverter::default());
    let mymodecv = Rc::new(MyModeConverter::default());
    let pixelcv = Rc::new(PixelBin2Ascii::default());
    let sixteenboolcv = Rc::new(SixteenBool2Ascii::default());

    convertmodes.insert(nonecv.to_string(), nonecv);
    convertmodes.insert(bytecolorcv.to_string(), bytecolorcv);
    convertmodes.insert(mymodecv.to_string(), mymodecv);
    convertmodes.insert(pixelcv.to_string(), pixelcv);
    convertmodes.insert(sixteenboolcv.to_string(),sixteenboolcv);

    for b in [8, 16, 32, 64] {
        // bits
        for i in [1, 2, 4, 8] {
            // instances
            if (b / 8) * i <= 8 {
                let cv = UintConverter::new(i, b).unwrap();
                convertmodes.insert(cv.to_string(), Rc::new(cv));
                let cv = IntConverter::new(i, b).unwrap();
                convertmodes.insert(cv.to_string(), Rc::new(cv));
            }
        }
    }
    convertmodes
}
