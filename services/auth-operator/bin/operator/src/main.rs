use log::{info};
use auth_kube::{ AuthUserWatcher, User };
use auth_client::{ AuthClient };

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");

    let namespace = std::env::var("NAMESPACE").unwrap_or("default".into());

    let watcher = AuthUserWatcher::new(namespace).await?;
    let client = AuthClient::new();


    let update = |user: User| {
        info!("Updated user: {:?}", user);
    };

    let refresh = |users: Vec<User>| {
        info!("Refreshed users: {:?}", users);

        let users = get_users();
    };

    watcher.start(update, refresh).await?;
    let client = AuthClient::new();

    Ok(())
}
