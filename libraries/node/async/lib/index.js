const DEFAULT_TIMEOUT = 4000;
const POLLING_TIMEOUT = 100;

module.exports = {
    sleep,
    waitFor
};

function sleep(timeout) {
    return new Promise(resolve => {
        setTimeout(resolve, timeout);
    });
}

async function waitFor(func, timeout = DEFAULT_TIMEOUT) {
    let elapsed = 0;
    while (true) {
        try {
            return await func();
        } catch (err) {
            await sleep(POLLING_TIMEOUT);
            elapsed += POLLING_TIMEOUT;

            if (elapsed > timeout) {
                throw new Error(`Failed waiting for function to return: ${err}`);
            }
        }
    }
}
