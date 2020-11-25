use serde::{ Serialize, Deserialize };

#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub struct ClientPayload {
    pub name: String,
    #[serde(rename = "displayName")]
    pub display_name: String,
    #[serde(rename = "callbackUrl")]
    pub callback_url: String
}

pub async fn get_client(name: &str) -> anyhow::Result<Option<ClientPayload>> {
    let client_result = reqwest::get(format!("http://auth-server/api/clients/{}", name).as_str()).await;

    if let Err(e) = client_result {
        return Err(anyhow::anyhow!("Error getting client from auth server: {:?}", e));
    }

    let response = client_result.unwrap();

    if response.status() == reqwest::StatusCode::NOT_FOUND {
        return Ok(None);
    }

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code getting client '{}': {}", name, response.status()));
    }

    let body: ClientPayload = response.json().await?;

    Ok(Some(body))
}

pub async fn post_client(payload: ClientPayload) -> anyhow::Result<()> {
    let client = reqwest::Client::new();
    let post_result = client.post("http://auth-server/api/clients")
        .json(&payload)
        .send()
        .await;

    let response = match post_result {
        Ok(response) => response,
        Err(e) => return Err(anyhow::anyhow!("Error posting client to auth server: {:?}", e))
    };

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code posting client '{:?}': {}", payload, response.status()));
    }

    Ok(())
}

#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub struct ClientPutPayload {
    #[serde(rename = "displayName")]
    pub display_name: String,
    #[serde(rename = "callbackUrl")]
    pub callback_url: String
}

pub async fn put_client(name: &str, payload: ClientPutPayload) -> anyhow::Result<()> {
    let client = reqwest::Client::new();
    let put_result = client.put(format!("http://auth-server/api/clients/{}", name).as_str())
        .json(&payload)
        .send()
        .await;

    let response = match put_result {
        Ok(response) => response,
        Err(e) => return Err(anyhow::anyhow!("Error updating client to auth server: {:?}", e))
    };

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code updating client '{:?}': {}", payload, response.status()));
    }

    Ok(())
}

pub async fn delete_client(name: &str) -> anyhow::Result<()> {
    let client = reqwest::Client::new();
    let delete_result = client.delete(format!("http://auth-server/api/clients/{}", name).as_str())
        .send()
        .await;

    let response = match delete_result {
        Ok(response) => response,
        Err(e) => return Err(anyhow::anyhow!("Error deleting client on auth server: {:?}", e))
    };

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code deleting client {}: {}", name, response.status()));
    }

    Ok(())
}

#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub struct UserPayload {
    pub name: String,
    pub email: String
}

pub async fn get_user(name: &str) -> anyhow::Result<Option<UserPayload>> {
    let user_result = reqwest::get(format!("http://auth-server/api/users/{}", name).as_str()).await;

    if let Err(e) = user_result {
        return Err(anyhow::anyhow!("Error getting user from auth server: {:?}", e));
    }

    let response = user_result.unwrap();

    if response.status() == reqwest::StatusCode::NOT_FOUND {
        return Ok(None);
    }

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code getting user '{}': {}", name, response.status()));
    }

    let body: UserPayload = response.json().await?;

    Ok(Some(body))
}

pub async fn post_user_seed(payload: UserPayload) -> anyhow::Result<()> {
    let client = reqwest::Client::new();
    let post_result = client.post("http://auth-server/api/users")
        .json(&payload)
        .send()
        .await;

    let response = match post_result {
        Ok(response) => response,
        Err(e) => return Err(anyhow::anyhow!("Error posting user seed to auth server: {:?}", e))
    };

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code posting user seed '{:?}': {}", payload, response.status()));
    }

    Ok(())
}

#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub struct UserSeedPutPayload {
    pub email: String
}

pub async fn put_user_seed(name: &str, payload: UserSeedPutPayload) -> anyhow::Result<()> {
    let client = reqwest::Client::new();
    let put_result = client.put(format!("http://auth-server/api/users/{}", name).as_str())
        .json(&payload)
        .send()
        .await;

    let response = match put_result {
        Ok(response) => response,
        Err(e) => return Err(anyhow::anyhow!("Error updating user seed to auth server: {:?}", e))
    };

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code updating user seed '{:?}': {}", payload, response.status()));
    }

    Ok(())
}

pub async fn delete_user(name: &str) -> anyhow::Result<()> {
    let client = reqwest::Client::new();
    let delete_result = client.delete(format!("http://auth-server/api/users/{}", name).as_str())
        .send()
        .await;

    let response = match delete_result {
        Ok(response) => response,
        Err(e) => return Err(anyhow::anyhow!("Error deleting user on auth server: {:?}", e))
    };

    if !response.status().is_success() {
        return Err(anyhow::anyhow!("Auth server returned non-200 code deleting user {}: {}", name, response.status()));
    }

    Ok(())
}