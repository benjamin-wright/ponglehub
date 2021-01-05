use actix_web::{App, HttpRequest, HttpResponse, HttpMessage, HttpServer, middleware::Logger, web, get, post, delete};
use actix_cors::Cors;
use deadpool_redis::{Config, ConnectionWrapper, Pool, cmd, redis::RedisError, redis::ErrorKind};
use serde::{Serialize};
use deadpool::managed::{ Object };

const LOGIN_EXPIRY_SECONDS: &str = "60";

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    log::info!("Running server...");

    let redis_url = match std::env::var("REDIS_URL") {
        Ok(url) => url,
        Err(e) => panic!("Error loading REDIS_URL: {:?}", e)
    };

    HttpServer::new(move || {
        let mut cfg = Config::default();
        cfg.url = Some(redis_url.clone());
        let pool = cfg.create_pool().unwrap();

        let cors = Cors::default()
            .allowed_origin("http://localhost")
            .allowed_origin_fn(|origin, _req_head| {
                origin.as_bytes().ends_with(b".ponglehub.co.uk")
            })
            .supports_credentials();

        App::new()
            .wrap(Logger::default())
            .wrap(cors)
            .data(pool)
            .service(index)
            .service(get_login)
            .service(post_login)
            .service(delete_session)
    })
    .bind("0.0.0.0:80")?
    .run()
    .await
}

#[derive(Serialize)]
pub struct TokenResponse {
    pub token: String
}

async fn get_logged_in_client(redis: web::Data<Pool>) -> anyhow::Result<Object<ConnectionWrapper, RedisError>> {
    let mut r = redis.get().await.unwrap();
    let redis_password = match std::env::var("REDIS_PASSWORD") {
        Ok(password) => password,
        Err(_) => {
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
            if ErrorKind::TypeError == e.kind() {
                log::warn!("Token not found");
                HttpResponse::NotFound().finish()
            } else {
                log::error!("Failed to fetch token: {:?}", e);
                HttpResponse::InternalServerError().finish()
            }
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

#[delete("/session/{token}")]
pub async fn delete_session(redis: web::Data<Pool>, web::Path(token): web::Path<String>) -> HttpResponse {
    log::info!("Deleting session token!");
    let mut r = match get_logged_in_client(redis).await {
        Ok(r) => r,
        Err(e) => {
            log::error!("Failed to get client: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    match cmd("DEL").arg(&[token.as_str()]).query_async::<i32>(&mut r).await {
        Ok(result) => {
            log::info!("Worked: {}", result);
            HttpResponse::Ok().finish()
        },
        Err(e) => {
            log::error!("Failed to delete: {:?}", e);
            HttpResponse::InternalServerError().finish()
        }
    }
}

#[get("/loggedIn")]
pub async fn index(req: HttpRequest) -> HttpResponse {
    let cookie = match req.cookie("pongle_auth") {
        Some(cookie) => cookie,
        None => {
            log::info!("Auth cookie missing");
            return HttpResponse::Unauthorized().finish();
        }
    };

    log::info!("Cookie: {:?}", cookie.value());
    return HttpResponse::Ok().finish();
}