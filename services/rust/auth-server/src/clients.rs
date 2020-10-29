use rocket::{http::Status};
use rocket_contrib::json::Json;
use serde::{Serialize, Deserialize};
use uuid::Uuid;
use crate::database::AuthDB;

#[derive(Serialize, Deserialize)]
pub struct Client {
    id: Uuid,
    name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[get("/clients/<name>")]
pub fn get_clients(client: AuthDB, name: String) -> Result<Json<Client>, Status> {
    log::info!("Getting client {}", name);
    let client_rows = client.0.query("SELECT * FROM clients WHERE name = $1", &[ &name ]).unwrap();

    if client_rows.len() != 1 {
        log::error!("Error: Expected 1 client, got {}", client_rows.len());
        return Err(Status::NotFound);
    }

    let row = client_rows.get(0);

    Ok(Json(
        Client{
            id: row.get("id"),
            name: row.get("name"),
            callback_url: row.get("callback_url")
        }
    ))
}

#[derive(Deserialize, Debug)]
pub struct PostData {
    name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[post("/clients", data = "<body>")]
pub fn post_client(client: AuthDB, body: Json<PostData>) -> Result<Status, Status> {
    log::info!("Adding new client: {}", body.name);

    if let Err(err) = client.0.query("INSERT INTO clients (name, callback_url) VALUES ($1, $2)", &[ &body.name, &body.callback_url ]) {
        log::error!("Failed to add client: {:?}", err);
        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}