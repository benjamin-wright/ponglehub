use tokio::time::sleep;
use std::time::Duration;
use log::{info};
use anyhow::anyhow;
use serde::{Deserialize, Serialize};
use schemars::JsonSchema;
use kube::{
    api::{Api, ListParams, ResourceExt},
    Client, CustomResource,
};

#[derive(CustomResource, Deserialize, Serialize, Clone, Debug, JsonSchema)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "User")]
pub struct UserSpec {
    name: String,
    email: String,
    password: String
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");

    let client = Client::try_default().await?;
    info!("Got client");

    loop {
        info!("hey ho!");
        sleep(Duration::from_secs(10)).await
    }

    Err(anyhow!("them"))
}
