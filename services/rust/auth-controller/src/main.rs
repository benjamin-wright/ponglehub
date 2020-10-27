#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
extern crate serde;
extern crate serde_json;

mod resources;

use crate::resources::client::Client;

use futures::{StreamExt, TryStreamExt};
use kube::{ api::{ Api, Meta }, Client as KubeClient, Config };
use kube_runtime::{ watcher, watcher::Event };
use log::info;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");

    println!("Getting kube config...");
    let config = Config::from_cluster_env()?;

    println!("Getting client...");
    let client = KubeClient::new(config);

    println!("Getting namespaced API...");
    let api: Api<Client> = Api::namespaced(client, "ponglehub");

    println!("Starting watcher...");
    let mut watcher = watcher(api, kube::api::ListParams::default()).boxed();
    while let Some(event) = watcher.try_next().await? {
        match event {
            Event::Applied(e) => {
                info!("Applied: {:?}", Meta::name(&e));
            }
            Event::Deleted(e) => {
                info!("Deleted: {:?}", Meta::name(&e));
            }
            Event::Restarted(e) => {
                for r in e {
                    info!("Restarted: {:?}", Meta::name(&r));
                }
            }
        }
    }

    info!("Finished!");
    Ok(())
}
