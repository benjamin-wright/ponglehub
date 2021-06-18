const auth = require('./auth-mock');
const crds = require('./crds');
const faker = require('faker');

let counter = 0;

function makeUser() {
    return {
        name: `user-${counter++}`,
        username: faker.unique(faker.internet.userName),
        email: faker.unique(faker.internet.email),
        password: faker.unique(faker.internet.password)
    };
}

describe('user CRD', () => {
    beforeAll(async () => {
        await crds.deleteAll();
        await auth.reset();
    });

    describe('adding a new user', () => {
        it('should add the user to the database', async () => {
            const user = makeUser();
            await crds.addUser(user);
            await auth.waitForUser(user.username);
        });

        it('should add a different user to the database', async () => {
            const user = makeUser();
            await crds.addUser(user);
            await auth.waitForUser(user.username);
        });
    });
});
