import {GoPtr} from "./index";

export type WasmPtr = number;
export type RefIdValue = bigint;

export const RefVoid: RefIdValue = 0n;

export const RefIsNonPointer: RefIdValue = 0x80000000n;
export const RefIsObject: RefIdValue = 0x40000000n;
export const RefIsError: RefIdValue = 0x20000000n;
export const RefIsBytes: RefIdValue = 0x10000000n;
export const RefLengthOrSubTypeMask: RefIdValue = 0x0fffffffn;

export enum RefSubType {
  RefSubTypeMask = 0x0fff0000,

  RefSubTypeNumber          = 0x0010000,
  RefSubTypeNumberUnsigned = RefSubTypeNumber | 0x1000,
  RefSubTypeInt8            = RefSubTypeNumber | 0x0001,
  RefSubTypeUint8           = RefSubTypeNumberUnsigned | 0x0001,
  RefSubTypeInt16   = RefSubTypeNumber | 0x0002,
  RefSubTypeUint16          = RefSubTypeNumberUnsigned | 0x0002,
  RefSubTypeInt32           = RefSubTypeNumber | 0x0004,
  RefSubTypeUint32          = RefSubTypeNumberUnsigned | 0x0004,

  RefSubTypeFunction = 0x0010001,
}

interface GoWasmHelperLike {
  memory: WebAssembly.Memory;
  goPtrAllocate(size: number): GoPtr;
  goPtrFree(ptr: WasmPtr): void;
}

export class RefId {
  public readonly helper!: GoWasmHelperLike;
  private _bytesBuffer: Uint8Array | null = null;

  constructor(
    helper: GoWasmHelperLike | null,
    public readonly value: RefIdValue,
    bytesBuffer?: Uint8Array,
  ) {
    if (helper) {
      this.helper = helper;
    }
  }

  public static from(
    helper: GoWasmHelperLike | null,
    ptr: WasmPtr,
    flags: RefIdValue,
    lengthOrSubType: number,
    bytesBuffer?: Uint8Array,
  ): RefId {
    return new RefId(helper, (BigInt(ptr) << 32n) | BigInt(flags) | (BigInt(lengthOrSubType) & RefLengthOrSubTypeMask), bytesBuffer);
  }

  public static number(subtype: number, value: number): RefId {
    return RefId.from(null, value, RefIsNonPointer, RefSubType.RefSubTypeNumber | subtype);
  }

  public static int8(value: number): RefId {
    return RefId.number(RefSubType.RefSubTypeInt8, value);
  }

  public static uint8(value: number): RefId {
    return RefId.number(RefSubType.RefSubTypeUint8, value);
  }

  public static int16(value: number): RefId {
    return RefId.number(RefSubType.RefSubTypeInt16, value);
  }

  public static uint16(value: number): RefId {
    return RefId.number(RefSubType.RefSubTypeUint16, value);
  }

  public static int32(value: number): RefId {
    return RefId.number(RefSubType.RefSubTypeInt32, value);
  }

  public static uint32(value: number): RefId {
    return RefId.number(RefSubType.RefSubTypeUint32, value);
  }

  public static Void(): RefId {
    return new RefId(null, RefVoid);
  }

  public static copyBufferFrom(helper: GoWasmHelperLike, data: Uint8Array): RefId {
    const ptr = helper.goPtrAllocate(data.byteLength);
    new Uint8Array(helper.memory.buffer).set(data, ptr);
    return RefId.from(helper, ptr, RefIsBytes, data.byteLength);
  }

  isVoid(): boolean {
    return this.value === RefVoid;
  }

  getPointer(): WasmPtr {
    return Number(this.value >> 32n);
  }

  isNonPointer(): boolean {
    return (this.value & RefIsNonPointer) !== 0n;
  }

  isObject(): boolean {
    return (this.value & RefIsObject) !== 0n;
  }

  isError(): boolean {
    return (this.value & RefIsError) !== 0n;
  }

  isBytes(): boolean {
    return (this.value & RefIsBytes) !== 0n;
  }

  getLength(): number {
    if (!this.isBytes()) {
      return 0;
    }
    return Number(this.value & RefLengthOrSubTypeMask);
  }

  loadAndFreeBytes(): Uint8Array | null {
    if (!this.isBytes()) {
      throw new Error(`refid(0x${this.value.toString(16)}) is not a buffer`);
    }
    if (this._bytesBuffer) {
      return this._bytesBuffer;
    }
    const pointer = this.getPointer();
    if (pointer == 0) {
      return null;
    }
    const buffer = new Uint8Array(this.helper.memory.buffer, pointer, this.getLength());
    const clone = buffer.slice();
    this.helper.goPtrFree(pointer);
    this._bytesBuffer = clone;
    return clone;
  }

  getBuffer(): Uint8Array | null {
    if (this._bytesBuffer) {
      return this._bytesBuffer;
    }
    if (!this.isBytes()) {
      throw new Error(`refid(0x${this.value.toString(16)}) is not a buffer`);
    }
    const pointer = this.getPointer();
    if (pointer == 0) {
      return null;
    }
    return new Uint8Array(this.helper.memory.buffer, pointer, this.getLength());
  }

  getSubType(): number {
    if (!this.isNonPointer()) {
      return 0;
    }
    return Number(this.value & RefLengthOrSubTypeMask);
  }

  isNumber(): boolean {
    return (this.getSubType() & RefSubType.RefSubTypeMask) === RefSubType.RefSubTypeNumber;
  }

  getNumber(): number {
    if (!this.isNumber()) {
      throw new Error(`refid(0x${this.value.toString(16)}) is not a number`);
    }
    return this.getPointer();
  }

  isFunction(): boolean {
    return this.getSubType() === RefSubType.RefSubTypeFunction;
  }

  free(): void {
    if (this.isNonPointer()) {
      throw new Error('is non-pointer');
    }
    this.helper.goPtrFree(this.getPointer());
  }
  //
  // toFunction(): CallbackFunc {
  //   if (!this.isFunction()) {
  //     throw new Error(`refid(0x${this.value.toString(16)}) is not a function`);
  //   }
  //   return (...args: any[]): RefId => {
  //     return this.goCallbackJsHandler(this.value, ...args);
  //   };
  // }
  //
  // private goCallbackJsHandler(refId: RefId, ...args: any[]): RefId {
  //   // Placeholder for the actual Go callback handling logic
  //   return refId;
  // }
}
