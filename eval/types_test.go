package eval

import "testing"

func TestTypes_ToString(t *testing.T) {
	mustStr(t, "toString(42)", "42")
	mustStr(t, "toString(3.14)", "3.14")
	mustStr(t, "toString(true)", "true")
	mustStr(t, "toString(null)", "null")
	mustEvalErr(t, "toString()")
}

func TestTypes_ToInt(t *testing.T) {
	mustInt(t, "toInt('42')", 42)
	mustInt(t, "toInt(3.9)", 3)
	mustInt(t, "toInt(true)", 1)
	mustInt(t, "toInt(false)", 0)
	mustEvalErr(t, "toInt()")
}

func TestTypes_ToFloat(t *testing.T) {
	mustFloat(t, "toFloat('3.14')", 3.14)
	mustFloat(t, "toFloat(5)", 5.0)
	mustEvalErr(t, "toFloat()")
}

func TestTypes_ToBool(t *testing.T) {
	mustBool(t, "toBool(1)", true)
	mustBool(t, "toBool(0)", false)
	mustBool(t, "toBool('hello')", true)
	mustBool(t, "toBool('')", false)
	mustEvalErr(t, "toBool()")
}

func TestTypes_ToList(t *testing.T) {
	// already a list → returned as-is
	v := mustEval(t, "toList([1, 2])")
	if v.Len() != 2 {
		t.Errorf("toList([1,2]) length = %d, want 2", v.Len())
	}
	// null → empty list
	v = mustEval(t, "toList(null)")
	if v.Len() != 0 {
		t.Errorf("toList(null) length = %d, want 0", v.Len())
	}
	// scalar → wrapped in list
	v = mustEval(t, "toList(42)")
	if v.Len() != 1 {
		t.Errorf("toList(42) length = %d, want 1", v.Len())
	}
	mustStr(t, "toList(42)[0]", "42")
	mustEvalErr(t, "toList()")
}

func TestTypes_TypeOf(t *testing.T) {
	mustStr(t, "typeOf(null)", "null")
	mustStr(t, "typeOf(true)", "bool")
	mustStr(t, "typeOf(1)", "int")
	mustStr(t, "typeOf(1.0)", "float")
	mustStr(t, "typeOf('hello')", "string")
	mustStr(t, "typeOf([1, 2])", "list")
	mustStr(t, "typeOf({'a': 1})", "map")
	mustEvalErr(t, "typeOf()")
}

func TestTypes_IsString(t *testing.T) {
	mustBool(t, "isString('hello')", true)
	mustBool(t, "isString(1)", false)
	mustEvalErr(t, "isString()")
}

func TestTypes_IsInt(t *testing.T) {
	mustBool(t, "isInt(1)", true)
	mustBool(t, "isInt(1.0)", false)
	mustEvalErr(t, "isInt()")
}

func TestTypes_IsFloat(t *testing.T) {
	mustBool(t, "isFloat(1.0)", true)
	mustBool(t, "isFloat(1)", false)
	mustEvalErr(t, "isFloat()")
}

func TestTypes_IsBool(t *testing.T) {
	mustBool(t, "isBool(true)", true)
	mustBool(t, "isBool(1)", false)
	mustEvalErr(t, "isBool()")
}

func TestTypes_IsList(t *testing.T) {
	mustBool(t, "isList([1, 2])", true)
	mustBool(t, "isList(1)", false)
	mustEvalErr(t, "isList()")
}

func TestTypes_IsMap(t *testing.T) {
	mustBool(t, "isMap({'a': 1})", true)
	mustBool(t, "isMap(1)", false)
	mustEvalErr(t, "isMap()")
}

func TestTypes_IsNumber(t *testing.T) {
	mustBool(t, "isNumber(1)", true)
	mustBool(t, "isNumber(1.0)", true)
	mustBool(t, "isNumber('3')", false)
	mustEvalErr(t, "isNumber()")
}

func TestTypes_IsNull(t *testing.T) {
	mustBool(t, "isNull(null)", true)
	mustBool(t, "isNull(0)", false)
	mustBool(t, "isNull('')", false)
	mustEvalErr(t, "isNull()")
}
