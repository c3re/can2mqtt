
use bitvec::{prelude::*, view::BitView};
fn main() {
    let data = [8u8, 5];
    let x = data.view_bits::<Msb0>();
    let res: String = x.iter()
        .by_vals()
        .map(|b| if b {"1"} else {"0"})
        .collect::<Vec<&str>>()
        .join(" ");
    println!("{}", res);
}