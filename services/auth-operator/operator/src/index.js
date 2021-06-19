const k8s = require('./k8s');
const Client = require('./client');

const client = new Client();
const listener = new k8s.Listener();
const api = new k8s.Api();

listener.reconcile(
    event => {
        if (event.status && event.status.id) {
            console.debug(`User ${event.metadata.name} already has an id: ${event.status.id}`);
            return;
        }

        console.info(`Adding new user: "${event.spec.name}"`);
        client.addUser(event.spec).then(id => {
            if (id) {
                console.info(`Got new id for user "${event.spec.name}": "${id}"`);
                return api.setUserId(event, id);
            }
        }).catch(err => {
            if (err.body) {
                console.error(`Failed with status [${err.statusCode}]: ${JSON.stringify(err.body.message)}`);
            } else {
                console.error(err);
            }
        });
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
