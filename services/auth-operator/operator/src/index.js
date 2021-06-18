const Listener = require('./listener');
const Client = require('./client');
const k8s = require('@kubernetes/client-node');

const client = new Client();
const listener = new Listener();
const namespace = process.env.NAMESPACE;

listener.reconcile(
    event => {
        client.addUser(event.spec).then(id => {
            console.info('Got user id: ' + id);
            const options = { headers: { 'Content-type': k8s.PatchUtils.PATCH_FORMAT_JSON_PATCH } };
            const patch = [
                {
                    op: 'replace',
                    path: '/status/id',
                    value: id
                }
            ];

            return listener.k8sApi.patchNamespacedCustomObject(
                'ponglehub.co.uk',
                'v1alpha1',
                namespace,
                'authusers',
                event.metadata.name,
                patch,
                undefined, undefined, undefined, options
            );
        }).catch(err => {
            if (err.body) {
                console.error(`Failed with status: ${err.statusCode}: ${JSON.stringify(err.body.message)}`);
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
