const axios = require('axios');
const DB = require('../helpers/db');
const db = new DB();
const { compareStrings } = require('./helpers');

describe('users route', () => {
    beforeEach(async () => {
        await db.clearUsers();
    });

    describe('get', () => {
        beforeEach(async () => {
            await db.addUser({ name: 'test-user', password: 'pwd', email: 'user@notathing.com', verified: false });
            await db.addUser({ name: 'verified-user', password: 'pwd2', email: 'user2@notathing.com', verified: true });
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

        it('should return a verified user if it does exist', async () => {
            await expect(
                axios.get('http://auth-server/user/verified-user').then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: {
                    id: expect.any(String),
                    name: 'verified-user',
                    email: 'user2@notathing.com',
                    verified: true
                }
            });
        });
    });

    describe('post', () => {
        it('should add the user to the database', async () => {
            await expect(
                axios.post('http://auth-server/user', { name: 'user', password: 'pwd', email: 'user@notathing.com' }).then(res => ({ status: res.status }))
            ).resolves.toEqual({
                status: 202
            });

            await expect(db.getUsers()).resolves.toEqual([
                {
                    id: expect.any(String),
                    name: 'user',
                    email: 'user@notathing.com',
                    password: 'pwd',
                    verified: false
                }
            ]);
        });

        it('should fail to add a user with duplicate name', async () => {
            await db.addUser({ name: 'duplicate', email: 'same@email.com', password: 'pwd', verified: false });

            await expect(
                axios.post('http://auth-server/user', { name: 'duplicate', password: 'pwd', email: 'duplicate@notathing.com' }).then(res => ({ status: res.status }))
            ).rejects.toMatchError('Request failed with status code 400');
        });

        it('should fail to add a user with duplicate email', async () => {
            await db.addUser({ name: 'duplicate', email: 'same@email.com', password: 'pwd', verified: false });

            await expect(
                axios.post('http://auth-server/user', { name: 'different', password: 'pwd', email: 'same@email.com' }).then(res => ({ status: res.status }))
            ).rejects.toMatchError('Request failed with status code 400');
        });
    });

    describe('list', () => {
        it('should return empty if no users', async () => {
            await expect(
                axios.get('http://auth-server/user').then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: []
            });
        });

        it('should return all users when present', async () => {
            const users = ['user1', 'user2', 'user3'];

            await db.addUsers(users.map(user => ({ name: user, email: `${user}@email.com`, password: 'pwd', verified: false })));

            await expect(
                axios.get('http://auth-server/user').then(res => ({ status: res.status, data: res.data.sort((a, b) => compareStrings(a.name, b.name)) }))
            ).resolves.toEqual({
                status: 200,
                data: users.map(user => ({
                    id: expect.any(String),
                    name: user,
                    email: `${user}@email.com`,
                    verified: false
                }))
            });
        });
    });
});
