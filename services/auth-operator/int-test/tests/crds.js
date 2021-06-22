const k8s = require('@kubernetes/client-node');
const async = require('@pongle/async');

const namespace = process.env.NAMESPACE;
const kc = new k8s.KubeConfig();
kc.loadFromCluster();
const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi);

module.exports = {
    addUser,
    updateUser,
    waitForId,
    deleteAll
};

async function addUser({ meta, spec }) {
    await k8sApi.createNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers', {
        apiVersion: 'ponglehub.co.uk/v1alpha1',
        kind: 'AuthUser',
        metadata: {
            name: meta.name,
            namespace: 'int-auth-operator'
        },
        spec
    }).catch(err => {
        throw new Error(`Failed to get user: ${err.statusCode}\n${JSON.stringify(err.body, null, 2)}`);
    });
}

async function updateUser(metaName, updates) {
    const user = await getUser(metaName);

    user.spec = {
        ...user.spec,
        ...updates
    };

    await k8sApi.replaceNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers', metaName, user).catch(err => {
        throw new Error(`Failed to get user: ${err.statusCode}\n${JSON.stringify(err.body, null, 2)}`);
    });
}

async function deleteAll() {
    await k8sApi.deleteCollectionNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers').catch(err => {
        throw new Error(`Failed to delete users: ${err.statusCode}\n${JSON.stringify(err.body, null, 2)}`);
    });
}

async function getUser(metaName) {
    const response = await k8sApi.getNamespacedCustomObject('ponglehub.co.uk', 'v1alpha1', namespace, 'authusers', metaName);
    return response.body;
}

async function waitForId(metaName, id) {
    let waiting = true;
    let waitedFor = 0;
    let lastError;
    let user;
    const timeout = 3000;

    while (waiting) {
        await async.sleep(50);
        waitedFor += 50;

        try {
            user = await getUser(metaName);
            if (user.status && user.status.id === id) {
                waiting = false;
                continue;
            }
        } catch (err) {
            lastError = err;
            continue;
        }

        if (waiting && waitedFor > timeout) {
            if (lastError) {
                throw new Error(`Timed out waiting for user ${metaName} to have id ${id}: ${lastError.message}`);
            } else {
                throw new Error(`Timed out waiting for user ${metaName} to have id ${id}`);
            }
        }
    };
}
