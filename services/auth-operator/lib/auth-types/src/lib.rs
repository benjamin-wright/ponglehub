use serde::{ Serialize, Deserialize };

#[derive(Debug, Deserialize, Serialize)]
pub struct User {
    pub name: String,
    pub email: String,
    pub password: String
}