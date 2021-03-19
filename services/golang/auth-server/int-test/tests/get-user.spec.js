const axios = require('axios');

describe('users route', () => {
    describe('get', () => {
        it('should return 404 if user doesn\'t exist', async () => {
            await expect(
                axios.get('http://auth-server/users/something')
            ).rejects.toMatchError('Request failed with status code 404');
        });
    });
});
