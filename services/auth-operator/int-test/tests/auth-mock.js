const axios = require('axios');
const mockUrl = 'http://mock';
const async = require('@pongle/async');

module.exports = {
    getUsers,
    waitForUser,
    reset
};

async function getUsers() {
    return await axios
        .get(`${mockUrl}/users`)
        .then(res => res.data)
        .catch(err => { throw new Error(`Failed to list mock users: ${err}`); });
}

async function waitForUser(name) {
    let waiting = true;
    let waitedFor = 0;
    let mockUser;
    const timeout = 3000;

    while (waiting) {
        await async.sleep(50);
        waitedFor += 50;

        const users = await getUsers();

        users.forEach(user => {
            if (user.name === name) {
                waiting = false;
                mockUser = user;
            }
        });

        if (waiting && waitedFor > timeout) {
            throw new Error(`Timed out waiting for user ${name}: have ${JSON.stringify(users.map(u => u.name))}`);
        }
    };

    return mockUser;
}

async function reset() {
    return await axios
        .post(`${mockUrl}/reset`)
        .catch(err => { throw new Error(`Failed to reset mock users: ${err}`); });
}
