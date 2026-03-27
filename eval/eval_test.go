package eval

import (
	"gowcode/value"
	"testing"
)

func TestEvaluator_EvaluateSimpleExpression(t *testing.T) {
	res, err := Eval("upper('test')", nil)
	if err != nil {
		t.Fatal(err)
	}

	exp := "TEST"
	if res.AsString() != exp {
		t.Errorf("expected %s, but got %v", exp, res.AsString())
	}
}

func TestEvaluator_EvaluateWithOperator(t *testing.T) {
	res, err := Eval("upper('test') + upper('test')", nil)
	if err != nil {
		t.Fatal(err)
	}

	exp := "TESTTEST"
	if res.AsString() != exp {
		t.Errorf("expected %s, but got %v", exp, res.AsString())
	}
}

func TestEvaluator_EvaluateWithМфкы(t *testing.T) {
	ev := NewEvaluator(map[string]*value.Value{"var": value.StringVal("data")})

	res, err := ev.Evaluate("upper({var}) + '!'")
	if err != nil {
		t.Fatal(err)
	}

	exp := "DATA!"
	if res.AsString() != exp {
		t.Errorf("expected %s, but got %v", exp, res.AsString())
	}
}
