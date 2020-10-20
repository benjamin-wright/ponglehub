#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
#[macro_use] extern crate serde;
extern crate serde_json;

use kube::{ api::Api, Client as KubeClient, CustomResource, api::PostParams };

#[tokio::main]
async fn main() -> Result<(), kube::Error> {
    println!("Starting...");

    let client = ClientApi::new().await?;
    client.post(String::from("this"), String::from("thing")).await?;

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
