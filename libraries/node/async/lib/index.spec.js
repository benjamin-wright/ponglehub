const async = require('./index');

beforeAll(() => {
    jest.useFakeTimers();
});

describe('sleep', () => {
    it('should not return before timeout elapses', async () => {
        const promise = async.sleep(500);
        promise.then(() => { promise.done = true; });

        await Promise.resolve();

        return expect(promise.done).not.toBe(true);
    });

    it('should return after the timeout elapses', async () => {
        const promise = async.sleep(500);
        promise.then(() => { promise.done = true; });

        jest.runAllTimers();
        await Promise.resolve();

        return expect(promise.done).toBe(true);
    });
});

describe('waitFor', () => {
    it('should return immediately if the function passes', async () => {
        const func = jest.fn(() => Promise.resolve('result'));
        const promise = async.waitFor(func);

        await Promise.resolve();

        expect(func).toHaveBeenCalledTimes(1);
        await expect(promise).resolves.toBe('result');
    });

    it('should wait if the function fails the first time', async () => {
        const func = jest.fn(() => Promise.reject(new Error('oops')));
        func.mockReturnValueOnce(Promise.reject(new Error('oops'))).mockReturnValue(Promise.resolve('result'));
        const promise = async.waitFor(func);

        await Promise.resolve();
        jest.runAllTimers();
        await Promise.resolve();

        expect(func).toHaveBeenCalledTimes(2);
        await expect(promise).resolves.toBe('result');
    });

    it('should error if the function takes too long to return', async () => {
        const func = jest.fn(() => Promise.reject(new Error('oops')));
        func.mockReturnValueOnce(Promise.reject(new Error('oops'))).mockReturnValue(Promise.resolve('result'));
        const promise = async.waitFor(func, 50);

        await Promise.resolve();
        jest.runAllTimers();
        expect(func).toHaveBeenCalledTimes(1);
        await expect(promise).rejects.toEqual(new Error('Failed waiting for function to return: Error: oops'));
    });
});
