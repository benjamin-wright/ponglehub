const url = require('url');
const axios = require('axios');
const MockServer = require('./mockServer');
const parser = require('node-html-parser');

describe('login page', () => {
    let mockServer = null;

    beforeEach(async () => {
        mockServer = new MockServer();
        mockServer.addRoute('/login', { token: 'test-token' }, 200);
        await mockServer.start();
    });

    afterEach(async () => {
        await mockServer.stop();
    });

    it('should work', async () => {
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
        mockServer = new MockServer();
        mockServer.addRoute('/login/test-token', {}, 200);
        mockServer.addRoute('/redirect-target', 'redirected', 200);
        await mockServer.start();
    });

    afterEach(async () => {
        await mockServer.stop();
    });

    it('should pass', async () => {
        const params = new url.URLSearchParams({
            token: 'test-token',
            redirect: 'http://localhost/redirect-target',
            username: 'test-user',
            password: 'test-pass'
        });

        const response = await axios.post(
            'http://doorman/api/login',
            params.toString()
        );

        expect(response.data).toEqual('redirected');
        expect(response.cookies).toEqual({});
    });
});
