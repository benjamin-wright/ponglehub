use postgres::{Client, NoTls, Error};
use std::fs;

use serde::{Serialize, Deserialize};

fn main() {
    println!("Running user migrations...");

    if let Err(e) = migrate("/config.yaml") {
        panic!("err: {:?}", e);
    }

    println!("Done");
}

fn migrate(config_path: &'static str) -> Result<(), Box<dyn std::error::Error>> {
    let conf = read_config(config_path)?;
    println!("Conf: {:?}", conf);

    let mut admin_client = connect_to_db(conf.host.as_str(), conf.port)?;

    for data in conf.data {
        create_user(&mut admin_client, data.user.as_str())?;
        create_database(&mut admin_client, data.database.as_str())?;
        grant_permissions(&mut admin_client, data.database.as_str(), data.user.as_str())?;
    }

    Ok(())
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
struct FileData {
    host: String,
    port: u16,
    data: Vec<DBData>
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
struct DBData {
    database: String,
    user: String
}

fn read_config(path: &'static str) -> Result<FileData, Box<dyn std::error::Error>> {
    let file_data = fs::read_to_string(path)?;

    let data: FileData = serde_yaml::from_str(file_data.as_str())?;
    Ok(data)
}

fn connect_to_db(host: &str, port: u16) -> Result<Client, Error> {
    let path = format!("postgresql://root@{}:{}", host, port);
    return Client::connect(path.as_ref(), NoTls);
}

fn create_user(client: &mut Client, user: &str) -> Result<(), Error> {
    let rows = client.query("SHOW USERS", &[])?;
    for row in rows {
        let name: &str = row.get("username");
        if name == user {
            println!("User {} already exists!", user);
            return Ok(());
        }
    }

    client.execute("CREATE USER $1", &[ &user ])?;
    println!("Created user: {}!", user);
    Ok(())
}

fn create_database(client: &mut Client, database: &str) -> Result<(), Error> {
    let rows = client.query("SELECT datname FROM pg_database", &[])?;
    for row in rows {
        let database_name: &str = row.get("datname");
        if database_name == database {
            println!("Database {} already exists!", database);
            return Ok(());
        }
    }

    let query = format!("CREATE DATABASE {}", database);
    client.execute(query.as_str(), &[])?;
    println!("Created database: {}!", database);
    Ok(())
}

fn grant_permissions(client: &mut Client, database: &str, user: &str) -> Result<(), Error> {
    let query = format!("GRANT ALL ON DATABASE {} TO {}", database, user);
    client.execute(query.as_str(), &[])?;
    println!("Granted '{}' permission to read/write to '{}'!", user, database);
    Ok(())
}
