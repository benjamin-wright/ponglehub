expect.extend({
    toMatchError(received, expected) {
        if (received.message === expected) {
            return {
                pass: true,
                message: () => `Expected "${received.message}" not to match "${expected}"`
            };
        } else {
            return {
                pass: false,
                message: () => `Expected "${received.message}" to match "${expected}"`
            };
        }
    }
});
