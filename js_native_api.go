package napi

/*
#include <stdlib.h>
#include <node/node_api.h>
*/
import "C"

import (
	"unsafe"
)

func GetUndefined(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_get_undefined(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetNull(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_get_null(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetGlobal(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_get_global(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetBoolean(env Env, value bool) (Value, Status) {
	var result Value
	status := Status(C.napi_get_boolean(
		C.napi_env(env),
		C.bool(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateObject(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_create_object(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateArray(env Env) (Value, Status) {
	var result Value
	status := Status(C.napi_create_array(
		C.napi_env(env),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateArrayWithLength(env Env, length int) (Value, Status) {
	var result Value
	status := Status(C.napi_create_array_with_length(
		C.napi_env(env),
		C.size_t(length),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateDouble(env Env, value float64) (Value, Status) {
	var result Value
	status := Status(C.napi_create_double(
		C.napi_env(env),
		C.double(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateStringUtf8(env Env, str string) (Value, Status) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	var result Value
	status := Status(C.napi_create_string_utf8(
		C.napi_env(env),
		cstr,
		C.size_t(len([]byte(str))), // must pass number of bytes
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateSymbol(env Env, description Value) (Value, Status) {
	var result Value
	status := Status(C.napi_create_symbol(
		C.napi_env(env),
		C.napi_value(description),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateFunction(env Env, name string, cb Callback) (Value, Status) {
	provider, status := getInstanceData(env)
	if status != StatusOK || provider == nil {
		return nil, status
	}

	return provider.GetCallbackData().CreateCallback(env, name, cb)
}

func CreateError(env Env, code, msg Value) (Value, Status) {
	var result Value
	status := Status(C.napi_create_error(
		C.napi_env(env),
		C.napi_value(code),
		C.napi_value(msg),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func Typeof(env Env, value Value) (ValueType, Status) {
	var result ValueType
	status := Status(C.napi_typeof(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_valuetype)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueDouble(env Env, value Value) (float64, Status) {
	var result float64
	status := Status(C.napi_get_value_double(
		C.napi_env(env),
		C.napi_value(value),
		(*C.double)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueBool(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_get_value_bool(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueStringUtf8(env Env, value Value) (string, Status) {
	// call napi_get_value_string_utf8 twice
	// first is to get number of bytes
	// second is to populate the actual string buffer
	bufsize := C.size_t(0)
	var strsize C.size_t

	status := Status(C.napi_get_value_string_utf8(
		C.napi_env(env),
		C.napi_value(value),
		nil,
		bufsize,
		&strsize,
	))

	if status != StatusOK {
		return "", status
	}

	// ensure there is room for the null terminator as well
	strsize++
	cstr := (*C.char)(C.malloc(C.sizeof_char * strsize))
	defer C.free(unsafe.Pointer(cstr))

	status = Status(C.napi_get_value_string_utf8(
		C.napi_env(env),
		C.napi_value(value),
		cstr,
		strsize,
		&strsize,
	))

	if status != StatusOK {
		return "", status
	}

	return C.GoStringN(
		(*C.char)(cstr),
		(C.int)(strsize),
	), status
}

func SetProperty(env Env, object, key, value Value) Status {
	return Status(C.napi_set_property(
		C.napi_env(env),
		C.napi_value(object),
		C.napi_value(key),
		C.napi_value(value),
	))
}

func SetElement(env Env, object Value, index int, value Value) Status {
	return Status(C.napi_set_element(
		C.napi_env(env),
		C.napi_value(object),
		C.uint32_t(index),
		C.napi_value(value),
	))
}

func StrictEquals(env Env, lhs, rhs Value) (bool, Status) {
	var result bool
	status := Status(C.napi_strict_equals(
		C.napi_env(env),
		C.napi_value(lhs),
		C.napi_value(rhs),
		(*C.bool)(&result),
	))
	return result, status
}

type GetCbInfoResult struct {
	Args []Value
	This Value
}

func GetCbInfo(env Env, info CallbackInfo) (GetCbInfoResult, Status) {
	// call napi_get_cb_info twice
	// first is to get total number of arguments
	// second is to populate the actual arguments
	argc := C.size_t(0)
	status := Status(C.napi_get_cb_info(
		C.napi_env(env),
		C.napi_callback_info(info),
		&argc,
		nil,
		nil,
		nil,
	))

	if status != StatusOK {
		return GetCbInfoResult{}, status
	}

	argv := make([]Value, int(argc))
	var cArgv unsafe.Pointer
	if argc > 0 {
		cArgv = unsafe.Pointer(&argv[0]) // must pass element pointer
	}

	var thisArg Value

	status = Status(C.napi_get_cb_info(
		C.napi_env(env),
		C.napi_callback_info(info),
		&argc,
		(*C.napi_value)(cArgv),
		(*C.napi_value)(unsafe.Pointer(&thisArg)),
		nil,
	))

	return GetCbInfoResult{
		Args: argv,
		This: thisArg,
	}, status
}

func Throw(env Env, err Value) Status {
	return Status(C.napi_throw(
		C.napi_env(env),
		C.napi_value(err),
	))
}

func ThrowError(env Env, code, msg string) Status {
	codeCStr, msgCCstr := C.CString(code), C.CString(msg)
	defer C.free(unsafe.Pointer(codeCStr))
	defer C.free(unsafe.Pointer(msgCCstr))

	return Status(C.napi_throw_error(
		C.napi_env(env),
		codeCStr,
		msgCCstr,
	))
}

func CreatePromise(env Env) (Promise, Status) {
	var result Promise
	status := Status(C.napi_create_promise(
		C.napi_env(env),
		(*C.napi_deferred)(unsafe.Pointer(&result.Deferred)),
		(*C.napi_value)(unsafe.Pointer(&result.Value)),
	))
	return result, status
}

func ResolveDeferred(env Env, deferred Deferred, resolution Value) Status {
	return Status(C.napi_resolve_deferred(
		C.napi_env(env),
		C.napi_deferred(deferred),
		C.napi_value(resolution),
	))
}

func RejectDeferred(env Env, deferred Deferred, rejection Value) Status {
	return Status(C.napi_reject_deferred(
		C.napi_env(env),
		C.napi_deferred(deferred),
		C.napi_value(rejection),
	))
}

func SetInstanceData(env Env, data any) Status {
	provider, status := getInstanceData(env)
	if status != StatusOK || provider == nil {
		return status
	}

	provider.SetUserData(data)
	return status
}

func GetInstanceData(env Env) (any, Status) {
	provider, status := getInstanceData(env)
	if status != StatusOK || provider == nil {
		return nil, status
	}

	return provider.GetUserData(), status
}

func CreateBuffer(env Env, length int) (Value, *byte, Status) {
	var result Value
	var data *byte
	status := Status(C.napi_create_buffer(
		C.napi_env(env),
		C.size_t(length),
		(**C.void)(unsafe.Pointer(&data)),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, data, status
}

func CreateExternal(env Env, data unsafe.Pointer, finalize Finalize, finalizeHint unsafe.Pointer) (Value, Status) {
	var result Value
	status := Status(C.napi_create_external(
		C.napi_env(env),
		data,
		C.napi_finalize(finalize),
		finalizeHint,
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueInt32(env Env, value Value) (int32, Status) {
	var result int32
	status := Status(C.napi_get_value_int32(
		C.napi_env(env),
		C.napi_value(value),
		(*C.int32_t)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueUint32(env Env, value Value) (uint32, Status) {
	var result uint32
	status := Status(C.napi_get_value_uint32(
		C.napi_env(env),
		C.napi_value(value),
		(*C.uint32_t)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueInt64(env Env, value Value) (int64, Status) {
	var result int64
	status := Status(C.napi_get_value_int64(
		C.napi_env(env),
		C.napi_value(value),
		(*C.int64_t)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueBigIntInt64(env Env, value Value) (int64, bool, Status) {
	var result int64
	var lossless bool
	status := Status(C.napi_get_value_bigint_int64(
		C.napi_env(env),
		C.napi_value(value),
		(*C.int64_t)(unsafe.Pointer(&result)),
		(*C.bool)(unsafe.Pointer(&lossless)),
	))
	return result, lossless, status
}

func GetValueBigIntWords(env Env, value Value, signBit int, wordCount int, words *uint64) Status {
	return Status(C.napi_get_value_bigint_words(
		C.napi_env(env),
		C.napi_value(value),
		(*C.int)(unsafe.Pointer(&signBit)),
		(*C.size_t)(unsafe.Pointer(&wordCount)),
		(*C.uint64_t)(unsafe.Pointer(words)),
	))
}

func GetValueExternal(env Env, value Value) (unsafe.Pointer, Status) {
	var result unsafe.Pointer
	status := Status(C.napi_get_value_external(
		C.napi_env(env),
		C.napi_value(value),
		&result,
	))
	return result, status
}

func CoerceToBool(env Env, value Value) (Value, Status) {
	var result Value
	status := Status(C.napi_coerce_to_bool(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CoerceToNumber(env Env, value Value) (Value, Status) {
	var result Value
	status := Status(C.napi_coerce_to_number(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CoerceToObject(env Env, value Value) (Value, Status) {
	var result Value
	status := Status(C.napi_coerce_to_object(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CoerceToString(env Env, value Value) (Value, Status) {
	var result Value
	status := Status(C.napi_coerce_to_string(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateBufferCopy(env Env, data []byte) (Value, *byte, Status) {
	var result Value
	var copiedData *byte
	status := Status(C.napi_create_buffer_copy(
		C.napi_env(env),
		C.size_t(len(data)),
		unsafe.Pointer(&data[0]),
		(**C.void)(unsafe.Pointer(&copiedData)),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, copiedData, status
}

func GetBufferInfo(env Env, value Value) (*byte, int, Status) {
	var data *byte
	var length C.size_t
	status := Status(C.napi_get_buffer_info(
		C.napi_env(env),
		C.napi_value(value),
		(**C.void)(unsafe.Pointer(&data)),
		&length,
	))
	return data, int(length), status
}

func GetArrayLength(env Env, value Value) (int, Status) {
	var length C.uint32_t
	status := Status(C.napi_get_array_length(
		C.napi_env(env),
		C.napi_value(value),
		&length,
	))
	return int(length), status
}

func GetPrototype(env Env, value Value) (Value, Status) {
	var result Value
	status := Status(C.napi_get_prototype(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func InstanceOf(env Env, object, constructor Value) (bool, Status) {
	var result bool
	status := Status(C.napi_instanceof(
		C.napi_env(env),
		C.napi_value(object),
		C.napi_value(constructor),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func IsArray(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_is_array(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func IsBuffer(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_is_buffer(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func IsError(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_is_error(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func IsPromise(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_is_promise(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func IsTypedArray(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_is_typedarray(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetTypedArrayInfo(env Env, value Value) (TypedArrayType, int, *byte, Status) {
	var type_ TypedArrayType
	var length C.size_t
	var data *byte
	status := Status(C.napi_get_typedarray_info(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_typedarray_type)(unsafe.Pointer(&type_)),
		&length,
		(**C.void)(unsafe.Pointer(&data)),
	))
	return type_, int(length), data, status
}

func CreateTypedArray(env Env, type_ TypedArrayType, length int, arrayBuffer Value, byteOffset int) (Value, Status) {
	var result Value
	status := Status(C.napi_create_typedarray(
		C.napi_env(env),
		C.napi_typedarray_type(type_),
		C.size_t(length),
		C.napi_value(arrayBuffer),
		C.size_t(byteOffset),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func AdjustExternalMemory(env Env, change int64) (int64, Status) {
	var result int64
	status := Status(C.napi_adjust_external_memory(
		C.napi_env(env),
		C.int64_t(change),
		(*C.int64_t)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateDataView(env Env, length int, arrayBuffer Value, byteOffset int) (Value, Status) {
	var result Value
	status := Status(C.napi_create_dataview(
		C.napi_env(env),
		C.size_t(length),
		C.napi_value(arrayBuffer),
		C.size_t(byteOffset),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetDataViewInfo(env Env, value Value) (int, *byte, Status) {
	var length C.size_t
	var data *byte
	status := Status(C.napi_get_dataview_info(
		C.napi_env(env),
		C.napi_value(value),
		&length,
		(**C.void)(unsafe.Pointer(&data)),
	))
	return int(length), data, status
}

func GetAllPropertyNames(env Env, object Value) (Value, Status) {
	var result Value
	status := Status(C.napi_get_all_property_names(
		C.napi_env(env),
		C.napi_value(object),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func HasOwnProperty(env Env, object, key Value) (bool, Status) {
	var result bool
	status := Status(C.napi_has_own_property(
		C.napi_env(env),
		C.napi_value(object),
		C.napi_value(key),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func HasProperty(env Env, object, key Value) (bool, Status) {
	var result bool
	status := Status(C.napi_has_property(
		C.napi_env(env),
		C.napi_value(object),
		C.napi_value(key),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetPropertyNames(env Env, object Value) (Value, Status) {
	var result Value
	status := Status(C.napi_get_property_names(
		C.napi_env(env),
		C.napi_value(object),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func DefineProperties(env Env, object Value, properties []PropertyDescriptor) Status {
	return Status(C.napi_define_properties(
		C.napi_env(env),
		C.napi_value(object),
		C.size_t(len(properties)),
		(*C.napi_property_descriptor)(unsafe.Pointer(&properties[0])),
	))
}

func Wrap(env Env, jsObject Value, nativeObject unsafe.Pointer, finalize Finalize, finalizeHint unsafe.Pointer) Status {
	return Status(C.napi_wrap(
		C.napi_env(env),
		C.napi_value(jsObject),
		nativeObject,
		C.napi_finalize(finalize),
		finalizeHint,
		(*C.napi_value)(unsafe.Pointer(&jsObject)),
	))
}

func Unwrap(env Env, jsObject Value) (unsafe.Pointer, Status) {
	var nativeObject unsafe.Pointer
	status := Status(C.napi_unwrap(
		C.napi_env(env),
		C.napi_value(jsObject),
		&nativeObject,
	))
	return nativeObject, status
}

func RemoveWrap(env Env, jsObject Value) Status {
	return Status(C.napi_remove_wrap(
		C.napi_env(env),
		C.napi_value(jsObject),
	))
}

func OpenHandleScope(env Env) (HandleScope, Status) {
	var scope HandleScope
	status := Status(C.napi_open_handle_scope(
		C.napi_env(env),
		(*C.napi_handle_scope)(unsafe.Pointer(&scope)),
	))
	return scope, status
}

func CloseHandleScope(env Env, scope HandleScope) Status {
	return Status(C.napi_close_handle_scope(
		C.napi_env(env),
		C.napi_handle_scope(scope),
	))
}

func OpenEscapableHandleScope(env Env) (EscapableHandleScope, Status) {
	var scope EscapableHandleScope
	status := Status(C.napi_open_escapable_handle_scope(
		C.napi_env(env),
		(*C.napi_escapable_handle_scope)(unsafe.Pointer(&scope)),
	))
	return scope, status
}

func CloseEscapableHandleScope(env Env, scope EscapableHandleScope) Status {
	return Status(C.napi_close_escapable_handle_scope(
		C.napi_env(env),
		C.napi_escapable_handle_scope(scope),
	))
}

func EscapeHandle(env Env, scope EscapableHandleScope, escapee Value) (Value, Status) {
	var result Value
	status := Status(C.napi_escape_handle(
		C.napi_env(env),
		C.napi_escapable_handle_scope(scope),
		C.napi_value(escapee),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateReference(env Env, value Value, initialRefcount int) (Reference, Status) {
	var ref Reference
	status := Status(C.napi_create_reference(
		C.napi_env(env),
		C.napi_value(value),
		C.uint32_t(initialRefcount),
		(*C.napi_ref)(unsafe.Pointer(&ref)),
	))
	return ref, status
}

func DeleteReference(env Env, ref Reference) Status {
	return Status(C.napi_delete_reference(
		C.napi_env(env),
		C.napi_ref(ref),
	))
}

func ReferenceRef(env Env, ref Reference) (int, Status) {
	var result C.uint32_t
	status := Status(C.napi_reference_ref(
		C.napi_env(env),
		C.napi_ref(ref),
		&result,
	))
	return int(result), status
}

func ReferenceUnref(env Env, ref Reference) (int, Status) {
	var result C.uint32_t
	status := Status(C.napi_reference_unref(
		C.napi_env(env),
		C.napi_ref(ref),
		&result,
	))
	return int(result), status
}

func GetReferenceValue(env Env, ref Reference) (Value, Status) {
	var result Value
	status := Status(C.napi_get_reference_value(
		C.napi_env(env),
		C.napi_ref(ref),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetValueBigIntUint64(env Env, value Value) (uint64, bool, Status) {
	var result uint64
	var lossless bool
	status := Status(C.napi_get_value_bigint_uint64(
		C.napi_env(env),
		C.napi_value(value),
		(*C.uint64_t)(unsafe.Pointer(&result)),
		(*C.bool)(unsafe.Pointer(&lossless)),
	))
	return result, lossless, status
}

func CreateBigIntInt64(env Env, value int64) (Value, Status) {
	var result Value
	status := Status(C.napi_create_bigint_int64(
		C.napi_env(env),
		C.int64_t(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateBigIntUint64(env Env, value uint64) (Value, Status) {
	var result Value
	status := Status(C.napi_create_bigint_uint64(
		C.napi_env(env),
		C.uint64_t(value),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func CreateBigIntWords(env Env, signBit int, wordCount int, words *uint64) (Value, Status) {
	var result Value
	status := Status(C.napi_create_bigint_words(
		C.napi_env(env),
		C.int(signBit),
		C.size_t(wordCount),
		(*C.uint64_t)(unsafe.Pointer(words)),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func IsDate(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_is_date(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func IsDetachedArrayBuffer(env Env, value Value) (bool, Status) {
	var result bool
	status := Status(C.napi_is_detached_arraybuffer(
		C.napi_env(env),
		C.napi_value(value),
		(*C.bool)(unsafe.Pointer(&result)),
	))
	return result, status
}

func DetachArrayBuffer(env Env, value Value) Status {
	return Status(C.napi_detach_arraybuffer(
		C.napi_env(env),
		C.napi_value(value),
	))
}

func CreateArrayBuffer(env Env, length int) (Value, *byte, Status) {
	var result Value
	var data *byte
	status := Status(C.napi_create_arraybuffer(
		C.napi_env(env),
		C.size_t(length),
		(**C.void)(unsafe.Pointer(&data)),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, data, status
}

func GetArrayBufferInfo(env Env, value Value) (*byte, int, Status) {
	var data *byte
	var length C.size_t
	status := Status(C.napi_get_arraybuffer_info(
		C.napi_env(env),
		C.napi_value(value),
		(**C.void)(unsafe.Pointer(&data)),
		&length,
	))
	return data, int(length), status
}

func CreateExternalArrayBuffer(env Env, data unsafe.Pointer, length int, finalize Finalize, finalizeHint unsafe.Pointer) (Value, Status) {
	var result Value
	status := Status(C.napi_create_external_arraybuffer(
		C.napi_env(env),
		data,
		C.size_t(length),
		C.napi_finalize(finalize),
		finalizeHint,
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}

func GetTypedArrayInfo(env Env, value Value) (TypedArrayType, int, *byte, Status) {
	var type_ TypedArrayType
	var length C.size_t
	var data *byte
	status := Status(C.napi_get_typedarray_info(
		C.napi_env(env),
		C.napi_value(value),
		(*C.napi_typedarray_type)(unsafe.Pointer(&type_)),
		&length,
		(**C.void)(unsafe.Pointer(&data)),
	))
	return type_, int(length), data, status
}

func CreateTypedArray(env Env, type_ TypedArrayType, length int, arrayBuffer Value, byteOffset int) (Value, Status) {
	var result Value
	status := Status(C.napi_create_typedarray(
		C.napi_env(env),
		C.napi_typedarray_type(type_),
		C.size_t(length),
		C.napi_value(arrayBuffer),
		C.size_t(byteOffset),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}
