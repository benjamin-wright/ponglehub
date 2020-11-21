mod users;
mod clients;
mod kafka;
use actix_web::{App, HttpServer};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();

    let mut cfg = deadpool_postgres::Config::default();
    cfg.host = Some("cockroach-cockroachdb-public.infra.svc.cluster.local".to_string());
    cfg.port = Some(26257);
    cfg.dbname = Some("auth".to_string());
    cfg.user = Some("authserver".to_string());

    let pool = cfg.create_pool(tokio_postgres::NoTls).unwrap();
    match pool.get().await {
        Ok(client) => log::info!("Postgres connection available"),
        Err(e) => {
            panic!(format!("Postgres connection not available: {:?}", e));
        }
    };

    let producer = kafka::new();

    HttpServer::new(move || {
        App::new()
            .data(pool.clone())
            .data(producer.clone())
            .service(clients::get_client)
            .service(clients::post_client)
    })
    .bind("0.0.0.0:80")?
    .run()
    .await
}
