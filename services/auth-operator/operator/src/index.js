const { Api, Listener } = require('./k8s');
const { Client } = require('./auth-server');
const { Manager, UserState, States } = require('./user-manager');

const client = new Client();
const listener = new Listener();
const api = new Api();
const state = new UserState();
const manager = new Manager(api, client, state);

listener.reconcile(
    event => {
        const name = `${event.metadata.name}-${event.metadata.namespace}`;
        const opCode = state.update(name, event.spec);

        switch (opCode) {
            case States.USER_ADDED:
                console.info(`Adding/updating user: "${name}"`);
                manager.addUser(event).catch(err => console.error(err));
                break;
            case States.USER_UPDATED:
                console.info(`Updating user: ${name}`);
                manager.updateUser(event).catch(err => console.error(err));
                break;
            default:
                console.info(`Nothing to do: ${name}`);
                break;
        }
    }, event => {
        const name = `${event.metadata.name}-${event.metadata.namespace}`;
        console.log(`Deleting user: ${name}`);
        state.remove(name);
    }, err => {
        console.error(`Something went wrong: ${err.message}`);
    }
);

listener.start();

process.on('SIGINT', () => {
    listener.stop();
});
