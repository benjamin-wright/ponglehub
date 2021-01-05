// const url = require('url');
const axios = require('axios');
const faker = require('faker');
const { fail } = require('assert');

describe('login', () => {
    it('should return a 404 if the token doesn\'t exist', async () => {
        const token = faker.random.uuid();

        try {
            await axios.get(`http://gatekeeper/login/${token}`);
            fail('expected the request to fail');
        } catch (err) {
            expect(err.response.status).toEqual(404);
        }
    });
});
