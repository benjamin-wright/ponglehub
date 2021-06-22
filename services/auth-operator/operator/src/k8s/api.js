const k8s = require('@kubernetes/client-node');
const namespace = process.env.NAMESPACE;

module.exports = class Api {
    constructor() {
        const kc = new k8s.KubeConfig();
        kc.loadFromCluster();
        this.k8sApi = kc.makeApiClient(k8s.CustomObjectsApi);
    }

    setUserId(user, id) {
        const patch = [];
        const options = { headers: { 'Content-type': k8s.PatchUtils.PATCH_FORMAT_JSON_PATCH } };

        if (!user.status) {
            patch.push({
                op: 'add',
                path: '/status',
                value: { id }
            });
        } else if (!user.status.id) {
            patch.push({
                op: 'add',
                path: '/status/id',
                value: id
            });
        } else {
            patch.push({
                op: 'replace',
                path: '/status/id',
                value: id
            });
        }

        return this.k8sApi.patchNamespacedCustomObject(
            'ponglehub.co.uk',
            'v1alpha1',
            namespace,
            'authusers',
            user.metadata.name,
            patch,
            undefined, undefined, undefined, options
        ).catch(err => {
            throw new Error(`Failed to update user CRD ID for ${user.metadata.name} [${err.statusCode}]: ${err.message}`);
        });
    }
};
