export const wasmBrowserInstantiate = async (wasmModuleUrl, importObject) => {
  let response = undefined;

  // Check if the browser supports streaming instantiation
  if (WebAssembly.instantiateStreaming) {
    // Fetch the module, and instantiate it as it is downloading
    response = await WebAssembly.instantiateStreaming(
      fetch(wasmModuleUrl),
      importObject
    );
  } else {
    // Fallback to using fetch to download the entire module
    // And then instantiate the module
    const fetchAndInstantiateTask = async () => {
      const wasmArrayBuffer = await fetch(wasmModuleUrl).then(response =>
        response.arrayBuffer()
      );
      return WebAssembly.instantiate(wasmArrayBuffer, importObject);
    };

    response = await fetchAndInstantiateTask();
  }

  return response;
};

const go = new Go(); // Defined in wasm_exec.js. Don't forget to add this in your index.html.
const WASM_URL = './compiled/naughts-and-crosses.wasm';

export const runWasmAdd = async (ponglehub) => {
  // Get the importObject from the go instance.
  const importObject = go.importObject;

  Object.assign(importObject.env, {
    "main.add": function(a, b) {
      console.info("from wasm!");
      return a + b;
    }
  });

  const wasmModule = await fetch(WASM_URL).then(response => 
    response.arrayBuffer()
  ).then(bytes => 
    WebAssembly.instantiate(bytes, importObject)
  ).catch(err => console.error(err));

  // Allow the wasm_exec go instance, bootstrap and execute our wasm module
  go.run(wasmModule.instance);

  // // Call the Add function export from wasm, save the result
  // const addResult = wasmModule.instance.exports.multiply(24, 24);

  // // Set the result onto the body
  // document.body.textContent = `Hello World! addResult: ${addResult}`;

  return wasmModule.instance.exports
};
