use serde::{ Deserialize };

#[derive(Clone)]
pub struct TokenApi {
    pub gatekeeper_url: String,
}

#[derive(Deserialize)]
pub struct TokenResponse {
    pub token: String,
}

impl TokenApi {
    pub fn new() -> anyhow::Result<TokenApi> {
        let gatekeeper_url = match std::env::var("GATEKEEPER_URL") {
            Ok(url) => url,
            Err(e) => {
                return Err(anyhow::anyhow!("Failed to fetch GATEKEEPER_URL: {:?}", e));
            }
        };

        log::info!("Creating Gatekeeper API with url: {}", gatekeeper_url);

        Ok(TokenApi{
            gatekeeper_url: gatekeeper_url,
        })
    }

    pub async fn get_token(&self) -> anyhow::Result<String> {
        let client = reqwest::Client::new();
        let post_result = client.post(format!("{}/login", self.gatekeeper_url).as_str())
            .send()
            .await;

        if let Err(e) = post_result {
            return Err(anyhow::anyhow!("Error getting token: {:?}", e));
        }

        let response = post_result.unwrap();

        if !response.status().is_success() {
            return Err(anyhow::anyhow!("Gatekeeper returned non-200 code getting new token: {}", response.status()));
        }

        let body: TokenResponse = response.json().await?;

        Ok(body.token)
    }

    pub async fn check_token(&self, token: String) -> anyhow::Result<bool> {
        let url = format!("{}/login/{}", self.gatekeeper_url, token);
        let get_result = reqwest::get(url.as_str())
            .await;

        if let Err(e) = get_result {
            return Err(anyhow::anyhow!("Error checking token: {:?}", e));
        }

        let response = get_result.unwrap();

        if response.status() == reqwest::StatusCode::NOT_FOUND {
            return Ok(false);
        }

        if !response.status().is_success() {
            return Err(anyhow::anyhow!("Gatekeeper returned non-200 code checking token '{}': {}", token, response.status()));
        }

        Ok(true)
    }
}
