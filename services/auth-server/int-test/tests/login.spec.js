const axios = require('axios');
const DB = require('../helpers/db');
const db = new DB();
const bcrypt = require('bcrypt');

describe('login routes', () => {
    beforeEach(async () => {
        await db.clearUsers();
    });

    describe('login', () => {
        it('should fail if the email is missing', async () => {
            await expect(
                axios.post('http://login.int-auth-server.svc.cluster.local', { password: 'pwd' }).then(res => ({ status: res.status }))
            ).rejects.toMatchError('Request failed with status code 400');
        });

        it('should fail if the password is missing', async () => {
            await expect(
                axios.post('http://login.int-auth-server.svc.cluster.local', { email: 'whatevs' }).then(res => ({ status: res.status }))
            ).rejects.toMatchError('Request failed with status code 400');
        });

        it('should fail if the user doesn\'t exist', async () => {
            await expect(
                axios.post('http://login.int-auth-server.svc.cluster.local', { password: 'pwd', email: 'user@notathing.com' }).then(res => ({ status: res.status }))
            ).rejects.toMatchError('Request failed with status code 401');
        });

        it('should fail if the user isn\'t verified', async () => {
            const hash = await bcrypt.hash('input-password', 10);
            await db.addUser({ name: 'test-user', password: hash, email: 'user@exists.com', verified: false });

            await expect(
                axios.post('http://login.int-auth-server.svc.cluster.local', { password: 'input-password', email: 'user@exists.com' }).then(res => ({ status: res.status, data: res.data }))
            ).rejects.toMatchError('Request failed with status code 401');
        });

        it('should fail if the password is wrong', async () => {
            const hash = await bcrypt.hash('input-password', 10);
            await db.addUser({ name: 'test-user', password: hash, email: 'user@exists.com', verified: true });

            await expect(
                axios.post('http://login.int-auth-server.svc.cluster.local', { password: 'wrong-password', email: 'user@exists.com' }).then(res => ({ status: res.status, data: res.data }))
            ).rejects.toMatchError('Request failed with status code 401');
        });

        it('should pass if everything works', async () => {
            const hash = await bcrypt.hash('input-password', 10);
            await db.addUser({ name: 'test-user', password: hash, email: 'user@exists.com', verified: true });

            await expect(
                axios.post('http://login.int-auth-server.svc.cluster.local', { password: 'input-password', email: 'user@exists.com' }).then(res => ({ status: res.status, data: res.data }))
            ).resolves.toEqual({
                status: 200,
                data: {
                    id: expect.any(String),
                    name: 'test-user'
                }
            });
        });
    });
});
