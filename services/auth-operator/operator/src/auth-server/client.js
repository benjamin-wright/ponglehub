const axios = require('axios');
const AUTH_ENDPOINT = process.env.AUTH_ENDPOINT;

module.exports = class Client {
    async addUser({ name, email, password }) {
        try {
            const response = await axios.post(`${AUTH_ENDPOINT}/users`, { name, email, password });
            return response.data.id;
        } catch (err) {
            throw new Error(`Failed to add user ${err.statusCode}: ${err.message}`);
        }
    }

    async getUser(id) {
        try {
            const response = await axios.get(`${AUTH_ENDPOINT}/users/${id}`);
            return response.data;
        } catch (err) {
            throw new Error(`Failed to get user ${err.statusCode}: ${err.message}`);
        }
    }

    async updateUser(id, { name, email }) {
        const user = await this.getUser(id);
        if (user.name === name && user.email === email) {
            return false;
        }

        try {
            await axios.put(`${AUTH_ENDPOINT}/users/${id}`, { name, email });
            return true;
        } catch (err) {
            throw new Error(`Failed to update user ${id} (${name}) ${err.statusCode}: ${err.message}`);
        }
    }
};
