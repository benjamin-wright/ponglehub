use std::pin::Pin;

use futures::{Stream, StreamExt};
use serde::{ Serialize, Deserialize };
use kube::{Api, Client as KubeClient, Config, CustomResource};
use kube_runtime::{ watcher, watcher::Event, watcher::Error };

#[derive(CustomResource, Serialize, Deserialize, Default, Debug, Clone)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "Client", namespaced)]
pub struct ClientSpec {
    name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String,
}

pub async fn get_client_events() -> anyhow::Result<Pin<Box<dyn Stream<Item = Result<Event<Client>, Error>> + Send>>> {
    log::debug!("Getting client events...");

    log::trace!("Getting kube config...");
    let config = Config::from_cluster_env()?;

    log::trace!("Getting client...");
    let client = KubeClient::new(config);

    log::trace!("Getting namespaced API...");
    let api: Api<Client> = Api::namespaced(client, "ponglehub");

    log::trace!("Starting watcher...");
    let watcher = watcher(api, kube::api::ListParams::default()).boxed();

    Ok(watcher)
}