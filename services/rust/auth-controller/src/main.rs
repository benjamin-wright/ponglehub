#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
extern crate serde;
extern crate serde_json;

mod resources;
mod auth;

use crate::resources::client::get_client_events;

use futures::TryStreamExt;
use kube::{ api::{ Meta }, };
use kube_runtime::{ watcher::Event };
use log::info;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");

    let mut watcher = get_client_events().await?;

    while let Some(event) = watcher.try_next().await? {
        match event {
            Event::Applied(e) => {
                if let Err(e) = assert_client(e).await {
                    log::error!("Failed to create client: {:?}", e);
                }
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

async fn assert_client(client: crate::resources::client::Client) -> anyhow::Result<()> {
    let name = Meta::name(&client);
    info!("Client created or modified: {:?}", name);

    match auth::api::get_client(name.as_str()).await? {
        Some(body) => info!("Got result from client: {:?}", body),
        None => info!("Client not found, creating...")
    }

    Ok(())
}