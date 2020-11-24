use deadpool_postgres::Pool;
use serde::{Serialize, Deserialize};
use crate::kafka::{ Kafka, KafkaMessage };
use uuid::Uuid;
use actix_web::{ web, get, post, put, delete, HttpResponse };

pub fn get_routes() -> actix_web::Scope {
    web::scope("/clients")
        .service(get_client)
        .service(post_client)
        .service(put_client)
        .service(delete_client)
}

#[derive(Serialize, Deserialize)]
pub struct Client {
    id: Uuid,
    name: String,
    #[serde(rename = "displayName")]
    display_name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[get("/{name}")]
pub async fn get_client(pool: web::Data<Pool>, web::Path(name): web::Path<String>) -> HttpResponse {
    let client = get_client!(pool);

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

#[post("")]
pub async fn post_client(body: web::Json<PostData>, pool: web::Data<Pool>, kafka: web::Data<Kafka>) -> HttpResponse {
    let data: PostData = body.into_inner();
    let client = get_client!(pool);

    if let Err(err) = client.query("INSERT INTO clients (name, display_name, callback_url) VALUES ($1, $2, $3)", &[ &data.name, &data.display_name, &data.callback_url ]).await {
        log::error!("Failed to add client: {:?}", err);
        return HttpResponse::InternalServerError().finish();
    }

    send_to_kafka!(kafka, "ponglehub.auth.create-user", format!("Created user: {}", data.name));

    HttpResponse::Ok().finish()
}

#[derive(Deserialize, Debug)]
pub struct PutData {
    #[serde(rename = "displayName")]
    display_name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

#[put("/{name}")]
pub async fn put_client(pool: web::Data<Pool>, body: web::Json<PutData>, web::Path(name): web::Path<String>, kafka: web::Data<Kafka>) -> HttpResponse {
    let client = get_client!(pool);

    if let Err(err) = client.query("UPDATE clients SET display_name = $2, callback_url = $3 WHERE name = $1", &[ &name, &body.display_name, &body.callback_url ]).await {
        log::error!("Failed to update client: {:?}", err);
        return HttpResponse::InternalServerError().finish();
    }

    send_to_kafka!(kafka, "ponglehub.auth.update-user", format!("Updated user: {}", name));

    HttpResponse::Ok().finish()
}

#[delete("/{name}")]
pub async fn delete_client(pool: web::Data<Pool>, web::Path(name): web::Path<String>, kafka: web::Data<Kafka>) -> HttpResponse {
    let client = get_client!(pool);

    if let Err(err) = client.query("DELETE FROM clients WHERE name = $1", &[ &name ]).await {
        log::error!("Failed to delete client: {:?}", err);
        return HttpResponse::InternalServerError().finish();
    }

    send_to_kafka!(kafka, "ponglehub.auth.delete-user", format!("Deleted user: {}", name));

    HttpResponse::Ok().finish()
}