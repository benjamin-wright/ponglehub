#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
extern crate serde;
extern crate serde_json;

mod resources;
mod auth;

use crate::resources::{ client::get_client_events, user_seed::get_user_seed_events };
use crate::auth::api;

use futures::TryStreamExt;
use kube_runtime::{ watcher::Event };
use log::info;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");

    let join = tokio::spawn(async move {
        let mut watcher = match get_client_events().await {
            Err(e) => {
                panic!("Failed to get client events: {:?}", e);
            },
            Ok(watcher) => watcher
        };

        loop {
            let event = match watcher.try_next().await {
                Err(e) => {
                    panic!("Failed to get watch event: {:?}", e);
                },
                Ok(event_option) => match event_option {
                    None => {
                        panic!("No more client events to get");
                    },
                    Some(event) => event
                }
            };

            match event {
                Event::Applied(client) => {
                    if let Err(e) = assert_client(&client).await {
                        log::error!("Failed to create client: {:?}", e);
                    }
                }
                Event::Deleted(client) => {
                    if let Err(e) = delete_client(&client).await {
                        log::error!("Failed to delete client: {:?}", e);
                    }
                }
                Event::Restarted(clients) => {
                    for client in clients {
                        if let Err(e) = assert_client(&client).await {
                            log::error!("Failed to create client: {:?}", e);
                        }
                    }
                }
            }
        }
    });

    let user_join = tokio::spawn(async move {
        let mut watcher = match get_user_seed_events().await {
            Err(e) => {
                panic!("Failed to get user_seed events: {:?}", e);
            },
            Ok(watcher) => watcher
        };

        loop {
            let event = match watcher.try_next().await {
                Err(e) => {
                    panic!("Failed to get watch event: {:?}", e);
                },
                Ok(event_option) => match event_option {
                    None => {
                        log::error!("No more user_seed events to get");
                        return;
                    },
                    Some(event) => event
                }
            };

            match event {
                Event::Applied(seed) => {
                    if let Err(e) = assert_user_seed(&seed).await {
                        log::error!("Failed to create user seed: {:?}", e);
                    }
                }
                Event::Deleted(seed) => {
                    if let Err(e) = delete_user(&seed).await {
                        log::error!("Failed to delete user: {:?}", e);
                    }
                }
                Event::Restarted(seeds) => {
                    for seed in seeds {
                        if let Err(e) = assert_user_seed(&seed).await {
                            log::error!("Failed to create user seed: {:?}", e);
                        }
                    }
                }
            }
        }
    });

    info!("All threads running!");

    if let Err(e) = user_join.await {
        return Err(anyhow::anyhow!(e));
    }

    if let Err(e) = join.await {
        return Err(anyhow::anyhow!(e));
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

async fn assert_user_seed(seed: &crate::resources::user_seed::UserSeed) -> anyhow::Result<()> {
    let name = match seed.metadata.name.as_ref() {
        Some(name) => name.clone(),
        None => return Err(anyhow::anyhow!("user seed had no name!"))
    };

    info!("User seed created or modified: {:?}", name);

    match api::get_user(name.as_str()).await? {
        Some(body) => update_user(seed, name.as_str(), &body).await,
        None => create_user_seed(name.as_str(), &seed.spec).await
    }
}

async fn create_user_seed(name: &str, seed: &crate::resources::user_seed::UserSeedSpec) -> anyhow::Result<()> {
    info!("Creating user seed!");

    api::post_user_seed(api::UserPayload{
        name: String::from(name),
        email: seed.email.clone()
    }).await?;

    crate::resources::user_seed::set_user_seeded(name).await?;

    info!("User seed created");

    return Ok(())
}

impl api::UserPayload {
    fn same(&self, client: &crate::resources::user_seed::UserSeed) -> bool {
        return self.email == client.spec.email;
    }
}

async fn update_user(seed: &crate::resources::user_seed::UserSeed, name: &str, existing: &api::UserPayload) -> anyhow::Result<()> {
    info!("Updating user seed!");

    if existing.same(seed) {
        info!("User seed details have not changed");
        return Ok(());
    }

    api::put_user_seed(name, api::UserSeedPutPayload{
        email: seed.spec.email.clone()
    }).await?;

    crate::resources::user_seed::set_user_seeded(name).await?;

    info!("Updated user seed");

    Ok(())
}

async fn delete_user(seed: &crate::resources::user_seed::UserSeed) -> anyhow::Result<()> {
    info!("Deleting user!");

    let name = match seed.metadata.name.as_ref() {
        Some(name) => name.clone(),
        None => return Err(anyhow::anyhow!("user seed had no name!"))
    };

    api::delete_user(name.as_str()).await?;

    info!("Deleted user");

    Ok(())
}