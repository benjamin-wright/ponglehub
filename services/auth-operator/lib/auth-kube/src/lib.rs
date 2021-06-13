use std::fmt::Debug;

use reflector::store::Writer;
use log::{info, warn};
use futures::{Future, StreamExt, TryStreamExt};
use serde::{Deserialize, Serialize};
use schemars::JsonSchema;
use kube::{
    api::{Api, ListParams},
    Client, CustomResource,
};
use kube_runtime::{reflector, utils::try_flatten_applied, watcher};
pub use auth_types::{User};

#[derive(CustomResource, Deserialize, Serialize, Clone, Debug, JsonSchema)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "AuthUser", plural="auth-users", singular="auth-user")]
pub struct AuthUserSpec {
    name: String,
    email: String,
    password: String
}

fn from_crd(user: &AuthUser) -> User {
    User{
        name: user.spec.name.to_string(),
        email: user.spec.email.to_string(),
        password: user.spec.password.to_string()
    }
}

pub struct AuthUserWatcher {
    client: Client,
    namespace: String
}

impl AuthUserWatcher {
    pub async fn new(namespace: String) -> anyhow::Result<AuthUserWatcher> {
        info!("Creating new client...");

        let client = Client::try_default().await?;
        info!("Done");

        Ok(AuthUserWatcher{
            client,
            namespace: namespace
        })
    }

    pub async fn start<T, U, Fut>(self, update: T, refresh: U) -> anyhow::Result<()> where
        T: Fn(User),
        U: Fn(Vec<User>)->Fut + Send + 'static,
        Fut: Future<Output = anyhow::Result<()>> + Send {
        let store = Writer::<AuthUser>::default();
        let reader = store.as_reader();
        let users: Api<AuthUser> = Api::namespaced(self.client, &self.namespace);

        let lp = ListParams::default().timeout(60);
        let rf = reflector(store, watcher(users, lp));

        tokio::spawn(async move {
            loop {
                info!("Refreshing state...");
                // Periodically read our state
                let users = reader.state().iter().map(|user| from_crd(user)).collect::<Vec<_>>();
                let users_result = refresh(users);
                match users_result.await {
                    Ok(()) => info!("Done"),
                    Err(err) => warn!("Failed to refresh state: {:?}", err)
                };

                tokio::time::sleep(std::time::Duration::from_secs(60)).await;
            }
        });

        let mut rfa = try_flatten_applied(rf).boxed();
        while let Some(event) = rfa.try_next().await? {
            info!("Received user update for {}", event.spec.name);
            update(from_crd(&event));
            info!("Done");
        }

        warn!("Shouldn't have gotten here, bailed out of infinite kube watch loop");

        Ok(())
    }
}
