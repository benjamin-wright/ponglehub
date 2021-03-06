#[macro_use]
extern crate serde_json;
extern crate time;

use actix_web::{App, HttpRequest, HttpMessage, HttpResponse, HttpServer, cookie::{self, SameSite}, get, http, middleware::Logger, post, web};
use actix_cors::Cors;
use serde::{ Serialize, Deserialize };
use handlebars::Handlebars;

mod api;

use api::{ TokenApi };
use time::OffsetDateTime;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    log::info!("Running server...");

    let mut handlebars = Handlebars::new();
    handlebars
        .register_templates_directory(".html", "./static")
        .unwrap();
    let handlebars_ref = web::Data::new(handlebars);

    let api = match TokenApi::new() {
        Ok(api) => api,
        Err(e) => {
            panic!("{:?}", e);
        }
    };
    let api_ref = web::Data::new(api);

    HttpServer::new(move || {
        let cors = Cors::default()
            .allowed_origin("http://localhost")
            .allowed_origin_fn(|origin, _req_head| {
                origin.as_bytes().ends_with(b".ponglehub.co.uk")
            })
            .allowed_header(http::header::SET_COOKIE)
            .supports_credentials();

        App::new()
            .wrap(Logger::default())
            .wrap(cors)
            .app_data(api_ref.clone())
            .app_data(handlebars_ref.clone())
            .service(login)
            .service(login_api)
            .service(logout_api)
    })
    .bind("0.0.0.0:80")?
    .run()
    .await
}

#[derive(Serialize, Deserialize)]
pub struct LoginQuery {
    redirect: String
}

#[get("/login")]
pub async fn login(query: web::Query<LoginQuery>, api: web::Data<TokenApi>, hb: web::Data<Handlebars<'_>>) -> HttpResponse {
    log::info!("Getting token from gatekeeper...");
    let token = match api.get_login_token().await {
        Ok(token) => token,
        Err(e) => {
            log::error!("Failed to get token: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    log::info!("Got new token");
    let data = json!({
        "login_token": token,
        "redirect": query.redirect
    });

    let body = hb.render("index", &data).unwrap();

    HttpResponse::Ok().body(body)
}

#[derive(Serialize, Deserialize)]
pub struct LoginData {
    token: String,
    redirect: String,
    username: String,
    password: String
}

#[post("/api/login")]
pub async fn login_api(body: web::Form<LoginData>, api: web::Data<TokenApi>) -> HttpResponse {
    log::info!("Hit auth endpoint -> {}:{} ({} => {})", body.username, body.password, body.token, body.redirect);

    let ok = match api.check_login_token(body.token.clone()).await {
        Ok(ok) => ok,
        Err(e) => {
            log::error!("Failed to check login token: {:?}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    if !ok {
        log::warn!("Login token {} not found", body.token);
        return HttpResponse::Unauthorized().finish();
    }

    let session_token = "special_token";
    let session_cookie = cookie::Cookie::build("pongle_auth", session_token)
        .domain("ponglehub.co.uk")
        .path("/")
        .secure(true)
        .http_only(true)
        .same_site(SameSite::None)
        .finish();

        return HttpResponse::Found()
            .header(http::header::LOCATION, body.redirect.clone())
            .cookie(session_cookie)
            .finish();
}

#[post("/api/logout")]
pub async fn logout_api(request: HttpRequest, api: web::Data<TokenApi>) -> HttpResponse {
    let cookie = match request.cookie("pongle_auth") {
        Some(cookie) => cookie.value().to_string(),
        None => {
            log::info!("Logged out witout being logged in!");
            return HttpResponse::Unauthorized().finish();
        }
    };

    log::info!("Hit logout endpoint -> Cookie: {}", cookie);

    match api.delete_session_token(cookie.clone()).await {
        Ok(_) => log::info!("Session token deleted"),
        Err(e) => log::error!("Failed to delete login token: {:?}", e)
    }

    let session_cookie = cookie::Cookie::build("pongle_auth", "")
        .domain("ponglehub.co.uk")
        .path("/")
        .secure(true)
        .http_only(true)
        .same_site(SameSite::None)
        .expires(OffsetDateTime::now_utc())
        .finish();

    return HttpResponse::Ok()
        .cookie(session_cookie)
        .finish();
}
