package eval

import (
	"math"
	"testing"

	"gowcode/ast"
	"gowcode/parser"
	"gowcode/value"
)

func parseExprHelper(expr string) (ast.Node, error) {
	return parser.Parse(expr)
}

// mustEval evaluates expr and fails the test if there is an error.
func mustEval(t *testing.T, expr string) *value.Value {
	t.Helper()
	v, err := Eval(expr, nil)
	if err != nil {
		t.Fatalf("Eval(%q) unexpected error: %v", expr, err)
	}
	return v
}

// mustEvalErr asserts that evaluating expr returns an error.
func mustEvalErr(t *testing.T, expr string) {
	t.Helper()
	_, err := Eval(expr, nil)
	if err == nil {
		t.Fatalf("Eval(%q): expected error, got nil", expr)
	}
}

// mustInt evaluates expr and asserts the int64 result.
func mustInt(t *testing.T, expr string, want int64) {
	t.Helper()
	v := mustEval(t, expr)
	if got := v.AsInt(); got != want {
		t.Errorf("Eval(%q) = %d, want %d", expr, got, want)
	}
}

// mustFloat evaluates expr and asserts the float64 result with an epsilon.
func mustFloat(t *testing.T, expr string, want float64) {
	t.Helper()
	v := mustEval(t, expr)
	got := v.AsFloat()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Eval(%q) = %v, want %v", expr, got, want)
	}
}

// mustBool evaluates expr and asserts the bool result.
func mustBool(t *testing.T, expr string, want bool) {
	t.Helper()
	v := mustEval(t, expr)
	if got := v.AsBool(); got != want {
		t.Errorf("Eval(%q) = %v, want %v", expr, got, want)
	}
}

// mustStr evaluates expr and asserts the string result.
func mustStr(t *testing.T, expr string, want string) {
	t.Helper()
	v := mustEval(t, expr)
	if got := v.AsString(); got != want {
		t.Errorf("Eval(%q) = %q, want %q", expr, got, want)
	}
}
