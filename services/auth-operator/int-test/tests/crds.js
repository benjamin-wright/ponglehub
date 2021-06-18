const k8s = require('@kubernetes/client-node');

const namespace = process.env.NAMESPACE;
const kc = new k8s.KubeConfig();
kc.loadFromCluster();
const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi);

module.exports = {
    addUser,
    deleteAll
};

async function addUser({ name, username, email, password }) {
    await k8sApi.createNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers', {
        apiVersion: 'ponglehub.co.uk/v1alpha1',
        kind: 'AuthUser',
        metadata: {
            name,
            namespace: 'int-auth-operator'
        },
        spec: {
            name: username,
            email,
            password
        },
        status: {
            id: 'something'
        }
    }).catch(err => {
        throw new Error(`Failed to get user: ${err.statusCode}\n${JSON.stringify(err.body, null, 2)}`);
    });
}

async function deleteAll() {
    await k8sApi.deleteCollectionNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers').catch(err => {
        throw new Error(`Failed to delete users: ${err.statusCode}\n${JSON.stringify(err.body, null, 2)}`);
    });
}
