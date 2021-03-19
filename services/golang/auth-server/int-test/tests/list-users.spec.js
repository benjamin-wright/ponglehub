const axios = require('axios');

describe('users route', () => {
    describe('list', () => {
        it('should return return all the users', async () => {
            await expect(
                axios.get('http://auth-server/users')
            ).resolves.toEqual({
                message: 'bingo!'
            });
        });
    });
});
