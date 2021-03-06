use std::pin::Pin;

use futures::{Stream, StreamExt};
use serde::{ Serialize, Deserialize };
use kube::{Api, Client, Config, CustomResource, api::PatchParams};
use kube_runtime::{watcher, watcher::Error, watcher::Event};

#[derive(CustomResource, Serialize, Deserialize, Default, Debug, Clone)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "UserSeed", namespaced)]
pub struct UserSeedSpec {
    pub email: String,
}

#[derive(Deserialize, Serialize, Clone, Debug, Default)]
pub struct UserSeedStatus {
    seeded: bool
}

pub async fn get_user_seed_events() -> anyhow::Result<Pin<Box<dyn Stream<Item = Result<Event<UserSeed>, Error>> + Send>>> {
    log::debug!("Getting client events...");

    log::trace!("Getting kube config...");
    let config = Config::from_cluster_env()?;

    log::trace!("Getting client...");
    let client = Client::new(config);

    log::trace!("Getting namespaced API...");
    let api: Api<UserSeed> = Api::namespaced(client, "ponglehub");

    log::trace!("Starting watcher...");
    let watcher = watcher(api, kube::api::ListParams::default()).boxed();

    Ok(watcher)
}

pub async fn set_user_seeded(name: &str) -> anyhow::Result<()> {
    log::debug!("Getting client events...");

    log::trace!("Getting kube config...");
    let config = Config::from_cluster_env()?;

    log::trace!("Getting client...");
    let client = Client::new(config);

    log::trace!("Getting namespaced API...");
    let api: Api<UserSeed> = Api::namespaced(client, "ponglehub");

    let data = serde_json::json!({
        "status": {
            "seeded": true
        }
    });

    log::info!("Setting seeded");
    api.patch(name, &PatchParams::default(), serde_json::to_vec(&data)?).await?;
    log::info!("Set");
    Ok(())
}