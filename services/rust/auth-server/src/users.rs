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

#[get("/users/<name>")]
pub fn get_user(client: AuthDB, name: String) -> Result<Json<User>, Status> {
    log::info!("Getting user {}", name);
    let user_rows = client.0.query("SELECT * FROM users WHERE name = $1", &[ &name ]).unwrap();

    if user_rows.len() != 1 {
        log::error!("Error: Expected 1 user, got {}", user_rows.len());
        return Err(Status::NotFound);
    }

    let row = user_rows.get(0);

    Ok(Json(
        User{
            id: row.get("id"),
            name: row.get("name"),
            email: row.get("email")
        }
    ))
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

        if e.code() == Some(&postgres::error::UNIQUE_VIOLATION) {
            return Err(Status::Conflict);
        }

        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}

#[derive(Deserialize, Debug)]
pub struct PutData {
    email: String
}

#[put("/users/<name>", data = "<body>")]
pub fn put_user(client: AuthDB, body: Json<PutData>, name: String) -> Result<Status, Status> {
    log::info!("Updating user: {}", name);

    if let Err(err) = client.0.query("UPDATE USERS SET email = $2, verified = false WHERE name = $1", &[ &name, &body.email ]) {
        log::error!("Failed to update client: {:?}", err);
        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}

#[delete("/users/<name>")]
pub fn delete_user(client: AuthDB, name: String) -> Result<Status, Status> {
    log::info!("Deleting user: {}", name);

    let result = client.0.execute("DELETE FROM users WHERE name = $1", &[ &name ]);
    if let Err(err) = result {
        log::error!("Failed to delete user: {:?}", err);
        return Err(Status::InternalServerError);
    }

    let modified = result.unwrap();
    if modified < 1 {
        log::error!("Failed to delete user: 0 rows affected");
        return Err(Status::NotFound);
    }

    if modified > 1 {
        log::error!("Failed to delete user: {} rows affected", modified);
        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}