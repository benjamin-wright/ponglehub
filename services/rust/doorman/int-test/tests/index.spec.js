const url = require('url');
const axios = require('axios');
const MockServer = require('./mockServer');
const parser = require('node-html-parser');
const { fail } = require('assert');

describe('login page', () => {
    let mockServer = null;

    beforeEach(async () => {
        mockServer = new MockServer('app.ponglehub.co.uk');
        mockServer.addRoute({
            path: '/login',
            host: 'mock-gatekeeper',
            body: { token: 'test-token' },
            status: 200
        });
        await mockServer.start();
    });

    afterEach(async () => {
        await mockServer.stop();
    });

    it('should add token and redirect hidden fields', async () => {
        const response = await axios.get('http://doorman/login?redirect=test-redirect');
        const page = parser.parse(response.data);

        expect({
            token: page.querySelector('input[name="token"]').getAttribute('value'),
            redirect: page.querySelector('input[name="redirect"]').getAttribute('value')
        }).toEqual({
            token: 'test-token',
            redirect: 'test-redirect'
        });
    });
});

describe('login api endpoint', () => {
    let mockServer = null;

    beforeEach(async () => {
        mockServer = new MockServer('app.ponglehub.co.uk');
        await mockServer.start();
    });

    afterEach(async () => {
        await mockServer.stop();
    });

    describe('login token is good', () => {
        beforeEach(() => {
            mockServer.addRoute({
                path: '/login/test-token',
                host: 'mock-gatekeeper',
                body: {},
                status: 200
            });
        });

        it('should redirect properly when username and password are provided', async () => {
            const params = new url.URLSearchParams({
                token: 'test-token',
                redirect: 'test-redirect',
                username: 'test-user',
                password: 'test-pass'
            });

            try {
                await axios.post(
                    'http://doorman/api/login',
                    params.toString(),
                    { maxRedirects: 0 }
                );

                fail('should have failed to redirect');
            } catch (err) {
                expect(mockServer.calls).toEqual([{
                    host: 'mock-gatekeeper',
                    method: 'GET',
                    path: '/login/test-token'
                }]);
                expect(err.response.status).toEqual(302);
                expect(err.response.headers['set-cookie']).toEqual(['pongle_auth=special_token; Domain=ponglehub.co.uk']);
            }
        });
    });
});
