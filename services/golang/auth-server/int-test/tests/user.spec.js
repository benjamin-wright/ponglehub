const axios = require('axios');

async function listUsers() {
    const result = await axios.get('http://auth-server/user');

    if (result.status !== 200) {
        throw new Error(`Failed listing users: ${result.status}`);
    }

    return result.data;
}

async function deleteUser(username) {
    const result = await axios.delete(`http://auth-server/user/${username}`);

    if (result.status !== 202) {
        throw new Error(`Failed deleting user: ${result.status}`);
    }
}

async function expectUsers(users) {
    await expect(
        axios.get('http://auth-server/user').then(res => ({ status: res.status, data: res.data }))
    ).resolves.toEqual({
        status: 200,
        data: users
    });
}

describe('users route', () => {
    beforeEach(async () => {
        const users = await listUsers();
        await Promise.all(users.map(user => deleteUser(user.name)));
    });

    describe('get', () => {
        beforeEach(async () => {
            await expect(
                axios.post('http://auth-server/user', { name: 'test-user', password: 'pwd', email: 'user@notathing.com' }).then(res => ({ status: res.status }))
            ).resolves.toEqual({
                status: 202
            });
        });

        it('should return 404 if user doesn\'t exist', async () => {
            await expect(
                axios.get('http://auth-server/user/something')
            ).rejects.toMatchError('Request failed with status code 404');
        });

        it('should return the user if it does exist', async () => {
            await expect(
                axios.get('http://auth-server/user/test-user').then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: {
                    id: expect.any(String),
                    name: 'test-user',
                    email: 'user@notathing.com',
                    verified: false
                }
            });
        });
    });

    describe('post', () => {
        it('should add the user to the database', async () => {
            await expect(
                axios.post('http://auth-server/user', { name: 'random', password: 'pwd', email: 'user@notathing.com' }).then(res => ({ status: res.status }))
            ).resolves.toEqual({
                status: 202
            });

            await expectUsers([
                {
                    id: expect.any(String),
                    name: 'random',
                    email: 'user@notathing.com',
                    verified: false
                }
            ]);
        });
    });

    describe('list', () => {
        it('should return return all the users', async () => {
            await expect(
                axios.get('http://auth-server/user').then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: []
            });
        });
    });
});
