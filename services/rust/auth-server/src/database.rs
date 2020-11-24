use deadpool_postgres::Pool;

#[macro_export]
macro_rules! get_client {
    ($pool: expr) => {
        match $pool.get().await {
            Ok(client) => client,
            Err(e) => {
                log::error!("Failed to get connection from pool: {:?}", e);
                return HttpResponse::InternalServerError().finish();
            }
        }
    }
}

pub async fn create_pool() -> std::io::Result<Pool> {
    let mut cfg = deadpool_postgres::Config::default();
    cfg.host = Some("cockroach-cockroachdb-public.infra.svc.cluster.local".to_string());
    cfg.port = Some(26257);
    cfg.dbname = Some("auth".to_string());
    cfg.user = Some("authserver".to_string());

    let pool = cfg.create_pool(tokio_postgres::NoTls).unwrap();
    match pool.get().await {
        Ok(client) => log::info!("Postgres connection available"),
        Err(e) => {
            return Err(std::io::Error::new(std::io::ErrorKind::NotConnected, anyhow::anyhow!("Postgres connection not available: {:?}", e)));
        }
    };

    Ok(pool)
}