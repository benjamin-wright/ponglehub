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

async function waitForUser(username) {
    let waiting = true;
    let waitedFor = 0;
    const timeout = 3000;

    while (waiting) {
        await async.sleep(50);
        waitedFor += 50;

        const users = await getUsers();

        users.forEach(user => {
            if (user.name === username) {
                waiting = false;
            }
        });

        if (waiting && waitedFor > timeout) {
            throw new Error(`Timed out waiting for user ${username}: have ${JSON.stringify(users.map(u => u.name))}`);
        }
    };
}

async function reset() {
    return await axios
        .post(`${mockUrl}/reset`)
        .catch(err => { throw new Error(`Failed to reset mock users: ${err}`); });
}
