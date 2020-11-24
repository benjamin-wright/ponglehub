#[macro_use]
mod database;
#[macro_use]
mod kafka;
mod users;
mod clients;
use actix_web::{App, HttpServer, middleware::Logger};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();

    let pool = database::create_pool().await?;
    let producer = kafka::new();

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .data(pool.clone())
            .data(producer.clone())
            .service(clients::get_routes())
            .service(users::get_routes())
    })
    .bind("0.0.0.0:80")?
    .run()
    .await
}
