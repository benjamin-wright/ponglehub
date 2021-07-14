module.exports = class Manager {
    constructor(k8sApi, authClient) {
        this.k8sApi = k8sApi;
        this.authClient = authClient;
    }

    async addUser(event) {
        const name = `${event.metadata.name}-${event.metadata.namespace}`;
        const user = event.spec;

        const id = await this.authClient.addUser(user);
        console.debug(`Got id ${id} for user ${name}`);

        await this.k8sApi.setUserId(event, id);
        console.debug(`Updated user ${name} status`);
    }

    async updateUser(event) {
        const name = `${event.metadata.name}-${event.metadata.namespace}`;
        const user = event.spec;

        const success = await this.authClient.updateUser(event.status.id, user);
        if (!success) {
            console.debug(`Ignoring user ${name} request, already up to date`);
        }
    }
};
