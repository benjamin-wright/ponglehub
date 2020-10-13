use postgres::{Client as PGClient, NoTls, Error};

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}

pub struct Client {
    client: PGClient
}

#[derive(Debug, Clone)]
pub struct DBError;

impl Client {
    pub fn new(user: &'static str, database: &'static str, hostname: &'static str, port: u16) -> Result<Client, Error> {
        let path = format!("postgresql://{}@{}:{}/{}", user, hostname, port, database);
        let client = PGClient::connect(path.as_ref(), NoTls)?;

        return Ok(Client{ client });
    }

    pub fn get_tables(&mut self) -> Result<Vec<String>, Error> {
        let rows = self.client.query("SHOW TABLES", &[])?;

        let mut result: Vec<String> = vec!();
        for row in rows {
            result.push(row.get("table_name"));
        }

        return Ok(result);
    }
}
