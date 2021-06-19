const axios = require('axios');
const AUTH_ENDPOINT = process.env.AUTH_ENDPOINT;

module.exports = class Client {
    constructor() {
        this.requests = [];
        this.current = undefined;
    }

    isProcessing(name) {
        if (this.current && this.current.user === name) {
            return true;
        }

        this.requests.forEach(request => {
            if (request.user === name) {
                return true;
            }
        });

        return false;
    }

    addUser({ name, email, password }) {
        if (this.isProcessing(name)) {
            console.debug(`Not adding user ${name}, a request is already in flight`);
            return Promise.resolve();
        }
        console.debug(`Queueing up request for user ${name}`);

        return new Promise((resolve, reject) => {
            const request = {
                user: name,
                action: () => {
                    axios.post(`${AUTH_ENDPOINT}/users`, { name, email, password })
                        .then(response => {
                            const id = response.data.id;
                            this.current = undefined;
                            this.run();
                            resolve(id);
                        })
                        .catch(err => {
                            this.current = undefined;
                            this.run();
                            reject(new Error(`Failed sending user to auth service: ${err.message}`));
                        });
                }
            };

            this.requests.push(request);

            this.run();
        });
    }

    run() {
        if (this.current) {
            console.debug('[client] already running, deferring action');
            return;
        }

        const request = this.requests.shift();
        if (request) {
            console.debug(`[client] running request for user ${request.user}`);
            this.current = request;
            request.action();
        } else {
            console.debug('[client] finished last request');
        }
    }
};
