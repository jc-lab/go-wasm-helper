import * as fs from 'node:fs';
import * as path from 'node:path';
import {
  GoError, GoPtr,
  GoWasmHelper,
} from '../src';

import './tinygo_wasm_exec';

async function loadInstance(): Promise<GoWasmHelper> {
  const go = new Go();
  go.argv = process.argv.slice(2);
  go.env = Object.assign({ TMPDIR: require("os").tmpdir() }, process.env);

  const testAppWasmBytes = fs.readFileSync(path.join(__dirname, './testapp.wasm'));
  const result = await WebAssembly.instantiate(testAppWasmBytes, go.importObject);

  go.run(result.instance); // do not wait (it is tinygo's bug...?)

  return new GoWasmHelper(go, result.instance);
}


describe('test app', () => {
  test('sampleData', async () => {
    const helper = await loadInstance();
    const result = helper.callFunction<Uint8Array>('sampleData');
    expect(result).toBeInstanceOf(Uint8Array);
    expect(Buffer.from(result).toString('utf8')).toBe('HELLO WORLD!!!');
  });

  test('sampleError', async () => {
    const helper = await loadInstance();
    expect(() => helper.callFunction<Uint8Array>('sampleError')).toThrow(new GoError('sample error'));
  });

  test('sampleObject', async () => {
    const helper = await loadInstance();
    const result = helper.callFunction<GoPtr>('sampleObject');
    expect(result).toBeGreaterThan(0);
  });

  // test('goroutineTestA', async () => {
  //   const helper = await loadInstance();
  //   console.log(helper.instance.exports)
  //
  //   helper.callFunction<GoPtr>('goroutineTestA');
  //   while (true) {
  //     await new Promise<void>((resolve) => {
  //       setTimeout(() => {
  //         helper.schedule();
  //         resolve();
  //       }, 10);
  //     });
  //   }
  // });
});