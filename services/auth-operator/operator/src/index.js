const { Api, Listener } = require('./k8s');
const { Client } = require('./auth-server');
const { Manager } = require('./user-manager');

const client = new Client();
const listener = new Listener();
const api = new Api();
const manager = new Manager(api, client);

listener.reconcile(
    event => {
        if (event.status && event.status.id) {
            console.info(`Updating user: ${event.spec.name}`);
            manager.updateUser(event).catch(err => console.error(err));
            return;
        }

        console.info(`Adding new user: "${event.spec.name}"`);
        manager.addUser(event).catch(err => console.error(err));
    }, event => {
        console.log('something got deleted!');
    }, () => {
        console.log('fuck!');
    }
);

listener.start();

process.on('SIGINT', () => {
    listener.stop();
});
