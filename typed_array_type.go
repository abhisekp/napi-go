package napi

type TypedArrayType int

const (
	TypedArrayInt8Array TypedArrayType = iota
	TypedArrayUint8Array
	TypedArrayUint8ClampedArray
	TypedArrayInt16Array
	TypedArrayUint16Array
	TypedArrayInt32Array
	TypedArrayUint32Array
	TypedArrayFloat32Array
	TypedArrayFloat64Array
	TypedArrayBigInt64Array
	TypedArrayBigUint64Array
)
