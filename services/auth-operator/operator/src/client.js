const axios = require('axios');
const AUTH_ENDPOINT = process.env.AUTH_ENDPOINT;

module.exports = class Client {
    constructor() {
        this.requests = [];
        this.processing = false;
    }

    addUser({ name, email, password }) {
        return new Promise((resolve, reject) => {
            this.requests.push(() => {
                axios.post(`${AUTH_ENDPOINT}/users`, { name, email, password })
                    .then(response => {
                        const id = response.data.id;
                        this.processing = false;
                        this.run();
                        resolve(id);
                    })
                    .catch(err => {
                        this.processing = false;
                        this.run();
                        reject(new Error(`Failed sending user to auth service: ${err.message}`));
                    });
            });

            this.run();
        });
    }

    run() {
        if (this.processing) {
            console.debug('[client] already running, deferring action');
            return;
        }

        const request = this.requests.shift();
        if (request) {
            this.processing = true;
            request();
        } else {
            console.debug('[client] finished last request');
            this.processing = false;
        }
    }
};
