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
                if let Err(e) = delete_client(&e).await {
                    log::error!("Failed to delete client: {:?}", e);
                }
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
    let name = match client.metadata.name.as_ref() {
        Some(name) => name.clone(),
        None => return Err(anyhow::anyhow!("client had no name!"))
    };

    info!("Client created or modified: {:?}", name);

    match api::get_client(name.as_str()).await? {
        Some(body) => update_client(client, name.as_str(), &body).await,
        None => create_client(name.as_str(), &client.spec).await
    }
}

async fn create_client(name: &str, client: &crate::resources::client::ClientSpec) -> anyhow::Result<()> {
    info!("Creating client!");

    api::post_client(api::ClientPayload{
        name: String::from(name),
        display_name: client.display_name.clone(),
        callback_url: client.callback_url.clone()
    }).await?;

    info!("Client created");

    return Ok(())
}

impl api::ClientPayload {
    fn same(&self, client: &crate::resources::client::Client) -> bool {
        return self.display_name == client.spec.display_name && self.callback_url == client.spec.callback_url;
    }
}

async fn update_client(client: &crate::resources::client::Client, name: &str, existing: &api::ClientPayload) -> anyhow::Result<()> {
    info!("Updating client!");

    if existing.same(client) {
        info!("Client details have not changed");
        return Ok(());
    }

    api::put_client(name, api::ClientPutPayload{
        display_name: client.spec.display_name.clone(),
        callback_url: client.spec.callback_url.clone()
    }).await?;

    info!("Updated client");

    Ok(())
}

async fn delete_client(client: &crate::resources::client::Client) -> anyhow::Result<()> {
    info!("Deleting client!");

    let name = match client.metadata.name.as_ref() {
        Some(name) => name.clone(),
        None => return Err(anyhow::anyhow!("client had no name!"))
    };

    api::delete_client(name.as_str()).await?;

    info!("Deleted client");

    Ok(())
}