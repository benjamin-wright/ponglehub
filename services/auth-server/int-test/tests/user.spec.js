const axios = require('axios');
const DB = require('../helpers/db');
const db = new DB();
const { compareStrings } = require('../helpers/strings');
const missingId = '00000000-0000-0000-0000-000000000000';

describe('users route', () => {
    beforeEach(async () => {
        await db.clearUsers();
    });

    describe('get', () => {
        let testUser;
        let verifiedUser;

        beforeEach(async () => {
            testUser = await db.addUser({ name: 'test-user', password: 'pwd', email: 'user@notathing.com', verified: false });
            verifiedUser = await db.addUser({ name: 'verified-user', password: 'pwd2', email: 'user2@notathing.com', verified: true });
        });

        it('should return 404 if user doesn\'t exist', async () => {
            await expect(
                axios.get(`http://auth-server/users/${missingId}`)
            ).rejects.toMatchError('Request failed with status code 404');
        });

        it('should return the user if it does exist', async () => {
            await expect(
                axios.get(`http://auth-server/users/${testUser}`).then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: {
                    name: 'test-user',
                    email: 'user@notathing.com',
                    verified: false
                }
            });
        });

        it('should return a verified user if it does exist', async () => {
            await expect(
                axios.get(`http://auth-server/users/${verifiedUser}`).then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: {
                    name: 'verified-user',
                    email: 'user2@notathing.com',
                    verified: true
                }
            });
        });
    });

    describe('put', () => {

    });

    describe('post', () => {
        it('should add the user to the database', async () => {
            await expect(
                axios.post('http://auth-server/users', { name: 'user', password: 'pwd', email: 'user@notathing.com' }).then(res => ({ status: res.status }))
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
                axios.post('http://auth-server/users', { name: 'duplicate', password: 'pwd', email: 'duplicate@notathing.com' })
            ).rejects.toMatchError('Request failed with status code 400');
        });

        it('should fail to add a user with duplicate email', async () => {
            await db.addUser({ name: 'duplicate', email: 'same@email.com', password: 'pwd', verified: false });

            await expect(
                axios.post('http://auth-server/users', { name: 'different', password: 'pwd', email: 'same@email.com' })
            ).rejects.toMatchError('Request failed with status code 400');
        });
    });

    describe('delete', () => {
        it('should return a 404 if the user doesn\'t exist', async () => {
            await expect(
                axios.delete(`htth://auth-server/users/${missingId}`)
            ).rejects.toMatchError('Request failed with status code 404');
        });

        it('should delete the specified user', async () => {
            await db.addUser({ name: 'not-delete', email: 'not@email.com', password: 'pwd', verified: false });
            const id = await db.addUser({ name: 'to-delete', email: 'do@email.com', password: 'pwd', verified: false });

            await expect(
                axios.delete(`http://auth-server/users/${id}`).then(res => ({ status: res.status }))
            ).resolves.toEqual({ status: 204 });

            expect((await db.getUsers()).map(x => x.name)).toEqual(['not-delete']);
        });
    });

    describe('list', () => {
        it('should return empty if no users', async () => {
            await expect(
                axios.get('http://auth-server/users').then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: []
            });
        });

        it('should return all users when present', async () => {
            const users = ['user1', 'user2', 'user3'];

            await db.addUsers(users.map(user => ({ name: user, email: `${user}@email.com`, password: 'pwd', verified: false })));

            await expect(
                axios.get('http://auth-server/users').then(res => ({ status: res.status, data: res.data.sort((a, b) => compareStrings(a.name, b.name)) }))
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
