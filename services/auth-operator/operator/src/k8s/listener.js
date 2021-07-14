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
        this.informer.on(k8s.ADD, event => { console.debug(`ADD - ${event.metadata.name}`); upsertCallback(event); });
        this.informer.on(k8s.UPDATE, event => {
            console.debug(`UPDATE - ${event.metadata.name}`);
            upsertCallback(event);
        });
        this.informer.on(k8s.CHANGE, event => {
            console.debug(`CHANGE - ${event.metadata.name}`);
            upsertCallback(event);
        });
        this.informer.on(k8s.DELETE, event => { console.debug(`DELETE - ${event.metadata.name}`); deleteCallback(event); });
        this.informer.on(k8s.ERROR, err => { console.debug('USERWATCH ERROR'); errorCallback(err); });
    }

    start() {
        this.informer.start()
            .then(() => {
                console.info('Informer started');
            })
            .catch(err => {
                console.error(`Something went wrong: ${err.message}`);
                console.info('Waiting 5 seconds and restarting');

                setTimeout(this.start.bind(this), 5000);
            });
    }

    stop() {
        this.informer.stop();
    }
};
