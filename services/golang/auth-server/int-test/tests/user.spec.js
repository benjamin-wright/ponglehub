const axios = require('axios');
const DB = require('../helpers/db');
const db = new DB();

describe('users route', () => {
    beforeEach(async () => {
        await db.clearUsers();
    });

    describe('get', () => {
        beforeEach(async () => {
            await db.addUser({ name: 'test-user', password: 'pwd', email: 'user@notathing.com', verified: false });
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

            expect(db.getUsers()).resolves.toEqual([
                {
                    id: expect.any(String),
                    name: 'random',
                    email: 'user@notathing.com',
                    password: 'pwd',
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
