use reflector::store::Writer;
use log::{info};
use futures::{StreamExt, TryStreamExt};
use serde::{Deserialize, Serialize};
use schemars::JsonSchema;
use kube::{
    api::{Api, ListParams, ResourceExt},
    Client, CustomResource,
};
use kube_runtime::{reflector, utils::try_flatten_applied, watcher};

#[derive(CustomResource, Deserialize, Serialize, Clone, Debug, JsonSchema)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "AuthUser", plural="auth-users", singular="auth-user")]
pub struct AuthUserSpec {
    name: String,
    email: String,
    password: String
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::init();
    info!("Starting...");

    let client = Client::try_default().await?;
    info!("Got client");

    let namespace = std::env::var("NAMESPACE").unwrap_or("default".into());

    let store = Writer::<AuthUser>::default();
    let reader = store.as_reader();
    let users: Api<AuthUser> = Api::namespaced(client, &namespace);

    let lp = ListParams::default().timeout(60);
    let rf = reflector(store, watcher(users, lp));

    tokio::spawn(async move {
        loop {
            // Periodically read our state
            tokio::time::sleep(std::time::Duration::from_secs(60)).await;
            let crds = reader.state().iter().map(ResourceExt::name).collect::<Vec<_>>();
            info!("Current crds: {:?}", crds);
        }
    });

    let mut rfa = try_flatten_applied(rf).boxed();
    while let Some(event) = rfa.try_next().await? {
        info!("Changed {}", event.name());
    }

    Ok(())
}
