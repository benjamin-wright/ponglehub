use log::{info, warn};
use reqwest::Client;

pub struct AuthClient {
    url: String
}

impl AuthClient {
    pub fn new(url: String) -> AuthClient {
        AuthClient {
            url
        }
    }

    pub fn get_users(&self) {
        let client = Client::new();

        let res = client.get(String::from(self.url));
    }
}