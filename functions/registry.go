package functions

import (
	"fmt"

	"gowcode/value"
)

// Func is a built-in function signature
type Func func(args []*value.Value) (*value.Value, error)

type Registry struct {
	funcs map[string]Func
}

func NewRegistry() *Registry {
	r := &Registry{funcs: make(map[string]Func)}
	registerAll(r)
	return r
}

func (r *Registry) Register(name string, fn Func) {
	r.funcs[name] = fn
}

func (r *Registry) Get(name string) (Func, bool) {
	fn, ok := r.funcs[name]
	return fn, ok
}

func (r *Registry) Call(name string, args []*value.Value) (*value.Value, error) {
	fn, ok := r.funcs[name]
	if !ok {
		return nil, fmt.Errorf("unknown function: %q", name)
	}
	return fn(args)
}

func argsExact(args []*value.Value, n int, name string) error {
	if len(args) != n {
		return fmt.Errorf("%s: expected %d argument(s), got %d", name, n, len(args))
	}
	return nil
}

func argsMin(args []*value.Value, n int, name string) error {
	if len(args) < n {
		return fmt.Errorf("%s: expected at least %d argument(s), got %d", name, n, len(args))
	}
	return nil
}
