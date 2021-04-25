const axios = require('axios');

describe.skip('status route', () => {
    it('should add token and redirect hidden fields', async () => {
        const response = await axios.get('http://auth-server:8080/status');

        expect({
            status: response.status,
            data: response.data
        }).toEqual({
            status: 200,
            data: { message: 'OK' }
        });
    });
});
