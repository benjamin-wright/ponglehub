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
        await manager.addUser({ spec: { name: 'someuser', email: 'whatevs', password: 'you-too' } });

        expect(client.addUser.mock.calls.length).toEqual(1);
    });

    it('should add the same user multiple times sequentially', async () => {
        await manager.addUser({ spec: { name: 'someuser', email: 'whatevs', password: 'you-too' } });
        await manager.addUser({ spec: { name: 'someuser', email: 'whatevs', password: 'you-too' } });

        expect(client.addUser.mock.calls.length).toEqual(2);
    });

    it('should not add the same user multiple times simultaneously', async () => {
        const promise1 = manager.addUser({ spec: { name: 'someuser', email: 'whatevs', password: 'you-too' } });
        const promise2 = manager.addUser({ spec: { name: 'someuser', email: 'whatevs', password: 'you-too' } });

        await promise1;
        await promise2;

        expect(client.addUser.mock.calls.length).toEqual(1);
    });

    it('should add different users multiple times simultaneously', async () => {
        const promise2 = manager.addUser({ spec: { name: 'someuser', email: 'whatevs', password: 'you-too' } });
        const promise1 = manager.addUser({ spec: { name: 'anotheruser', email: 'whatevs', password: 'you-too' } });

        await promise1;
        await promise2;

        expect(client.addUser.mock.calls.length).toEqual(2);
    });
});
