package js

import (
	"github.com/abhisekp/napi-go"
)

type Value struct {
	Env   Env
	Value napi.Value
}

func (v Value) GetEnv() Env {
	return v.Env
}
