const { Pool } = require('pg');
const connectionString = process.env.DB_URL;

class DB {
    constructor() {
        this.pool = new Pool({
            connectionString
        });
    }

    async addUser({ name, email, password, verified }) {
        const result = await this.pool.query(
            'INSERT INTO users (name, email, password, verified) VALUES ($1, $2, $3, $4) RETURNING id',
            [name, email, password, verified]
        );

        if (result.rows.length !== 1) {
            throw new Error(`Error adding user: expected 1 row response, got ${result.rows.length}`);
        }

        return result.rows[0].id;
    }

    async addUsers(users) {
        await Promise.all(users.map(user => this.addUser(user)));
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
