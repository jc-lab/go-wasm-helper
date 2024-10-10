import * as msgpack from '@msgpack/msgpack';

export type GoPtr = number;

const RetIsObject = 0x80000000n;
const RetIsError = 0x40000000n;
const RetIsBytes = 0x20000000n;
const RetLengthMask = 0x0fffffffn;
const RetLengthSign = 0x08000000n;

type GoRefId = bigint;

export class GoError extends Error {
  constructor(message: string) {
    super(message);
  }
}

type GoPtrAllocateFunc = (size: number) => GoPtr;
type GoPtrFreeFunc = (ptr: GoPtr) => void;

export class GoWasmHelper {
  private readonly fnGoPtrAllocate: GoPtrAllocateFunc;
  private readonly fnGoPtrFree: GoPtrFreeFunc;
  private readonly memory: WebAssembly.Memory;

  constructor(
    public readonly go: Go,
    public readonly instance: WebAssembly.Instance
  ) {
    this.fnGoPtrAllocate = instance.exports['goPtrAllocate'] as GoPtrAllocateFunc;
    if (!this.fnGoPtrAllocate) {
      throw new Error('could not find goPtrAllocate in exports');
    }

    this.fnGoPtrFree = instance.exports['goPtrFree'] as GoPtrFreeFunc;
    if (!this.fnGoPtrFree) {
      throw new Error('could not find goPtrFree in exports');
    }

    this.memory = this.instance.exports.memory as WebAssembly.Memory;
    if (!this.memory) {
      throw new Error('could not find memory in exports');
    }
  }

  public goPtrAllocate(size: number): GoPtr {
    return this.fnGoPtrAllocate(size);
  }

  public goPtrFree(ptr: GoPtr): void {
    this.fnGoPtrFree(ptr);
  }

  public callFunction<T extends (Uint8Array | GoPtr | void)>(name: string, ...args: any): T {
     
    const fn = this.instance.exports[name] as Function;
    if (!fn) {
      throw new Error(`no function: ${name}`);
    }

    const retval = fn.call(null, args);
    if (typeof retval !== 'bigint') {
      throw new Error('return value is not bigint');
    }
    const refId: GoRefId = retval;
    const isObject = (refId & RetIsObject) != 0n;
    const isBytes = (refId & RetIsBytes) != 0n;
    const isError = (refId & RetIsError) != 0n;
    const pointer = Number(refId >> 32n); // it is keep unsigned. do not &0xffffffff.
    let lengthBig = refId & RetLengthMask;
    if ((lengthBig & RetLengthSign) != 0n) {
      lengthBig = -(~lengthBig & RetLengthMask) - 1n;
    }
    const length = Number(lengthBig);

    let bytesData: Uint8Array | null = null;
    if (isBytes) {
      const view = new Uint8Array(this.memory.buffer, pointer, length);
      bytesData = view.slice();
      this.goPtrFree(pointer);
    }

    if (isError) {
      if (isBytes) {
        const packedError = msgpack.decode(bytesData!) as any;
        throw new GoError(packedError['message']);
      } else {
        throw new GoError('unknown error');
      }
    }

    if (isBytes) {
      return bytesData! as any;
    } else {
      // isObject or zero
      return pointer as any;
    }
  }

  public schedule() {
    const fnScheduler = this.instance.exports['go_scheduler'] as Function;
    if (fnScheduler) {
      fnScheduler();
    }
  }
}
