const axios = require('axios');

async function listUsers() {
    const result = await axios.get('http://auth-server/user');

    if (result.status !== 200) {
        throw new Error(`Failed listing users: ${result.status}`);
    }

    return result.data;
}

async function deleteUser(username) {
    console.log(`Deleting user: ${username}`);
    const result = await axios.delete(`http://auth-server/user/${username}`);

    if (result.status !== 202) {
        throw new Error(`Failed deleting user: ${result.status}`);
    }
}

describe('users route', () => {
    beforeEach(async () => {
        const users = await listUsers();
        await Promise.all(users.map(user => deleteUser(user.name)));
    });

    describe('get', () => {
        it('should return 404 if user doesn\'t exist', async () => {
            await expect(
                axios.get('http://auth-server/user/something')
            ).rejects.toMatchError('Request failed with status code 404');
        });
    });

    describe('post', () => {
        it('should return a dummy ok code', async () => {
            await expect(
                axios.post('http://auth-server/user', { name: 'random', password: 'pwd', email: 'user@notathing.com' }).then(res => ({ status: res.status }))
            ).resolves.toEqual({
                status: 202
            });
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
