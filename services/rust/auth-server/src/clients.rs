use deadpool_postgres::Pool;
use serde::{Serialize, Deserialize};
use crate::kafka::Kafka;
use uuid::Uuid;
use actix_web::{ web, get, post, HttpResponse };

#[derive(Serialize, Deserialize)]
pub struct Client {
    id: Uuid,
    name: String,
    #[serde(rename = "displayName")]
    display_name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[get("/clients/{name}")]
pub async fn get_client(pool: web::Data<Pool>, web::Path(name): web::Path<String>) -> HttpResponse {
    log::info!("Getting client {}", name);
    let client = match pool.get().await {
        Ok(client) => client,
        Err(e) => {
            log::error!("Failed to get connection from pool: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    let client_rows = match client.query("SELECT * FROM clients WHERE name = $1", &[ &name ]).await {
        Ok(client_rows) => client_rows,
        Err(e) => {
            log::error!("Failed to get clients: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    if client_rows.len() != 1 {
        log::error!("Error: Expected 1 client, got {}", client_rows.len());
        return HttpResponse::NotFound().finish();
    }

    match client_rows.get(0) {
        Some(row) => HttpResponse::Ok().json(
            Client{
                id: row.get("id"),
                name: row.get("name"),
                display_name: row.get("display_name"),
                callback_url: row.get("callback_url")
            }
        ),
        None => HttpResponse::NotFound().finish()
    }
}

#[derive(Deserialize, Debug)]
pub struct PostData {
    name: String,
    #[serde(rename = "displayName")]
    display_name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[post("/clients")]
pub async fn post_client(body: web::Json<PostData>, pool: web::Data<Pool>, kafka: web::Data<Kafka>) -> HttpResponse {
    let data: PostData = body.into_inner();
    log::info!("Adding new client: {}", data.name);

    let client = match pool.get().await {
        Ok(client) => client,
        Err(e) => {
            log::error!("Failed to get connection from pool: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    if let Err(e) = kafka.send("Hi there!".to_string()).await {
        log::error!("Failed to post to kafka: {:?}", e);
    }

    if let Err(err) = client.query("INSERT INTO clients (name, display_name, callback_url) VALUES ($1, $2, $3)", &[ &data.name, &data.display_name, &data.callback_url ]).await {
        log::error!("Failed to add client: {:?}", err);
        return HttpResponse::InternalServerError().finish();
    }

    HttpResponse::Ok().finish()
}

// #[derive(Deserialize, Debug)]
// pub struct PutData {
//     #[serde(rename = "displayName")]
//     display_name: String,
//     #[serde(rename = "callbackUrl")]
//     callback_url: String
// }

// #[put("/clients/<name>", data = "<body>")]
// pub fn put_client(client: AuthDB, body: Json<PutData>, name: String) -> Result<Status, Status> {
//     log::info!("Updating client: {}", name);

//     if let Err(err) = client.0.query("UPDATE clients SET display_name = $2, callback_url = $3 WHERE name = $1", &[ &name, &body.display_name, &body.callback_url ]) {
//         log::error!("Failed to update client: {:?}", err);
//         return Err(Status::InternalServerError);
//     }

//     Ok(Status::Ok)
// }

// #[delete("/clients/<name>")]
// pub fn delete_client(client: AuthDB, name: String) -> Result<Status, Status> {
//     log::info!("Deleting client: {}", name);

//     if let Err(err) = client.0.query("DELETE FROM clients WHERE name = $1", &[ &name ]) {
//         log::error!("Failed to delete client: {:?}", err);
//         return Err(Status::InternalServerError);
//     }

//     Ok(Status::Ok)
// }