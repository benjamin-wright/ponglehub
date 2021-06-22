const auth = require('./auth-mock');
const crds = require('./crds');
const faker = require('faker');
const async = require('@pongle/async');

function makeUser() {
    return {
        meta: {
            name: `user-${faker.unique(() => faker.random.alphaNumeric(10))}`
        },
        spec: {
            name: faker.unique(faker.internet.userName),
            email: faker.unique(faker.internet.email),
            password: faker.unique(faker.internet.password)
        }
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
            await auth.waitForUser(user.spec.name);
        });

        it('should add a different user to the database', async () => {
            const user = makeUser();
            await crds.addUser(user);
            await auth.waitForUser(user.spec.name);
        });
    });

    describe('updating a user', () => {
        let user;

        beforeEach(async () => {
            user = makeUser();
            await crds.addUser(user);
            const mockUser = await auth.waitForUser(user.spec.name);
            await crds.waitForId(user.meta.name, mockUser.id);
        });

        it('should update the user name in the database', async () => {
            user.spec.name = 'updated.name';

            await crds.updateUser(user.meta.name, user.spec);
            await auth.waitForUser(user.spec.name);
        });

        it('should update the user email in the database', async () => {
            user.spec.email = 'updated@email';

            await crds.updateUser(user.meta.name, user.spec);

            await async.waitFor(async () => {
                const actual = await auth.waitForUser(user.spec.name);
                expect(actual.email).toEqual(user.spec.email);
            });
        });
    });
});
