const { Pool } = require('pg');
const connectionString = process.env.DB_URL;

class DB {
    constructor() {
        this.pool = new Pool({
            connectionString
        });
    }

    async addUser({ name, email, password, verified }) {
        await this.pool.query(
            'INSERT INTO users (name, email, password, verified) VALUES ($1, $2, $3, $4)',
            [name, email, password, verified]
        );
    }

    async getUsers() {
        const result = await this.pool.query(
            'SELECT id, name, email, password, verified FROM users'
        );

        return result.rows;
    }

    async clearUsers() {
        await this.pool.query('DELETE FROM users');
    }
}

module.exports = DB;
