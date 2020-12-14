use actix_web::{App, HttpServer, HttpResponse, middleware::Logger, get, put, post};
use actix_files as fs;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    log::info!("Running server...");

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .service(index)
            .service(fs::Files::new("/", "/static").index_file("index.html"))
    })
    .bind("0.0.0.0:80")?
    .run()
    .await
}

#[get("/api/login")]
pub async fn index() -> HttpResponse {
    log::info!("Hit auth endpoint!");
    return HttpResponse::Unauthorized().finish();
}
