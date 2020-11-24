use deadpool_postgres::Pool;
use serde::{ Serialize, Deserialize };
use uuid::Uuid;
use actix_web::{ web, get, post, put, delete, HttpResponse };

pub fn get_routes() -> actix_web::Scope {
    web::scope("/users")
        .service(get_users)
        .service(get_user)
        .service(post_user)
        .service(put_user)
        .service(delete_user)
}

#[derive(Serialize, Deserialize)]
pub struct User {
    id: Uuid,
    name: String,
    email: String
}

#[get("")]
pub async fn get_users(pool: web::Data<Pool>) -> HttpResponse {
    let client = get_client!(pool);
    let user_rows = match client.query("SELECT * FROM users", &[]).await {
        Ok(user_rows) => user_rows,
        Err(e) => {
            log::error!("Failed to get user rows: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    let mut usernames: Vec<User> = vec!();
    for row in user_rows.iter() {
        usernames.push(User{
            id: row.get("id"),
            name: row.get("name"),
            email: row.get("email")
        });
    }

    HttpResponse::Ok().json(usernames)
}

#[get("/{name}")]
pub async fn get_user(pool: web::Data<Pool>, web::Path(name): web::Path<String>) -> HttpResponse {
    let client = get_client!(pool);
    let user_rows = match client.query("SELECT * FROM users WHERE name = $1", &[ &name ]).await {
        Ok(user_rows) => user_rows,
        Err(e) => {
            log::error!("Failed to fetch user {}: {:?}", name, e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    if user_rows.len() != 1 {
        log::error!("Error: Expected 1 user, got {}", user_rows.len());
        return HttpResponse::NotFound().finish();
    }

    let row = match user_rows.get(0) {
        Some(row) => row,
        None => {
            log::error!("Failed to get user from user_rows: {}", name);
            return HttpResponse::InternalServerError().finish();
        }
    };

    HttpResponse::Ok().json(User{
        id: row.get("id"),
        name: row.get("name"),
        email: row.get("email")
    })
}

#[derive(Serialize, Deserialize)]
pub struct UserSeed {
    name: String,
    email: String
}

#[post("")]
pub async fn post_user(pool: web::Data<Pool>, user: web::Json<UserSeed>) -> HttpResponse {
    let client = get_client!(pool);
    if let Err(e) = client.query("INSERT INTO users (name, email, verified) VALUES ($1, $2, false)", &[ &user.name, &user.email ]).await {
        log::error!("Failed to add new user: {:?}", e);

        if e.code() == Some(&tokio_postgres::error::SqlState::UNIQUE_VIOLATION) {
            return HttpResponse::Conflict().finish();
        }

        return HttpResponse::InternalServerError().finish();
    }

    HttpResponse::Ok().finish()
}

#[derive(Deserialize, Debug)]
pub struct PutData {
    email: String
}

#[put("/{name}")]
pub async fn put_user(pool: web::Data<Pool>, body: web::Json<PutData>, web::Path(name): web::Path<String>) -> HttpResponse {
    let client = get_client!(pool);

    if let Err(err) = client.query("UPDATE USERS SET email = $2, verified = false WHERE name = $1", &[ &name, &body.email ]).await {
        log::error!("Failed to update client: {:?}", err);
        return HttpResponse::InternalServerError().finish();
    }

    HttpResponse::Ok().finish()
}

#[delete("/{name}")]
pub async fn delete_user(pool: web::Data<Pool>, web::Path(name): web::Path<String>) -> HttpResponse {
    let client = get_client!(pool);
    match client.execute("DELETE FROM users WHERE name = $1", &[ &name ]).await {
        Err(err) => {
            log::error!("Failed to delete user: {:?}", err);
            HttpResponse::InternalServerError().finish()
        }
        Ok(modified) => {
            if modified < 1 {
                log::error!("Failed to delete user: 0 rows affected");
                return HttpResponse::NotFound().finish();
            }

            if modified > 1 {
                log::error!("Failed to delete user: {} rows affected", modified);
                return HttpResponse::InternalServerError().finish();
            }

            HttpResponse::Ok().finish()
        }
    }
}