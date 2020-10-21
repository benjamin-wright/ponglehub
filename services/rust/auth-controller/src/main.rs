#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
#[macro_use] extern crate serde;
extern crate serde_json;

use tokio::prelude::*;
use kube::{ api::Api, Client as KubeClient, CustomResource, api::PostParams };

#[tokio::main]
async fn main() -> Result<(), kube::Error> {
    println!("Spinning ups...");

    println!("Getting clients...");
    let client = match ClientApi::new().await {
        Err(e) => {
            println!("Failed to get client: {:?}", e);
            return Err(e);
        }
        Ok(client) => client
    };

    println!("Posting things...");
    let result = client.post(String::from("this"), String::from("thing")).await;

    println!("Match for results...");
    match result {
        Err(e) => println!("Error: {:?}", e),
        Ok(_) => println!("It worked!")
    };

    Ok(())
}

struct ClientApi {
    api: Api<Client>
}

impl ClientApi {
    async fn new() -> Result<ClientApi, kube::Error> {
        let client = KubeClient::try_default().await?;
        let api: Api<Client> = Api::namespaced(client, "ponglehub");

        Ok(ClientApi{
            api: api
        })
    }

    async fn post(&self, name: String, callback_url: String) -> Result<(), kube::Error> {
        let client = serde_json::from_value(serde_json::json!({
            "apiVersion": "auth.ponglehub.co.uk/v1beta1",
            "kind": "Client",
            "metadata": {
                "name": "my-pod"
            },
            "spec": {
                "name": name,
                "callback_url": callback_url
            }
        }))?;

        // Create the client
        self.api.create(&PostParams::default(), &client).await?;

        return Ok(());
    }
}

#[derive(CustomResource, Serialize, Deserialize, Default, Debug, Clone)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "Client", namespaced)]
pub struct ClientSpec {
    name: String,
    callback_url: String,
}
