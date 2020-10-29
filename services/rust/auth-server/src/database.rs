use rocket_contrib::databases::postgres;

#[database("auth")]
pub struct AuthDB (postgres::Connection);

pub fn fairing() -> impl rocket::fairing::Fairing {
    AuthDB::fairing()
}