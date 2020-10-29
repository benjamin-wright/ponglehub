use rocket::http::Status;
use serde::{ Serialize, Deserialize };
use uuid::Uuid;
use crate::database::AuthDB;
use rocket_contrib::json::Json;

#[derive(Serialize, Deserialize)]
pub struct User {
    id: Uuid,
    name: String,
    email: String
}


#[get("/users")]
pub fn get_users(client: AuthDB) -> Json<Vec<User>> {
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

#[derive(Serialize, Deserialize)]
pub struct UserSeed {
    name: String,
    email: String
}

#[post("/users", data = "<user>")]
pub fn post_user(client: AuthDB, user: Json<UserSeed>) -> Result<Status, Status> {
    if let Err(e) = client.0.query("INSERT INTO users (name, email, verified) VALUES ($1, $2, false)", &[ &user.name, &user.email ]) {
        log::error!("Failed to add new user: {:?}", e);
        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}