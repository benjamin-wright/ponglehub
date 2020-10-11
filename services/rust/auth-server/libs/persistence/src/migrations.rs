use postgres::{Client as PGClient, NoTls, Error};

pub fn run(user: &'static str, database: &'static str, hostname: &'static str, port: u16) -> Result<(), Error> {
    let mut admin_client = connect_to_db(hostname, port, "postgres")?;
    create_user(&mut admin_client, user)?;
    create_database(&mut admin_client, database)?;
    grant_permissions(&mut admin_client, database, user)?;

    drop(admin_client);

    let mut db_client = connect_to_db(hostname, port, database)?;
    create_table(&mut db_client, "users")?;

    Ok(())
}

fn connect_to_db(hostname: &'static str, port: u16, database: &'static str) -> Result<PGClient, Error> {
    let path = format!("postgresql://root@{}:{}/{}", hostname, port, database);
    return PGClient::connect(path.as_ref(), NoTls);
}

fn create_user(client: &mut PGClient, user: &'static str) -> Result<(), Error> {
    let rows = client.query("SHOW USERS", &[])?;
    for row in rows {
        let name: &str = row.get("username");
        if name == user {
            println!("Found user!");
            return Ok(());
        }
    }

    client.execute("CREATE USER $1", &[ &user ])?;
    println!("Created user!");
    Ok(())
}

fn create_database(client: &mut PGClient, database: &'static str) -> Result<(), Error> {
    let rows = client.query("SELECT datname FROM pg_database", &[])?;
    for row in rows {
        let database_name: &str = row.get("datname");
        if database_name == database {
            println!("Found database!");
            return Ok(());
        }
    }

    let query = format!("CREATE DATABASE {}", database);
    client.execute(query.as_str(), &[])?;
    println!("Created database!");
    Ok(())
}

fn grant_permissions(client: &mut PGClient, database: &'static str, user: &'static str) -> Result<(), Error> {
    let query = format!("GRANT ALL ON DATABASE {} TO {}", database, user);
    client.execute(query.as_str(), &[])?;
    println!("Granted permissions!");
    Ok(())
}

fn create_table(client: &mut PGClient, table: &'static str) -> Result<(), Error> {
    client.execute(
        "CREATE TABLE IF NOT EXISTS $1",
        &[ &table ]
    )?;
    println!("Created table!");
    Ok(())
}