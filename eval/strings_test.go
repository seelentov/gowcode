package eval

import "testing"

func TestStrings_Upper(t *testing.T) {
	mustStr(t, "upper('hello')", "HELLO")
	mustStr(t, "upper('')", "")
	mustEvalErr(t, "upper()")
	mustEvalErr(t, "upper('a', 'b')")
}

func TestStrings_Lower(t *testing.T) {
	mustStr(t, "lower('HELLO')", "hello")
	mustEvalErr(t, "lower()")
}

func TestStrings_Trim(t *testing.T) {
	mustStr(t, "trim('  hi  ')", "hi")
	mustStr(t, "trim('xxhixx', 'x')", "hi")
	mustEvalErr(t, "trim()")
}

func TestStrings_TrimLeft(t *testing.T) {
	mustStr(t, "trimLeft('  hi  ')", "hi  ")
	mustStr(t, "trimLeft('xxhi', 'x')", "hi")
	mustEvalErr(t, "trimLeft()")
}

func TestStrings_TrimRight(t *testing.T) {
	mustStr(t, "trimRight('  hi  ')", "  hi")
	mustStr(t, "trimRight('hixx', 'x')", "hi")
	mustEvalErr(t, "trimRight()")
}

func TestStrings_TrimPrefix(t *testing.T) {
	mustStr(t, "trimPrefix('foobar', 'foo')", "bar")
	mustStr(t, "trimPrefix('foobar', 'baz')", "foobar")
	mustEvalErr(t, "trimPrefix('a')")
}

func TestStrings_TrimSuffix(t *testing.T) {
	mustStr(t, "trimSuffix('foobar', 'bar')", "foo")
	mustStr(t, "trimSuffix('foobar', 'baz')", "foobar")
	mustEvalErr(t, "trimSuffix('a')")
}

func TestStrings_Replace(t *testing.T) {
	mustStr(t, "replace('aabbcc', 'b', 'x')", "aaxxcc")
	mustStr(t, "replace('aabbcc', 'b', 'x', 1)", "aaxbcc")
	mustEvalErr(t, "replace('a', 'b')")
}

func TestStrings_ReplaceAll(t *testing.T) {
	mustStr(t, "replaceAll('aabbcc', 'b', 'x')", "aaxxcc")
	mustEvalErr(t, "replaceAll('a', 'b')")
}

func TestStrings_Contains(t *testing.T) {
	mustBool(t, "contains('hello', 'ell')", true)
	mustBool(t, "contains('hello', 'xyz')", false)
	mustEvalErr(t, "contains('a')")
}

func TestStrings_StartsWith(t *testing.T) {
	mustBool(t, "startsWith('hello', 'hel')", true)
	mustBool(t, "startsWith('hello', 'ell')", false)
	mustEvalErr(t, "startsWith('a')")
}

func TestStrings_EndsWith(t *testing.T) {
	mustBool(t, "endsWith('hello', 'llo')", true)
	mustBool(t, "endsWith('hello', 'hel')", false)
	mustEvalErr(t, "endsWith('a')")
}

func TestStrings_Split(t *testing.T) {
	v := mustEval(t, "split('a,b,c', ',')")
	if v.Len() != 3 {
		t.Errorf("expected 3 parts, got %d", v.Len())
	}
	mustStr(t, "split('a,b,c', ',')[0]", "a")
	mustStr(t, "split('a,b,c', ',')[2]", "c")
	mustEvalErr(t, "split('a')")
}

func TestStrings_Join(t *testing.T) {
	mustStr(t, "join(['a', 'b', 'c'], ',')", "a,b,c")
	mustStr(t, "join([], ',')", "")
	mustEvalErr(t, "join(['a'])")
}

func TestStrings_Len(t *testing.T) {
	mustInt(t, "len('hello')", 5)
	mustInt(t, "len('')", 0)
	mustInt(t, "len([1, 2, 3])", 3)
	mustInt(t, "len({'a': 1, 'b': 2})", 2)
	mustEvalErr(t, "len()")
}

func TestStrings_Substr(t *testing.T) {
	mustStr(t, "substr('hello', 1)", "ello")
	mustStr(t, "substr('hello', 1, 3)", "el")
	mustStr(t, "substr('hello', -2)", "lo")
	mustStr(t, "substr('hello', -100)", "hello")
	mustStr(t, "substr('hello', 100)", "")
	mustStr(t, "substr('hello', 3, 1)", "")
	mustStr(t, "substr('hello', 1, 100)", "ello")
	mustStr(t, "substr('hello', 1, -1)", "ell")
	mustEvalErr(t, "substr('a')")
}

func TestStrings_IndexOf(t *testing.T) {
	mustInt(t, "indexOf('hello', 'l')", 2)
	mustInt(t, "indexOf('hello', 'z')", -1)
	mustEvalErr(t, "indexOf('a')")
}

func TestStrings_LastIndexOf(t *testing.T) {
	mustInt(t, "lastIndexOf('hello', 'l')", 3)
	mustInt(t, "lastIndexOf('hello', 'z')", -1)
	mustEvalErr(t, "lastIndexOf('a')")
}

func TestStrings_Repeat(t *testing.T) {
	mustStr(t, "repeat('ab', 3)", "ababab")
	mustStr(t, "repeat('x', 0)", "")
	mustEvalErr(t, "repeat('x')")
}

func TestStrings_PadLeft(t *testing.T) {
	mustStr(t, "padLeft('hi', 5)", "   hi")
	mustStr(t, "padLeft('hello', 3)", "hello")
	mustStr(t, "padLeft('hi', 5, '0')", "000hi")
	mustEvalErr(t, "padLeft('x')")
}

func TestStrings_PadRight(t *testing.T) {
	mustStr(t, "padRight('hi', 5)", "hi   ")
	mustStr(t, "padRight('hello', 3)", "hello")
	mustStr(t, "padRight('hi', 5, '0')", "hi000")
	mustEvalErr(t, "padRight('x')")
}

func TestStrings_Format(t *testing.T) {
	mustStr(t, "format('%s %d', 'hello', 42)", "hello 42")
	mustStr(t, "format('%.2f', 3.14159)", "3.14")
	mustStr(t, "format('%v', true)", "true")
	mustStr(t, "format('hello')", "hello")
	mustEvalErr(t, "format()")
}

func TestStrings_CharAt(t *testing.T) {
	mustStr(t, "charAt('hello', 0)", "h")
	mustStr(t, "charAt('hello', -1)", "o")
	mustStr(t, "charAt('hello', 100)", "")
	mustStr(t, "charAt('hello', -100)", "")
	mustEvalErr(t, "charAt('a')")
}

func TestStrings_Reverse(t *testing.T) {
	mustStr(t, "reverse('hello')", "olleh")
	mustStr(t, "reverse('')", "")
	v := mustEval(t, "reverse([1, 2, 3])")
	if v.Len() != 3 || v.Get(0).AsInt() != 3 {
		t.Errorf("reverse([1,2,3]) wrong: %v", v.AsString())
	}
	// non-string, non-list → returns value as-is
	mustInt(t, "reverse(42)", 42)
	mustEvalErr(t, "reverse()")
}
