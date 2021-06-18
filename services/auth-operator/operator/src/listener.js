const k8s = require('@kubernetes/client-node');
const namespace = process.env.NAMESPACE;

module.exports = class Listener {
    constructor() {
        const kc = new k8s.KubeConfig();
        kc.loadFromCluster();

        this.k8sApi = kc.makeApiClient(k8s.CustomObjectsApi);

        const listFn = () => this.k8sApi.listNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers');
        this.informer = k8s.makeInformer(kc, `/apis/ponglehub.co.uk/v1alpha1/namespaces/${namespace}/authusers`, listFn);
    }

    reconcile(upsertCallback, deleteCallback, errorCallback) {
        this.informer.on(k8s.ADD, upsertCallback);
        this.informer.on(k8s.UPDATE, () => {});
        this.informer.on(k8s.CHANGE, () => {});
        this.informer.on(k8s.DELETE, deleteCallback);
        this.informer.on(k8s.ERROR, errorCallback);
    }

    start() {
        this.informer.start()
            .then(() => {
                console.info('Informer started');
            })
            .catch(err => {
                console.error(`Something went wrong: ${err.message}`);
                console.info('Waiting 5 seconds and restarting');

                setTimeout(this.start, 5000);
            });
    }

    stop() {
        this.informer.stop();
    }
};
