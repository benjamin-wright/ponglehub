import './components'
import { runWasmAdd } from './wasm';

global.ponglehub = {
  wasm: async () => {
    if (global.ponglehub._wasm) {
      return global.ponglehub._wasm;
    }

    global.ponglehub._wasm = await runWasmAdd();
    return global.ponglehub._wasm;
  }
}