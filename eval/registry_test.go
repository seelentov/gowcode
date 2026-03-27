package eval

import (
	"testing"

	"gowcode/functions"
	"gowcode/value"
)

func TestRegistry_UnknownFunction(t *testing.T) {
	mustEvalErr(t, "noSuchFunction(1, 2)")
}

func TestRegistry_RegisterFuncAndGet(t *testing.T) {
	reg := functions.NewRegistry()

	// RegisterFunc — public alias for Register
	reg.RegisterFunc("triple", func(args []*value.Value) (*value.Value, error) {
		if len(args) != 1 {
			return nil, nil
		}
		return value.IntVal(args[0].AsInt() * 3), nil
	})

	// Get — retrieve a registered function by name
	fn, ok := reg.Get("triple")
	if !ok {
		t.Fatal("Get('triple') returned false, want true")
	}
	result, err := fn([]*value.Value{value.IntVal(7)})
	if err != nil {
		t.Fatalf("triple(7) returned error: %v", err)
	}
	if result.AsInt() != 21 {
		t.Errorf("triple(7) = %d, want 21", result.AsInt())
	}

	// Get — missing function
	_, ok = reg.Get("nonexistent")
	if ok {
		t.Error("Get('nonexistent') returned true, want false")
	}

	// Call — unknown function error path
	_, err = reg.Call("nonexistent", nil)
	if err == nil {
		t.Error("Call('nonexistent') should return an error")
	}

	// Call — known function through registry
	out, err := reg.Call("triple", []*value.Value{value.IntVal(3)})
	if err != nil || out.AsInt() != 9 {
		t.Errorf("Call('triple', [3]) = %v, %v; want 9, nil", out, err)
	}
}

func TestRegistry_NewEvaluatorWithRegistry(t *testing.T) {
	reg := functions.NewRegistry()
	reg.RegisterFunc("double", func(args []*value.Value) (*value.Value, error) {
		return value.IntVal(args[0].AsInt() * 2), nil
	})
	ev := NewEvaluatorWithRegistry(nil, reg)
	node, err := parseExprHelper("double(5)")
	if err != nil {
		t.Fatal(err)
	}
	v, err := ev.eval(node)
	if err != nil {
		t.Fatal(err)
	}
	if v.AsInt() != 10 {
		t.Errorf("double(5) = %d, want 10", v.AsInt())
	}
}
