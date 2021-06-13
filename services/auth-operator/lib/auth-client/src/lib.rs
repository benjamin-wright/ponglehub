use log::{info};
use auth_types::User;

#[derive(Clone)]
pub struct AuthClient {
    url: String
}

impl AuthClient {
    pub fn new(url: String) -> AuthClient {
        AuthClient {
            url
        }
    }

    pub async fn get_users(&self) -> anyhow::Result<Vec<User>> {
        info!("Making request to user list endpoint");

        let result = match reqwest::get(format!("{}/users", self.url.as_str())).await {
            Ok(result) => result,
            Err(err) => {
                return Err(anyhow::anyhow!("Error fetching users: {:?}", err));
            }
        };

        let users: Vec<User> = match result.json().await {
            Ok(users) => users,
            Err(err) => {
                return Err(anyhow::anyhow!("Error parsing response body: {:?}", err));
            }
        };

        Ok(users)
    }
}
