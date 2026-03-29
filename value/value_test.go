package value

import (
	"math"
	"testing"
)

// --- Constructors and type checks ---

func TestConstructors_Nil(t *testing.T) {
	v := Nil()
	if !v.IsNull() {
		t.Error("Nil() should be null")
	}
	if v.IsBool() || v.IsInt() || v.IsFloat() || v.IsString() || v.IsList() || v.IsMap() {
		t.Error("Nil() should not match other type checks")
	}
}

func TestConstructors_Bool(t *testing.T) {
	vt := BoolVal(true)
	vf := BoolVal(false)
	if !vt.IsBool() || !vf.IsBool() {
		t.Error("BoolVal should be bool")
	}
	if vt.IsNull() || vt.IsInt() {
		t.Error("BoolVal(true) should not be null or int")
	}
}

func TestConstructors_Int(t *testing.T) {
	v := IntVal(42)
	if !v.IsInt() {
		t.Error("IntVal should be int")
	}
	if !v.IsNumber() {
		t.Error("IntVal should be number")
	}
	if v.IsFloat() {
		t.Error("IntVal should not be float")
	}
}

func TestConstructors_Float(t *testing.T) {
	v := FloatVal(3.14)
	if !v.IsFloat() {
		t.Error("FloatVal should be float")
	}
	if !v.IsNumber() {
		t.Error("FloatVal should be number")
	}
	if v.IsInt() {
		t.Error("FloatVal should not be int")
	}
}

func TestConstructors_String(t *testing.T) {
	v := StringVal("hello")
	if !v.IsString() {
		t.Error("StringVal should be string")
	}
}

func TestConstructors_List(t *testing.T) {
	v := ListVal(IntVal(1), IntVal(2))
	if !v.IsList() {
		t.Error("ListVal should be list")
	}
}

func TestConstructors_ListEmpty(t *testing.T) {
	v := ListVal()
	if !v.IsList() {
		t.Error("ListVal() should be list")
	}
	if v.Len() != 0 {
		t.Error("ListVal() should be empty")
	}
}

func TestConstructors_Map(t *testing.T) {
	v := MapVal()
	if !v.IsMap() {
		t.Error("MapVal should be map")
	}
}

// --- AsBool ---

func TestAsBool_Bool(t *testing.T) {
	if !BoolVal(true).AsBool() {
		t.Error("true should be truthy")
	}
	if BoolVal(false).AsBool() {
		t.Error("false should be falsy")
	}
}

func TestAsBool_Int(t *testing.T) {
	if !IntVal(1).AsBool() {
		t.Error("1 should be truthy")
	}
	if IntVal(0).AsBool() {
		t.Error("0 should be falsy")
	}
	if !IntVal(-1).AsBool() {
		t.Error("-1 should be truthy")
	}
}

func TestAsBool_Float(t *testing.T) {
	if !FloatVal(1.5).AsBool() {
		t.Error("1.5 should be truthy")
	}
	if FloatVal(0.0).AsBool() {
		t.Error("0.0 should be falsy")
	}
}

func TestAsBool_String(t *testing.T) {
	if !StringVal("hello").AsBool() {
		t.Error("non-empty string should be truthy")
	}
	if StringVal("").AsBool() {
		t.Error("empty string should be falsy")
	}
	if StringVal("false").AsBool() {
		t.Error(`"false" should be falsy`)
	}
	if StringVal("0").AsBool() {
		t.Error(`"0" should be falsy`)
	}
	if !StringVal("true").AsBool() {
		t.Error(`"true" should be truthy`)
	}
}

func TestAsBool_Null(t *testing.T) {
	if Nil().AsBool() {
		t.Error("null should be falsy")
	}
}

func TestAsBool_List(t *testing.T) {
	if !ListVal(IntVal(1)).AsBool() {
		t.Error("non-empty list should be truthy")
	}
	if ListVal().AsBool() {
		t.Error("empty list should be falsy")
	}
}

func TestAsBool_Map(t *testing.T) {
	m := MapVal()
	if m.AsBool() {
		t.Error("empty map should be falsy")
	}
	m.Set("k", IntVal(1))
	if !m.AsBool() {
		t.Error("non-empty map should be truthy")
	}
}

// --- AsInt ---

func TestAsInt_Int(t *testing.T) {
	if IntVal(99).AsInt() != 99 {
		t.Error("AsInt on int")
	}
}

func TestAsInt_Float(t *testing.T) {
	if FloatVal(3.9).AsInt() != 3 {
		t.Error("AsInt on float should truncate")
	}
}

func TestAsInt_Bool(t *testing.T) {
	if BoolVal(true).AsInt() != 1 {
		t.Error("true.AsInt() should be 1")
	}
	if BoolVal(false).AsInt() != 0 {
		t.Error("false.AsInt() should be 0")
	}
}

func TestAsInt_String(t *testing.T) {
	if StringVal("42").AsInt() != 42 {
		t.Error("'42'.AsInt() should be 42")
	}
	if StringVal("abc").AsInt() != 0 {
		t.Error("non-numeric string AsInt should be 0")
	}
}

func TestAsInt_Null(t *testing.T) {
	if Nil().AsInt() != 0 {
		t.Error("null.AsInt() should be 0")
	}
}

// --- AsFloat ---

func TestAsFloat_Float(t *testing.T) {
	if FloatVal(2.5).AsFloat() != 2.5 {
		t.Error("AsFloat on float")
	}
}

func TestAsFloat_Int(t *testing.T) {
	if IntVal(5).AsFloat() != 5.0 {
		t.Error("AsFloat on int")
	}
}

func TestAsFloat_Bool(t *testing.T) {
	if BoolVal(true).AsFloat() != 1.0 {
		t.Error("true.AsFloat() should be 1.0")
	}
	if BoolVal(false).AsFloat() != 0.0 {
		t.Error("false.AsFloat() should be 0.0")
	}
}

func TestAsFloat_String(t *testing.T) {
	if StringVal("3.14").AsFloat() != 3.14 {
		t.Error("'3.14'.AsFloat()")
	}
	if StringVal("bad").AsFloat() != 0 {
		t.Error("non-numeric string AsFloat should be 0")
	}
}

func TestAsFloat_Null(t *testing.T) {
	if Nil().AsFloat() != 0 {
		t.Error("null.AsFloat() should be 0")
	}
}

// --- AsString ---

func TestAsString_String(t *testing.T) {
	if StringVal("hello").AsString() != "hello" {
		t.Error("AsString on string")
	}
}

func TestAsString_Int(t *testing.T) {
	if IntVal(42).AsString() != "42" {
		t.Error("42.AsString() should be '42'")
	}
}

func TestAsString_Float(t *testing.T) {
	// whole float → x.y format
	s := FloatVal(3.0).AsString()
	if s != "3.0" {
		t.Errorf("3.0.AsString() = %q, want '3.0'", s)
	}
	s = FloatVal(3.14).AsString()
	if s != "3.14" {
		t.Errorf("3.14.AsString() = %q, want '3.14'", s)
	}
}

func TestAsString_Bool(t *testing.T) {
	if BoolVal(true).AsString() != "true" {
		t.Error("true.AsString() should be 'true'")
	}
	if BoolVal(false).AsString() != "false" {
		t.Error("false.AsString() should be 'false'")
	}
}

func TestAsString_Null(t *testing.T) {
	if Nil().AsString() != "null" {
		t.Error("null.AsString() should be 'null'")
	}
}

func TestAsString_List(t *testing.T) {
	v := ListVal(IntVal(1), IntVal(2))
	s := v.AsString()
	if s != "[1, 2]" {
		t.Errorf("list AsString = %q, want '[1, 2]'", s)
	}
}

func TestAsString_Map(t *testing.T) {
	m := MapVal()
	m.Set("a", IntVal(1))
	s := m.AsString()
	if s != `{"a": 1}` {
		t.Errorf("map AsString = %q, want {\"a\": 1}", s)
	}
}

// --- AsList ---

func TestAsList_List(t *testing.T) {
	v := ListVal(IntVal(1), IntVal(2))
	l := v.AsList()
	if len(l) != 2 {
		t.Error("AsList on list should return items")
	}
}

func TestAsList_NonList(t *testing.T) {
	v := IntVal(5)
	l := v.AsList()
	if len(l) != 1 || l[0].AsInt() != 5 {
		t.Error("AsList on non-list should wrap in slice")
	}
}

// --- AsMap ---

func TestAsMap_Map(t *testing.T) {
	m := MapVal()
	m.Set("k", IntVal(1))
	if len(m.AsMap()) != 1 {
		t.Error("AsMap on map")
	}
}

func TestAsMap_NonMap(t *testing.T) {
	v := IntVal(5)
	if len(v.AsMap()) != 0 {
		t.Error("AsMap on non-map should return empty map")
	}
}

// --- List operations ---

func TestList_Append(t *testing.T) {
	v := ListVal()
	v.Append(IntVal(1))
	v.Append(IntVal(2))
	if v.Len() != 2 {
		t.Error("Append should increase length")
	}
}

func TestList_Get(t *testing.T) {
	v := ListVal(IntVal(10), IntVal(20), IntVal(30))
	if v.Get(0).AsInt() != 10 {
		t.Error("Get(0)")
	}
	if v.Get(2).AsInt() != 30 {
		t.Error("Get(2)")
	}
	if v.Get(-1).AsInt() != 30 {
		t.Error("Get(-1) should return last element")
	}
	if v.Get(-3).AsInt() != 10 {
		t.Error("Get(-3)")
	}
}

func TestList_GetOutOfBounds(t *testing.T) {
	v := ListVal(IntVal(1))
	if !v.Get(5).IsNull() {
		t.Error("Get out of bounds should return null")
	}
	if !v.Get(-5).IsNull() {
		t.Error("Get negative out of bounds should return null")
	}
}

func TestList_Len(t *testing.T) {
	v := ListVal(IntVal(1), IntVal(2), IntVal(3))
	if v.Len() != 3 {
		t.Errorf("Len() = %d, want 3", v.Len())
	}
}

func TestString_Len(t *testing.T) {
	// UTF-8 character count
	v := StringVal("hello")
	if v.Len() != 5 {
		t.Errorf("'hello'.Len() = %d, want 5", v.Len())
	}
	// multi-byte UTF-8
	v2 := StringVal("привет")
	if v2.Len() != 6 {
		t.Errorf("'привет'.Len() = %d, want 6", v2.Len())
	}
}

func TestMap_Len(t *testing.T) {
	m := MapVal()
	m.Set("a", IntVal(1))
	m.Set("b", IntVal(2))
	if m.Len() != 2 {
		t.Errorf("map.Len() = %d, want 2", m.Len())
	}
}

// --- Map operations ---

func TestMap_SetAndGet(t *testing.T) {
	m := MapVal()
	m.Set("key", StringVal("value"))
	got := m.GetKey("key")
	if got.AsString() != "value" {
		t.Errorf("GetKey = %q, want 'value'", got.AsString())
	}
}

func TestMap_GetMissing(t *testing.T) {
	m := MapVal()
	got := m.GetKey("missing")
	if !got.IsNull() {
		t.Error("GetKey on missing key should return null")
	}
}

func TestMap_HasKey(t *testing.T) {
	m := MapVal()
	m.Set("a", IntVal(1))
	if !m.HasKey("a") {
		t.Error("HasKey existing key")
	}
	if m.HasKey("b") {
		t.Error("HasKey missing key")
	}
}

func TestMap_Keys(t *testing.T) {
	m := MapVal()
	m.Set("x", IntVal(1))
	m.Set("y", IntVal(2))
	m.Set("z", IntVal(3))
	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("Keys() len = %d, want 3", len(keys))
	}
	// insertion order preserved
	if keys[0] != "x" || keys[1] != "y" || keys[2] != "z" {
		t.Errorf("Keys() order = %v, want [x y z]", keys)
	}
}

func TestMap_SetOverwrite(t *testing.T) {
	m := MapVal()
	m.Set("k", IntVal(1))
	m.Set("k", IntVal(2))
	if m.GetKey("k").AsInt() != 2 {
		t.Error("Set should overwrite existing key")
	}
	// key should not be duplicated
	if len(m.Keys()) != 1 {
		t.Errorf("overwrite should not duplicate key, got %d keys", len(m.Keys()))
	}
}

// --- Repr ---

func TestRepr_String(t *testing.T) {
	v := StringVal("hello")
	if v.Repr() != `"hello"` {
		t.Errorf("Repr of string = %q, want %q", v.Repr(), `"hello"`)
	}
}

func TestRepr_Int(t *testing.T) {
	v := IntVal(42)
	if v.Repr() != "42" {
		t.Errorf("Repr of int = %q, want '42'", v.Repr())
	}
}

// --- Equal ---

func TestEqual_SameType(t *testing.T) {
	cases := []struct {
		a, b  *Value
		equal bool
	}{
		{Nil(), Nil(), true},
		{BoolVal(true), BoolVal(true), true},
		{BoolVal(true), BoolVal(false), false},
		{IntVal(5), IntVal(5), true},
		{IntVal(5), IntVal(6), false},
		{FloatVal(1.5), FloatVal(1.5), true},
		{FloatVal(1.5), FloatVal(2.5), false},
		{StringVal("a"), StringVal("a"), true},
		{StringVal("a"), StringVal("b"), false},
	}
	for _, c := range cases {
		got := c.a.Equal(c.b)
		if got != c.equal {
			t.Errorf("Equal(%v, %v) = %v, want %v", c.a, c.b, got, c.equal)
		}
	}
}

func TestEqual_CrossTypeNumeric(t *testing.T) {
	// int == float cross-type
	if !IntVal(5).Equal(FloatVal(5.0)) {
		t.Error("5 (int) should equal 5.0 (float)")
	}
	if IntVal(5).Equal(FloatVal(5.1)) {
		t.Error("5 should not equal 5.1")
	}
}

func TestEqual_DifferentTypes(t *testing.T) {
	if IntVal(1).Equal(BoolVal(true)) {
		t.Error("int and bool should not be equal")
	}
	if StringVal("1").Equal(IntVal(1)) {
		t.Error("string and int should not be equal")
	}
}

func TestEqual_List(t *testing.T) {
	a := ListVal(IntVal(1), IntVal(2))
	b := ListVal(IntVal(1), IntVal(2))
	c := ListVal(IntVal(1))
	if !a.Equal(b) {
		t.Error("equal lists should be equal")
	}
	if a.Equal(c) {
		t.Error("lists of different length should not be equal")
	}
}

func TestEqual_Map(t *testing.T) {
	a := MapVal()
	a.Set("k", IntVal(1))
	b := MapVal()
	b.Set("k", IntVal(1))
	c := MapVal()
	c.Set("k", IntVal(2))
	if !a.Equal(b) {
		t.Error("equal maps should be equal")
	}
	if a.Equal(c) {
		t.Error("maps with different values should not be equal")
	}
}

// --- Compare ---

func TestCompare_Numeric(t *testing.T) {
	cases := []struct {
		a, b *Value
		want int
	}{
		{IntVal(1), IntVal(2), -1},
		{IntVal(2), IntVal(1), 1},
		{IntVal(1), IntVal(1), 0},
		{FloatVal(1.5), FloatVal(2.5), -1},
		{FloatVal(2.5), FloatVal(1.5), 1},
		{IntVal(5), FloatVal(5.0), 0},
		{IntVal(5), FloatVal(5.1), -1},
	}
	for _, c := range cases {
		got := c.a.Compare(c.b)
		if got != c.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestCompare_String(t *testing.T) {
	if StringVal("a").Compare(StringVal("b")) != -1 {
		t.Error("'a' < 'b'")
	}
	if StringVal("b").Compare(StringVal("a")) != 1 {
		t.Error("'b' > 'a'")
	}
	if StringVal("a").Compare(StringVal("a")) != 0 {
		t.Error("'a' == 'a'")
	}
}

func TestCompare_MixedTypes(t *testing.T) {
	// non-numeric non-string compare returns 0
	if Nil().Compare(IntVal(1)) != 0 {
		t.Error("Compare of incompatible types should return 0")
	}
}

// --- String (Stringer) ---

func TestString_Method(t *testing.T) {
	v := IntVal(42)
	if v.String() != "42" {
		t.Error("String() should call AsString()")
	}
}

// --- Type() accessor ---

func TestType_Accessor(t *testing.T) {
	if Nil().Type() != Null {
		t.Error("Nil().Type() should be Null")
	}
	if BoolVal(true).Type() != Bool {
		t.Error("BoolVal.Type() should be Bool")
	}
	if IntVal(1).Type() != Int {
		t.Error("IntVal.Type() should be Int")
	}
	if FloatVal(1.0).Type() != Float {
		t.Error("FloatVal.Type() should be Float")
	}
	if StringVal("").Type() != String {
		t.Error("StringVal.Type() should be String")
	}
	if ListVal().Type() != List {
		t.Error("ListVal.Type() should be List")
	}
	if MapVal().Type() != Map {
		t.Error("MapVal.Type() should be Map")
	}
}

// --- Special float values ---

func TestAsString_NaN(t *testing.T) {
	v := FloatVal(math.NaN())
	_ = v.AsString() // should not panic
}

func TestAsString_Inf(t *testing.T) {
	v := FloatVal(math.Inf(1))
	_ = v.AsString() // should not panic
}
