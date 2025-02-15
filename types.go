package napi

import "unsafe"

type PropertyDescriptor struct {
	Utf8name   string
	Name       Value
	Method     Callback
	Getter     Callback
	Setter     Callback
	Value      Value
	Attributes PropertyAttributes
	Data       unsafe.Pointer
}

type Finalize func(env Env, finalizeData, finalizeHint unsafe.Pointer)

type Reference struct {
	Ref unsafe.Pointer
}

type EscapableHandleScope struct {
	Scope unsafe.Pointer
}

type HandleScope struct {
	Scope unsafe.Pointer
}
