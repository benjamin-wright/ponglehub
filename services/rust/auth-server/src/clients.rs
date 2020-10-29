use rocket::{http::Status};
use rocket_contrib::json::Json;
use serde::{Serialize, Deserialize};
use uuid::Uuid;
use crate::database::AuthDB;

#[derive(Serialize, Deserialize)]
pub struct Client {
    id: Uuid,
    name: String,
    #[serde(rename = "displayName")]
    display_name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[get("/clients/<name>")]
pub fn get_client(client: AuthDB, name: String) -> Result<Json<Client>, Status> {
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
            display_name: row.get("display_name"),
            callback_url: row.get("callback_url")
        }
    ))
}

#[derive(Deserialize, Debug)]
pub struct PostData {
    name: String,
    #[serde(rename = "displayName")]
    display_name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[post("/clients", data = "<body>")]
pub fn post_client(client: AuthDB, body: Json<PostData>) -> Result<Status, Status> {
    log::info!("Adding new client: {}", body.name);

    if let Err(err) = client.0.query("INSERT INTO clients (name, display_name, callback_url) VALUES ($1, $2, $3)", &[ &body.name, &body.display_name, &body.callback_url ]) {
        log::error!("Failed to add client: {:?}", err);
        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}

#[derive(Deserialize, Debug)]
pub struct PutData {
    #[serde(rename = "displayName")]
    display_name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[put("/clients/<name>", data = "<body>")]
pub fn put_client(client: AuthDB, body: Json<PutData>, name: String) -> Result<Status, Status> {
    log::info!("Updating client: {}", name);

    if let Err(err) = client.0.query("UPDATE clients SET display_name = $2, callback_url = $3 WHERE name = $1", &[ &name, &body.display_name, &body.callback_url ]) {
        log::error!("Failed to update client: {:?}", err);
        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}

#[delete("/clients/<name>")]
pub fn delete_client(client: AuthDB, name: String) -> Result<Status, Status> {
    log::info!("Deleting client: {}", name);

    if let Err(err) = client.0.query("DELETE FROM clients WHERE name = $1", &[ &name ]) {
        log::error!("Failed to delete client: {:?}", err);
        return Err(Status::InternalServerError);
    }

    Ok(Status::Ok)
}