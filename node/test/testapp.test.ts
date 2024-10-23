import * as fs from 'node:fs';
import * as path from 'node:path';
import {
  GoError, GoPtr,
  GoWasmHelper,
} from '../src';

import './tinygo_wasm_exec';
import {RefId} from "../src/refid";

async function loadInstance(): Promise<GoWasmHelper> {
  const go = new Go();
  go.argv = process.argv.slice(2);
  go.env = Object.assign({ TMPDIR: require("os").tmpdir() }, process.env);

  const helper = new GoWasmHelper(go);

  const testAppWasmBytes = fs.readFileSync(path.join(__dirname, './testapp.wasm'));
  const result = await WebAssembly.instantiate(testAppWasmBytes, go.importObject);

  helper.run(result.instance); // do not wait (it is tinygo's bug...?)

  return helper;
}


describe('test app', () => {
  test('sampleData', async () => {
    const helper = await loadInstance();
    const result = helper.callFunction('sampleData');
    expect(result.isBytes()).toBe(true);
    expect(Buffer.from(result.getBuffer()!!).toString('utf8')).toBe('HELLO WORLD!!!');
  });

  test('sampleError', async () => {
    const helper = await loadInstance();
    expect(() => helper.callFunction('sampleError')).toThrow(new GoError('sample error'));
  });

  test('sampleObject', async () => {
    const helper = await loadInstance();
    const result = helper.callFunction('sampleObject');
    expect(result.value).toBeGreaterThan(0);
    expect(result.isObject()).toBe(true);
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

  test('callJsTest', async () => {
    const helper = await loadInstance();
    const jsCallbackRefId = helper.registerCallback((a: bigint): RefId => {
      return RefId.uint32(Number(a) + 0x100);
    });
    const result = helper.callFunction('callbackTest', jsCallbackRefId, 0x12);
    expect(result.isNumber()).toBe(true);
    expect(result.getNumber()).toBe(0x12 + 0x100 + 0x1000 + 0x2000);
  });

  test('copyBufferFrom', async () => {
    const sampleData = Buffer.from('HELLO WORLD');

    const helper = await loadInstance();
    const param = RefId.copyBufferFrom(helper, sampleData);
    const result = helper.callFunction('bufferReadTest', param);
    param.free();

    expect(Buffer.compare(result.loadAndFreeBytes() as any, sampleData)).toBe(0);
  })
});