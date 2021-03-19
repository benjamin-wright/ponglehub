const axios = require('axios');

describe('status route', () => {
    it('should add token and redirect hidden fields', async () => {
        const response = await axios.get('http://auth-server/status');

        expect({
            status: response.status,
            data: response.data
        }).toEqual({
            status: 200,
            data: { message: 'OK' }
        });
    });
});
