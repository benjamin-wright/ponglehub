#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;

use std::{thread::sleep, time::Duration};

use rocket_contrib::json::Json;
use serde::{Serialize, Deserialize};

mod logger;
use logger::Logger;

#[derive(Serialize, Deserialize)]
struct IndexResp {
    message: &'static str
}

#[get("/")]
fn index() -> Json<IndexResp> {
    sleep(Duration::from_millis(15));

    Json(IndexResp{
        message: "Oh hai!"
    })
}

fn main() {
    rocket::ignite()
        .attach(Logger{})
        .mount("/", routes![index])
        .launch();
}
