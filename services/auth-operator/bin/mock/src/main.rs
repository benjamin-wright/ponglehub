use rocket::State;
use rocket::serde::{Serialize, Deserialize, json::{ Json }};
use rocket::tokio::sync::Mutex;
use uuid::Uuid;

#[macro_use] extern crate rocket;

#[derive(Deserialize, Serialize, Clone)]
struct User {
    id: String,
    name: String,
    email: String,
    password: String,
    verified: bool
}

type UserDB = Mutex<Vec<User>>;

type Users<'r> = &'r State<UserDB>;

#[get("/")]
async fn list(users: Users<'_>) -> Json<Vec<User>> {
    let users = users.lock().await;
    let mut result: Vec<User> = vec!();

    for user in users.iter() {
        result.push(user.clone());
    }

    Json(result)
}

#[derive(Serialize)]
struct PostResponse {
    id: String
}

#[post("/", data = "<user>", format = "json")]
async fn post(user: Json<User>, users: Users<'_>) -> Json<PostResponse> {
    let mut users = users.lock().await;

    let mut user = user.to_owned();
    let id = Uuid::new_v4().to_hyphenated().to_string();

    user.id = id.to_string();

    users.push(user);

    Json(PostResponse{ id })
}

#[launch]
fn rocket() -> _ {
    rocket::build()
        .manage(UserDB::new(vec![]))
        .mount("/users", routes![list, post])
}