# ASYNC

## Examples

### sleep

sleep for some time
```js
const { sleep } = require('@pongle/async');

await sleep(500); //milliseconds
```

### waitFor

Keep trying something until is resolves or the timeout expires
```js
const { waitFor } = require('@pongle/async');
const TIMEOUT = 3000; //milliseconds

await waitFor(async () => {
  await mightNoPassImmediately();
}, TIMEOUT);
```