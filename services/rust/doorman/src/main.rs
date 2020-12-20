#[macro_use]
extern crate serde_json;

use actix_web::{App, cookie, web, http, HttpServer, HttpResponse, middleware::Logger, get, post};
use actix_cors::Cors;
use serde::{ Serialize, Deserialize };
use handlebars::Handlebars;

mod api;

use api::{ TokenApi };

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
            .allowed_origin("https://auth.ponglehub.co.uk");

        App::new()
            .wrap(Logger::default())
            .wrap(cors)
            .app_data(api_ref.clone())
            .app_data(handlebars_ref.clone())
            .service(login)
            .service(login_api)
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
        .http_only(true)
        .finish();

    return HttpResponse::TemporaryRedirect()
        .header(http::header::LOCATION, body.redirect.clone())
        .cookie(session_cookie)
        .finish();
}
