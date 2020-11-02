#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;
#[macro_use] extern crate rocket_contrib;

mod database;
mod users;
mod clients;

fn main() {
    env_logger::init();

    rocket::ignite()
        .attach(database::fairing())
        .mount("/", routes![
            users::get_users, users::get_user, users::post_user, users::put_user, users::delete_user,
            clients::get_client, clients::post_client, clients::put_client, clients::delete_client
        ])
        .launch();
}
