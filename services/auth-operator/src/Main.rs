use tokio::time::sleep;
use std::time::Duration;
use log::{info};
use anyhow::anyhow;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");
    
    loop {
        info!("yofo!");
        sleep(Duration::from_secs(2)).await
    }

    Err(anyhow!("them"))
}
