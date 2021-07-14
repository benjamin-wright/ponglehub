const Manager = require('./manager');

describe('manager', () => {
    let manager;
    let client;
    let k8s;

    beforeEach(() => {
        k8s = {
            setUserId: jest.fn(async () => true)
        };

        client = {
            addUser: jest.fn(async () => 'abcde')
        };
        manager = new Manager(k8s, client);
    });

    it('should add the user', async () => {
        await manager.addUser({ metadata: { name: 'user1', namespace: 'my-namespace' }, spec: { name: 'someuser', email: 'whatevs', password: 'you-too' } });

        expect(client.addUser.mock.calls.length).toEqual(1);
    });
});
