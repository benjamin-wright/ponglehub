#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
extern crate serde;
extern crate serde_json;

mod resources;
mod auth;

use crate::resources::client::get_client_events;
use crate::auth::api;

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
                if let Err(e) = assert_client(&e).await {
                    log::error!("Failed to create client: {:?}", e);
                }
            }
            Event::Deleted(e) => {
                info!("Deleted: {:?}", Meta::name(&e));
            }
            Event::Restarted(e) => {
                for r in e {
                    if let Err(e) = assert_client(&r).await {
                        log::error!("Failed to create client: {:?}", e);
                    }
                }
            }
        }
    }

    info!("Finished!");
    Ok(())
}

async fn assert_client(client: &crate::resources::client::Client) -> anyhow::Result<()> {
    let name = client.spec.name.clone();
    info!("Client created or modified: {:?}", name);

    match api::get_client(name.as_str()).await? {
        Some(body) => update_client(client, &body).await,
        None => create_client(client).await
    }
}

async fn create_client(client: &crate::resources::client::Client) -> anyhow::Result<()> {
    info!("Creating client!");

    api::post_client(api::ClientPayload{
        name: client.spec.name.clone(),
        callback_url: client.spec.callback_url.clone()
    }).await?;

    info!("Client created");

    return Ok(())
}

async fn update_client(client: &crate::resources::client::Client, payload: &api::ClientPayload) -> anyhow::Result<()> {
    info!("Updating client!");

    return Ok(())
}