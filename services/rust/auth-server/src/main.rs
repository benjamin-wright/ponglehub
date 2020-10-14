#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;
#[macro_use] extern crate rocket_contrib;

use rocket_contrib::databases::postgres;
use rocket_contrib::json::Json;
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
struct User {
    id: String,
    name: String,
    email: String
}

#[get("/users")]
fn get_users(client: AuthDB) -> Json<Vec<User>> {
    let user_rows = client.0.query("SELECT * FROM users", &[]).unwrap();

    let mut usernames: Vec<User> = vec!();
    for row in user_rows.iter() {
        usernames.push(User{
            id: row.get("id"),
            name: row.get("name"),
            email: row.get("email")
        });
    }

    Json(usernames)
}

#[derive(Serialize)]
struct PostResponse {
    status: String
}

#[derive(Deserialize)]
struct PostData {
    name: String,
    email: String
}

#[post("/users", format = "json", data = "<data>")]
fn post_users(_client: AuthDB, data: Json<PostData>) -> Json<PostResponse> {
    Json(PostResponse{
        status: String::from("Alright then")
    })
}

#[database("auth")]
struct AuthDB (postgres::Connection);

fn main() {
    rocket::ignite()
        .attach(AuthDB::fairing())
        .mount("/", routes![get_users, post_users])
        .launch();
}
