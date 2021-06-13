const k8s = require('@kubernetes/client-node');

const namespace = process.env.NAMESPACE;

const kc = new k8s.KubeConfig();
kc.loadFromCluster();

const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi);

const listFn = () => k8sApi.listNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers');

const informer = k8s.makeInformer(kc, `/apis/ponglehub.co.uk/v1alpha1/namespaces/${namespace}/authusers`, listFn);

informer.on(k8s.ADD, () => console.log('something got added!'));
informer.on(k8s.UPDATE, () => console.log('something got updated!'));
informer.on(k8s.CHANGE, () => console.log('something got changed!'));
informer.on(k8s.DELETE, () => console.log('something got deleted!'));
informer.on(k8s.ERROR, () => console.log('fuck!'));

function start() {
    informer.start()
        .then(() => {
            console.info('Informer started');
        })
        .catch(err => {
            console.error(`Something went wrong: ${err.message}`);
            console.info('Waiting 5 seconds and restarting');

            setTimeout(start, 5000);
        });
}

start();

process.on('SIGINT', () => {
    informer.stop();
});
