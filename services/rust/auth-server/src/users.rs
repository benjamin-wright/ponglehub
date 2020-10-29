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