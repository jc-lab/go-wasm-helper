import * as msgpack from '@msgpack/msgpack';
import {
  RefId,
  RefIsNonPointer,
  RefSubType
} from './refid';

export * from './refid';

export type GoPtr = number;

export class GoError extends Error {
  constructor(message: string, public readonly refId?: RefId) {
    super(message);
  }

  public toString(): string {
    return `GoError: ${this.message}`;
  }
}

type GoPtrAllocateFunc = (size: number) => GoPtr;
type GoPtrFreeFunc = (ptr: GoPtr) => void;

export class GoWasmHelper {
  public instance!: WebAssembly.Instance;
  public memory!: WebAssembly.Memory;
  private fnGoPtrAllocate!: GoPtrAllocateFunc;
  private fnGoPtrFree!: GoPtrFreeFunc;

  private lastCallbackId = 0;
  private readonly callbacks: Record<number, Function> = {};

    constructor(
        public readonly go: Go,
    ) {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const self = this;

        const mem = () => {
            // The buffer may change when requesting more memory.
            return new DataView(this.memory.buffer);
        }

        const gojs = go.importObject.gojs || go.importObject.env;
        const isTinygo = !!gojs['runtime.sleepTicks'];

        const tinygoCallbackJsHandler = function (refIdValue: bigint, args_ptr: number, args_len: number): bigint {
            const args: any[] = [];
            for (let i = 0; i < args_len; i++) {
                args[i] = mem().getBigUint64(args_ptr + 8 * i, true);
            }

            const refId = new RefId(this, BigInt(refIdValue));
            if (!refId.isFunction()) {
                throw new Error(`refid(${BigInt(refIdValue).toString(16)}) is not function`);
            }
            const fn = self.callbacks[Number(refId.getPointer())];
            const out = fn.call(null, ...args);
            if (out instanceof RefId) {
                return out.value;
            }
            return BigInt(out);
        };
        if (isTinygo) {
            gojs['github.com/jc-lab/go-wasm-helper.goCallbackJsHandler'] = tinygoCallbackJsHandler;
            gojs['goCallbackJsHandler'] = tinygoCallbackJsHandler;
        } else {
            const goCallbackJsHandler = function (sp: number) {
                sp >>>= 0;
                const m = mem();
                const refId = m.getBigInt64(sp + 8, true);
                const args_ptr = m.getUint32(sp + 16, true);
                const args_len = m.getUint32(sp + 20, true);
                const result = tinygoCallbackJsHandler(refId, args_ptr, args_len);
                m.setBigUint64(sp + 24, result, true);
            }
            gojs['github.com/jc-lab/go-wasm-helper.goCallbackJsHandler'] = goCallbackJsHandler;
            gojs['goCallbackJsHandler'] = goCallbackJsHandler;
        }
    }

  public run(instance: WebAssembly.Instance): Promise<void> {
    this.instance = instance;

    this.memory = (this.instance.exports.mem || this.instance.exports.memory) as WebAssembly.Memory;
    if (!this.memory) {
      throw new Error('could not find memory in exports');
    }

    this.fnGoPtrAllocate = instance.exports['goPtrAllocate'] as GoPtrAllocateFunc;
    if (!this.fnGoPtrAllocate) {
      throw new Error('could not find goPtrAllocate in exports');
    }

    this.fnGoPtrFree = instance.exports['goPtrFree'] as GoPtrFreeFunc;
    if (!this.fnGoPtrFree) {
      throw new Error('could not find goPtrFree in exports');
    }

    return this.go.run(instance);
  }

  public goPtrAllocate(size: number): GoPtr {
    return this.fnGoPtrAllocate(size);
  }

  public goPtrFree(ptr: GoPtr): void {
    this.fnGoPtrFree(ptr);
  }

  public callFunction(name: string, ...args: any): RefId {
    const fn = this.instance.exports[name] as Function;
    if (!fn) {
      throw new Error(`no function: ${name}`);
    }

    const fixedArgs = args.map((v: any) => {
      if (v instanceof RefId) {
        return v.value;
      } else {
        return v;
      }
    });

    const retval = fn.call(null, ...fixedArgs);
    if (typeof retval !== 'bigint') {
      throw new Error('return value is not bigint');
    }
    const refId = new RefId(this, retval);
    if (refId.isBytes()) {
      refId.loadAndFreeBytes();
    }
    if (refId.isError()) {
      const packedError = msgpack.decode(refId.getBuffer()!) as any;
      throw new GoError(packedError['message'], refId);
    }
    return refId;
  }

  public schedule() {
    const fnScheduler = this.instance.exports['go_scheduler'] as Function;
    if (fnScheduler) {
      fnScheduler();
    }
  }

  public registerCallback(fn: (...args: bigint[]) => RefId): RefId {
    const callbackId = ++this.lastCallbackId;
    this.callbacks[callbackId] = fn;
    return RefId.from(this, callbackId, RefIsNonPointer, RefSubType.RefSubTypeFunction)
  }
}
