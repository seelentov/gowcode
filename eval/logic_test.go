package eval

import "testing"

func TestLogic_Not(t *testing.T) {
	mustBool(t, "not(true)", false)
	mustBool(t, "not(false)", true)
	mustBool(t, "not(0)", true)
	mustEvalErr(t, "not()")
}

func TestLogic_And(t *testing.T) {
	mustBool(t, "and(true, true)", true)
	mustBool(t, "and(true, false)", false)
	mustBool(t, "and(false, true, true)", false)
	mustBool(t, "and(true, true, true)", true)
	mustEvalErr(t, "and(true)")
}

func TestLogic_Or(t *testing.T) {
	mustBool(t, "or(false, false)", false)
	mustBool(t, "or(false, true)", true)
	mustBool(t, "or(true, false, false)", true)
	mustBool(t, "or(false, false, false)", false)
	mustEvalErr(t, "or(true)")
}

func TestLogic_Xor(t *testing.T) {
	mustBool(t, "xor(true, false)", true)
	mustBool(t, "xor(false, true)", true)
	mustBool(t, "xor(true, true)", false)
	mustBool(t, "xor(false, false)", false)
	mustEvalErr(t, "xor(true)")
}

func TestLogic_If(t *testing.T) {
	mustInt(t, "if(true, 1, 2)", 1)
	mustInt(t, "if(false, 1, 2)", 2)
	v := mustEval(t, "if(false, 1)")
	if !v.IsNull() {
		t.Errorf("if(false, 1) should return null, got %v", v.AsString())
	}
	mustInt(t, "if(true, 1)", 1)
	mustEvalErr(t, "if(true)")
}

func TestLogic_Coalesce(t *testing.T) {
	mustInt(t, "coalesce(null, 1)", 1)
	mustInt(t, "coalesce(null, null, 2)", 2)
	mustInt(t, "coalesce(3, 4)", 3)
	v := mustEval(t, "coalesce(null, null)")
	if !v.IsNull() {
		t.Errorf("coalesce(null, null) should be null, got %v", v.AsString())
	}
	mustEvalErr(t, "coalesce()")
}

func TestLogic_DefaultTo(t *testing.T) {
	mustInt(t, "defaultTo(null, 5)", 5)
	mustInt(t, "defaultTo(3, 5)", 3)
	mustEvalErr(t, "defaultTo(1)")
}

func TestLogic_IsNull(t *testing.T) {
	mustBool(t, "isNull(null)", true)
	mustBool(t, "isNull(0)", false)
	mustBool(t, "isNull('')", false)
	mustEvalErr(t, "isNull()")
}

func TestLogic_IsTruthy(t *testing.T) {
	mustBool(t, "isTruthy(1)", true)
	mustBool(t, "isTruthy(0)", false)
	mustBool(t, "isTruthy('hello')", true)
	mustEvalErr(t, "isTruthy()")
}

func TestLogic_IsFalsy(t *testing.T) {
	mustBool(t, "isFalsy(0)", true)
	mustBool(t, "isFalsy(1)", false)
	mustBool(t, "isFalsy('')", true)
	mustEvalErr(t, "isFalsy()")
}
