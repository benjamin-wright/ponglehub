#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
extern crate serde;
extern crate serde_json;

use futures::{TryStreamExt};
use kube::{ api::Api, Client as KubeClient, Config, CustomResource };
use kube_runtime::{ watcher, utils::try_flatten_applied };
use serde::{ Serialize, Deserialize };

use log::{info, error};

#[tokio::main]
async fn main() -> Result<(), kube::Error> {
    env_logger::init();
    info!("Starting...");

    println!("Getting kube config...");
    let config = Config::from_cluster_env()?;

    println!("Getting client...");
    let client = KubeClient::new(config);

    println!("Getting namespaced API...");
    let api: Api<Client> = Api::namespaced(client, "ponglehub");

    println!("Starting watcher...");
    let watcher = watcher(api, kube::api::ListParams::default());
    let result = try_flatten_applied(watcher)
        .try_for_each(|client| async move {
            log::debug!("Client: {}", kube::api::Meta::name(&client));
            Ok(())
        })
        .await;

    match result {
        Ok(_) => {
            info!("Finished!");
            Ok(())
        },
        Err(e) => {
            error!("This: {:?}", e);
            Ok(())
        }
    }
}

#[derive(CustomResource, Serialize, Deserialize, Default, Debug, Clone)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "Client", namespaced)]
pub struct ClientSpec {
    name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String,
}
