use serde::{ Serialize, Deserialize };

#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub struct ClientPayload {
    name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String
}

pub async fn get_client(name: &str) -> anyhow::Result<Option<ClientPayload>> {
    let user_result = reqwest::get(format!("http://auth-server/clients/{}", name).as_str()).await;

    if let Err(e) = user_result {
        return Err(anyhow::anyhow!("Error contacting auth server: {:?}", e));
    }

    let response = user_result.unwrap();

    if response.status() == reqwest::StatusCode::NOT_FOUND {
        return Ok(None);
    }

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code getting client '{}': {}", name, response.status()));
    }

    let body: ClientPayload = response.json().await?;

    Ok(Some(body))
}