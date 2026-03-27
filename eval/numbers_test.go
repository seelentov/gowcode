package eval

import (
	"math"
	"testing"
)

func TestNumbers_Abs(t *testing.T) {
	mustInt(t, "abs(5)", 5)
	mustInt(t, "abs(-5)", 5)
	mustFloat(t, "abs(-3.5)", 3.5)
	mustEvalErr(t, "abs()")
	mustEvalErr(t, "abs(1, 2)")
}

func TestNumbers_Floor(t *testing.T) {
	mustInt(t, "floor(3.9)", 3)
	mustInt(t, "floor(-3.1)", -4)
	mustEvalErr(t, "floor()")
}

func TestNumbers_Ceil(t *testing.T) {
	mustInt(t, "ceil(3.1)", 4)
	mustInt(t, "ceil(-3.9)", -3)
	mustEvalErr(t, "ceil()")
}

func TestNumbers_Round(t *testing.T) {
	mustInt(t, "round(3.5)", 4)
	mustInt(t, "round(3.4)", 3)
	mustFloat(t, "round(3.14159, 2)", 3.14)
	mustFloat(t, "round(3.145, 2)", 3.15)
	mustEvalErr(t, "round()")
}

func TestNumbers_Trunc(t *testing.T) {
	mustInt(t, "trunc(3.9)", 3)
	mustInt(t, "trunc(-3.9)", -3)
	mustEvalErr(t, "trunc()")
}

func TestNumbers_Sqrt(t *testing.T) {
	mustFloat(t, "sqrt(4.0)", 2.0)
	mustFloat(t, "sqrt(9.0)", 3.0)
	mustEvalErr(t, "sqrt()")
}

func TestNumbers_Pow(t *testing.T) {
	mustFloat(t, "pow(2.0, 10.0)", 1024.0)
	mustFloat(t, "pow(3.0, 0.0)", 1.0)
	mustEvalErr(t, "pow(2.0)")
}

func TestNumbers_Sign(t *testing.T) {
	mustInt(t, "sign(-5)", -1)
	mustInt(t, "sign(5)", 1)
	mustInt(t, "sign(0)", 0)
	mustEvalErr(t, "sign()")
}

func TestNumbers_Log(t *testing.T) {
	mustFloat(t, "log(1.0)", 0.0)
	mustFloat(t, "log(8.0, 2.0)", 3.0)
	mustEvalErr(t, "log()")
}

func TestNumbers_Log2(t *testing.T) {
	mustFloat(t, "log2(8.0)", 3.0)
	mustFloat(t, "log2(1.0)", 0.0)
	mustEvalErr(t, "log2()")
}

func TestNumbers_Log10(t *testing.T) {
	mustFloat(t, "log10(1000.0)", 3.0)
	mustFloat(t, "log10(1.0)", 0.0)
	mustEvalErr(t, "log10()")
}

func TestNumbers_PiAndE(t *testing.T) {
	mustFloat(t, "pi()", math.Pi)
	mustFloat(t, "e()", math.E)
}

func TestNumbers_IsNaN(t *testing.T) {
	mustBool(t, "isNaN(0.0)", false)
	mustBool(t, "isNaN(sqrt(-1.0))", true)
	mustEvalErr(t, "isNaN()")
}

func TestNumbers_IsInf(t *testing.T) {
	mustBool(t, "isInf(1.0)", false)
	mustBool(t, "isInf(log(0.0))", true)
	mustEvalErr(t, "isInf()")
}

func TestNumbers_Min(t *testing.T) {
	mustInt(t, "min(3, 1, 2)", 1)
	mustInt(t, "min([3, 1, 2])", 1)
	v := mustEval(t, "min([])")
	if !v.IsNull() {
		t.Errorf("min([]) should be null, got %v", v.AsString())
	}
	mustEvalErr(t, "min()")
}

func TestNumbers_Max(t *testing.T) {
	mustInt(t, "max(1, 3, 2)", 3)
	mustInt(t, "max([1, 3, 2])", 3)
	v := mustEval(t, "max([])")
	if !v.IsNull() {
		t.Errorf("max([]) should be null, got %v", v.AsString())
	}
	mustEvalErr(t, "max()")
}

func TestNumbers_Clamp(t *testing.T) {
	mustInt(t, "clamp(5, 1, 10)", 5)
	mustInt(t, "clamp(-5, 1, 10)", 1)
	mustInt(t, "clamp(15, 1, 10)", 10)
	mustFloat(t, "clamp(5.5, 1.0, 10.0)", 5.5)
	mustFloat(t, "clamp(-1.0, 0.0, 1.0)", 0.0)
	mustFloat(t, "clamp(2.0, 0.0, 1.0)", 1.0)
	mustEvalErr(t, "clamp(1, 2)")
}

func TestNumbers_Sum(t *testing.T) {
	mustInt(t, "sum(1, 2, 3)", 6)
	mustInt(t, "sum([1, 2, 3])", 6)
	mustFloat(t, "sum(1.0, 2.0)", 3.0)
	mustEvalErr(t, "sum()")
}

func TestNumbers_Avg(t *testing.T) {
	mustFloat(t, "avg(2, 4, 6)", 4.0)
	mustFloat(t, "avg([2, 4])", 3.0)
	v := mustEval(t, "avg([])")
	if !v.IsNull() {
		t.Errorf("avg([]) should be null, got %v", v.AsString())
	}
	mustEvalErr(t, "avg()")
}

func TestNumbers_Random(t *testing.T) {
	v := mustEval(t, "random()")
	f := v.AsFloat()
	if f < 0 || f >= 1 {
		t.Errorf("random() = %v, want [0, 1)", f)
	}
}

func TestNumbers_RandomInt(t *testing.T) {
	mustInt(t, "randomInt(5, 5)", 5)
	v := mustEval(t, "randomInt(1, 10)")
	n := v.AsInt()
	if n < 1 || n >= 10 {
		t.Errorf("randomInt(1, 10) = %v, want [1, 10)", n)
	}
	mustEvalErr(t, "randomInt(1)")
}

func TestNumbers_GCD(t *testing.T) {
	mustInt(t, "gcd(12, 8)", 4)
	mustInt(t, "gcd(-12, 8)", 4)
	mustInt(t, "gcd(12, -8)", 4)
	mustInt(t, "gcd(7, 5)", 1)
	mustEvalErr(t, "gcd(1)")
}

func TestNumbers_LCM(t *testing.T) {
	mustInt(t, "lcm(4, 6)", 12)
	mustInt(t, "lcm(0, 5)", 0)
	mustInt(t, "lcm(5, 0)", 0)
	mustInt(t, "lcm(-4, 6)", 12)
	mustInt(t, "lcm(4, -6)", 12)
	mustEvalErr(t, "lcm(1)")
}
