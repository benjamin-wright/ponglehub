#[macro_use]
extern crate serde_json;

use actix_web::{App, web, HttpServer, HttpResponse, middleware::Logger, get, post};

use serde::{ Serialize, Deserialize };

use handlebars::Handlebars;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    log::info!("Running server...");

    let mut handlebars = Handlebars::new();
    handlebars
        .register_templates_directory(".html", "./static")
        .unwrap();
    let handlebars_ref = web::Data::new(handlebars);

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
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
pub async fn login(query: web::Query<LoginQuery>, hb: web::Data<Handlebars<'_>>) -> HttpResponse {
    let data = json!({
        "login_token": "abcde",
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
pub async fn login_api(body: web::Form<LoginData>) -> HttpResponse {
    log::info!("Hit auth endpoint -> {}:{} ({} => {})", body.username, body.password, body.token, body.redirect);
    return HttpResponse::Ok().body("Logged in!");
}
