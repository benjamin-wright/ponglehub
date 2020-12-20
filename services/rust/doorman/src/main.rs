#[macro_use]
extern crate serde_json;

use actix_web::{App, web, HttpServer, HttpResponse, middleware::Logger, get, post};
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

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(api.clone())
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
    let token = match api.get_token().await {
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
pub async fn login_api(body: web::Form<LoginData>, _api: web::Data<TokenApi>) -> HttpResponse {
    log::info!("Hit auth endpoint -> {}:{} ({} => {})", body.username, body.password, body.token, body.redirect);
    return HttpResponse::Ok().body("Logged in!");
}
