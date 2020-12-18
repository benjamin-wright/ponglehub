use actix_web::{App, HttpResponse, HttpServer, get, post, middleware::Logger, web};
use deadpool_redis::{Config, ConnectionWrapper, Pool, cmd, redis::RedisError};
use serde::{Serialize};
use deadpool::managed::{ Object };

const LOGIN_EXPIRY_SECONDS: &str = "60";

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    log::info!("Running server...");

    HttpServer::new(move || {
        let mut cfg = Config::default();
        cfg.url = Some(String::from("redis://ponglehub-redis-headless:6379"));
        let pool = cfg.create_pool().unwrap();

        App::new()
            .wrap(Logger::default())
            .data(pool)
            .service(index)
            .service(get_login)
            .service(post_login)
    })
    .bind("0.0.0.0:80")?
    .run()
    .await
}

#[get("/")]
pub async fn index() -> HttpResponse {
    log::info!("Hit auth endpoint!");
    return HttpResponse::Unauthorized().finish();
}

#[derive(Serialize)]
pub struct TokenResponse {
    pub token: String
}

async fn get_logged_in_client(redis: web::Data<Pool>) -> anyhow::Result<Object<ConnectionWrapper, RedisError>> {
    let mut r = redis.get().await.unwrap();
    let redis_password = match std::env::var("REDIS_PASSWORD") {
        Ok(password) => password,
        Err(e) => {
            return Err(anyhow::anyhow!("REDIS_PASSWORD environment variable must be defined!"));
        }
    };

    if let Err(e) = cmd("AUTH").arg(&[redis_password]).execute_async(&mut r).await {
        return Err(anyhow::anyhow!("Error getting client: {:?}", e));
    }

    return Ok(r);
}

#[get("/login/{token}")]
pub async fn get_login(redis: web::Data<Pool>, web::Path(token): web::Path<String>) -> HttpResponse {
    log::info!("Submitting login token!");
    let mut r = match get_logged_in_client(redis).await {
        Ok(r) => r,
        Err(e) => {
            log::error!("Failed to get client: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    match cmd("GET").arg(&[token]).query_async::<String>(&mut r).await {
        Ok(result) => {
            log::info!("Worked: {}", result);
            HttpResponse::Ok().finish()
        },
        Err(e) => {
            log::error!("Failed to create: {:?}", e);
            HttpResponse::InternalServerError().finish()
        }
    }
}

#[post("/login")]
pub async fn post_login(redis: web::Data<Pool>) -> HttpResponse {
    log::info!("Creating new login token!");
    let mut r = match get_logged_in_client(redis).await {
        Ok(r) => r,
        Err(e) => {
            panic!(format!("Failed to get client: {:?}", e));
        }
    };

    let token = "test-token";

    if let Err(e) = cmd("SET").arg(&[token, "value", "EX", LOGIN_EXPIRY_SECONDS]).execute_async(&mut r).await {
        log::error!("Failed to create: {:?}", e);
        return HttpResponse::InternalServerError().finish();
    }

    return HttpResponse::Ok().json(TokenResponse{
        token: String::from(token),
    });
}