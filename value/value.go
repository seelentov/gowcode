package value

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Type int

const (
	Null Type = iota
	Bool
	Int
	Float
	String
	List
	Map
)

type Value struct {
	typ   Type
	bval  bool
	ival  int64
	fval  float64
	sval  string
	lval  []*Value
	mval  map[string]*Value
	mkeys []string // preserve insertion order
}

// Constructors

func Nil() *Value                  { return &Value{typ: Null} }
func BoolVal(b bool) *Value        { return &Value{typ: Bool, bval: b} }
func IntVal(i int64) *Value        { return &Value{typ: Int, ival: i} }
func FloatVal(f float64) *Value    { return &Value{typ: Float, fval: f} }
func StringVal(s string) *Value    { return &Value{typ: String, sval: s} }
func ListVal(items ...*Value) *Value {
	if items == nil {
		items = []*Value{}
	}
	return &Value{typ: List, lval: items}
}
func MapVal() *Value {
	return &Value{typ: Map, mval: map[string]*Value{}, mkeys: []string{}}
}

// Type checks

func (v *Value) Type() Type    { return v.typ }
func (v *Value) IsNull() bool  { return v.typ == Null }
func (v *Value) IsBool() bool  { return v.typ == Bool }
func (v *Value) IsInt() bool   { return v.typ == Int }
func (v *Value) IsFloat() bool { return v.typ == Float }
func (v *Value) IsNumber() bool { return v.typ == Int || v.typ == Float }
func (v *Value) IsString() bool { return v.typ == String }
func (v *Value) IsList() bool  { return v.typ == List }
func (v *Value) IsMap() bool   { return v.typ == Map }

// Raw accessors

func (v *Value) AsBool() bool {
	switch v.typ {
	case Bool:
		return v.bval
	case Int:
		return v.ival != 0
	case Float:
		return v.fval != 0
	case String:
		return v.sval != "" && v.sval != "false" && v.sval != "0"
	case Null:
		return false
	case List:
		return len(v.lval) > 0
	case Map:
		return len(v.mval) > 0
	}
	return false
}

func (v *Value) AsInt() int64 {
	switch v.typ {
	case Int:
		return v.ival
	case Float:
		return int64(v.fval)
	case Bool:
		if v.bval {
			return 1
		}
		return 0
	case String:
		i, _ := strconv.ParseInt(v.sval, 10, 64)
		return i
	}
	return 0
}

func (v *Value) AsFloat() float64 {
	switch v.typ {
	case Float:
		return v.fval
	case Int:
		return float64(v.ival)
	case Bool:
		if v.bval {
			return 1
		}
		return 0
	case String:
		f, _ := strconv.ParseFloat(v.sval, 64)
		return f
	}
	return 0
}

func (v *Value) AsString() string {
	switch v.typ {
	case String:
		return v.sval
	case Int:
		return strconv.FormatInt(v.ival, 10)
	case Float:
		f := v.fval
		if f == math.Trunc(f) {
			return strconv.FormatFloat(f, 'f', 1, 64)
		}
		return strconv.FormatFloat(f, 'f', -1, 64)
	case Bool:
		if v.bval {
			return "true"
		}
		return "false"
	case Null:
		return "null"
	case List:
		parts := make([]string, len(v.lval))
		for i, item := range v.lval {
			parts[i] = item.Repr()
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case Map:
		parts := make([]string, 0, len(v.mkeys))
		for _, k := range v.mkeys {
			parts = append(parts, fmt.Sprintf("%q: %s", k, v.mval[k].Repr()))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	}
	return ""
}

func (v *Value) AsList() []*Value {
	if v.typ == List {
		return v.lval
	}
	return []*Value{v}
}

func (v *Value) AsMap() map[string]*Value {
	if v.typ == Map {
		return v.mval
	}
	return map[string]*Value{}
}

// List operations

func (v *Value) Append(item *Value) {
	v.lval = append(v.lval, item)
}

func (v *Value) Get(index int) *Value {
	if index < 0 {
		index = len(v.lval) + index
	}
	if index < 0 || index >= len(v.lval) {
		return Nil()
	}
	return v.lval[index]
}

func (v *Value) Len() int {
	switch v.typ {
	case List:
		return len(v.lval)
	case Map:
		return len(v.mval)
	case String:
		return len([]rune(v.sval))
	}
	return 0
}

// Map operations

func (v *Value) Set(key string, val *Value) {
	if _, exists := v.mval[key]; !exists {
		v.mkeys = append(v.mkeys, key)
	}
	v.mval[key] = val
}

func (v *Value) GetKey(key string) *Value {
	if val, ok := v.mval[key]; ok {
		return val
	}
	return Nil()
}

func (v *Value) Keys() []string {
	return v.mkeys
}

func (v *Value) HasKey(key string) bool {
	_, ok := v.mval[key]
	return ok
}

// Repr returns a string representation suitable for display
func (v *Value) Repr() string {
	if v.typ == String {
		return fmt.Sprintf("%q", v.sval)
	}
	return v.AsString()
}

// Equal compares two values
func (v *Value) Equal(other *Value) bool {
	if v.typ != other.typ {
		// cross-type numeric comparison
		if v.IsNumber() && other.IsNumber() {
			return v.AsFloat() == other.AsFloat()
		}
		return false
	}
	switch v.typ {
	case Null:
		return true
	case Bool:
		return v.bval == other.bval
	case Int:
		return v.ival == other.ival
	case Float:
		return v.fval == other.fval
	case String:
		return v.sval == other.sval
	case List:
		if len(v.lval) != len(other.lval) {
			return false
		}
		for i := range v.lval {
			if !v.lval[i].Equal(other.lval[i]) {
				return false
			}
		}
		return true
	case Map:
		if len(v.mval) != len(other.mval) {
			return false
		}
		for k, val := range v.mval {
			oval, ok := other.mval[k]
			if !ok || !val.Equal(oval) {
				return false
			}
		}
		return true
	}
	return false
}

func (v *Value) Compare(other *Value) int {
	if v.IsNumber() && other.IsNumber() {
		a, b := v.AsFloat(), other.AsFloat()
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	if v.IsString() && other.IsString() {
		a, b := v.sval, other.sval
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	return 0
}

func (v *Value) String() string {
	return v.AsString()
}
