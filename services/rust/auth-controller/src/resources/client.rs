use serde::{ Serialize, Deserialize };
use kube::{CustomResource};

#[derive(CustomResource, Serialize, Deserialize, Default, Debug, Clone)]
#[kube(group = "auth.ponglehub.co.uk", version = "v1beta1", kind = "Client", namespaced)]
pub struct ClientSpec {
    name: String,
    #[serde(rename = "callbackUrl")]
    callback_url: String,
}
