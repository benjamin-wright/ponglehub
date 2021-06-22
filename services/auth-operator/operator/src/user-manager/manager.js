module.exports = class Manager {
    constructor(k8sApi, authClient) {
        this.k8sApi = k8sApi;
        this.authClient = authClient;

        this.mentalModel = [];
    }

    alreadyProcessing(user) {
        let found = false;

        this.mentalModel.forEach(userRequest => {
            if (
                userRequest.name === user.name &&
                userRequest.email === user.email &&
                userRequest.password === user.password
            ) {
                found = true;
            }
        });

        return found;
    }

    addToModel(user) {
        this.mentalModel.push(user);
        console.debug(`MentalModel: ${this.mentalModel.map(u => u.name)}`);
    }

    removeFromModel(user) {
        const idx = this.mentalModel.findIndex(u => u === user);
        this.mentalModel.splice(idx);
        console.debug(`MentalModel: ${this.mentalModel.map(u => u.name)}`);
    }

    async addUser(event) {
        const user = event.spec;

        if (this.alreadyProcessing(user)) {
            console.debug(`Ignoring user ${user.name} request, already in flight`);
            return;
        }

        this.addToModel(user);

        const id = await this.authClient.addUser(user);
        console.debug(`Got id ${id} for user ${user.name}`);

        try {
            await this.k8sApi.setUserId(event, id);
            console.debug(`Updated user ${user.name} status`);
        } finally {
            this.removeFromModel(user);
        }
    }

    async updateUser(event) {
        const user = event.spec;

        if (this.alreadyProcessing(user)) {
            console.debug(`Ignoring user ${user.name} request, already in flight`);
            return;
        }

        this.addToModel(user);

        try {
            const success = await this.authClient.updateUser(event.status.id, user);
            if (!success) {
                console.debug(`Ignoring user ${user.name} request, already up to date`);
            }
        } finally {
            this.removeFromModel(user);
        }
    }
};
