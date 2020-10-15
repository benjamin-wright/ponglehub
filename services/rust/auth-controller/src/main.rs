#![feature(proc_macro_hygiene, decl_macro)]

extern crate kube;
#[macro_use] extern crate serde;

use kube::{ api::Api, Client as KubeClient, CustomResource };

#[tokio::main]
async fn main() -> Result<(), kube::Error> {
    println!("Starting...");

    let client = ClientApi::new().await?;

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
}

#[derive(CustomResource, Serialize, Deserialize, Default, Debug, Clone)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", namespaced)]
pub struct ClientSpec {
    name: String,
    callback_url: String,
}
