use log::{info, warn};
use auth_kube::{ AuthUserWatcher, User };
use auth_client::{ AuthClient };

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");

    let namespace = std::env::var("NAMESPACE").unwrap_or("default".into());

    let watcher = AuthUserWatcher::new(namespace).await?;

    let update = |user: User| {
        info!("Updated user: {:?}", user);
    };

    let refresh = |users: Vec<User>| async move {
        info!("Refreshed users: {:?}", users);

        let auth_url = match std::env::var("AUTH_ENDPOINT"){
            Ok(url) => url,
            Err(err) => {
                return Err(anyhow::anyhow!("Failed to get AUTH_ENDPOINT: {:?}", err));
            }
        };
        let client = AuthClient::new(auth_url);

        match client.get_users().await {
            Ok(existing_users) => {
                info!("Existing users: {:?}", existing_users);
            },
            Err(err) => {
                warn!("Failed to fetch users: {:?}", err);
            }
        };

        Ok(())
    };

    watcher.start(update, refresh).await?;

    Ok(())
}
