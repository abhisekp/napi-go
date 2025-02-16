# Contributing

## TODO

### Unstable APIs

```go
// filename: node_api.go

package napi

func BasicFinalize(env Env, finalizeData, finalizeHint unsafe.Pointer) Status {
	return Status(C.node_api_basic_finalize(
		C.napi_env(env),
		finalizeData,
		finalizeHint,
	))
}


func BasicEnv(env Env) Status {
    return Status(C.node_api_basic_env(C.napi_env(env)))
}


func PostFinalizer(env Env, finalizeData, finalizeHint unsafe.Pointer) Status {
    return Status(C.node_api_post_finalizer(
        C.napi_env(env),
        finalizeData,
        finalizeHint,
    ))
}

```

### Fix the following functions

```go
// filename: node_api.go

package napi

/*
#include <stdlib.h>
#include <node/napi_api.h>
*/
import "C"
import "unsafe"

// OpenCallbackScope Function to open a callback scope
func OpenCallbackScope(env Env, resourceObject, context Value) (CallbackScope, Status) {
	var scope CallbackScope
	status := Status(C.napi_open_callback_scope(
		C.napi_env(env),
		C.napi_value(resourceObject),
		C.napi_value(context),
		(*C.napi_callback_scope)(unsafe.Pointer(&scope.scope)),
	))
	return scope, status
}

func CreateExternalStringLatin1(env Env, str string, finalize Finalize, finalizeHint unsafe.Pointer) (Value, Status) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	finalizer := FinalizeToFinalizer(finalize)
	var result Value
	status := Status(C.node_api_create_external_string_latin1(
		C.napi_env(env),
		cstr,
		C.size_t(len([]byte(str))),
		C.napi_finalize(unsafe.Pointer(&finalizer)),
		finalizeHint,
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}


func CreateExternalStringUtf16(env Env, str []uint16, finalize Finalize, finalizeHint unsafe.Pointer) (Value, Status) {
	var result Value
	finalizer := FinalizeToFinalizer(finalize)
	status := Status(C.node_api_create_external_string_utf16(
		C.napi_env(env),
		(*C.char16_t)(unsafe.Pointer(&str[0])),
		C.size_t(len(str)),
		C.napi_finalize(unsafe.Pointer(&finalizer)),
		finalizeHint,
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, status
}


func CreateBufferFromArrayBuffer(env Env, arrayBuffer Value, byteOffset, length int) (Value, *byte, Status) {
	var result Value
	var data *byte
	status := Status(C.node_api_create_buffer_from_arraybuffer(
		C.napi_env(env),
		C.napi_value(arrayBuffer),
		C.size_t(byteOffset),
		C.size_t(length),
		unsafe.Pointer(&data),
		(*C.napi_value)(unsafe.Pointer(&result)),
	))
	return result, data, status
}
```