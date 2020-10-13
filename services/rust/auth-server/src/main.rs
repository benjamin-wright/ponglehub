#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;

use persistence::{Client};

#[get("/")]
fn index() -> String {
    let mut client = Client::new(
        "authserver",
        "auth",
        "infra-cockroachdb-public.infra.svc.cluster.local",
        26257).unwrap();

    let tables = client.get_tables().unwrap();
    format!("Tables: {:?}", tables)
}

fn main() {
    rocket::ignite().mount("/", routes![index]).launch();
}
