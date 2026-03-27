package eval

import "testing"

func TestMap_Keys(t *testing.T) {
	mustInt(t, "len(keys({'a': 1, 'b': 2}))", 2)
	mustStr(t, "keys({'a': 1, 'b': 2})[0]", "a")
	mustStr(t, "keys({'a': 1, 'b': 2})[1]", "b")
	mustEvalErr(t, "keys()")
}

func TestMap_Values(t *testing.T) {
	mustInt(t, "len(values({'a': 1, 'b': 2}))", 2)
	mustStr(t, "values({'a': 1, 'b': 2})[0]", "1")
	mustStr(t, "values({'a': 1, 'b': 2})[1]", "2")
	mustEvalErr(t, "values()")
}

func TestMap_HasKey(t *testing.T) {
	mustBool(t, "hasKey({'a': 1}, 'a')", true)
	mustBool(t, "hasKey({'a': 1}, 'b')", false)
	mustEvalErr(t, "hasKey({'a': 1})")
}

func TestMap_Get(t *testing.T) {
	mustStr(t, "get({'a': 1}, 'a')", "1")
	v := mustEval(t, "get({'a': 1}, 'b')")
	if !v.IsNull() {
		t.Errorf("get missing key should be null, got %v", v.AsString())
	}
	mustStr(t, "get({'a': 1}, 'b', 99)", "99")
	mustEvalErr(t, "get({'a': 1})")
}

func TestMap_Set(t *testing.T) {
	mustBool(t, "hasKey(set({'a': 1}, 'b', 2), 'b')", true)
	mustStr(t, "set({'a': 1}, 'b', 2)['b']", "2")
	// existing keys are preserved
	mustBool(t, "hasKey(set({'a': 1}, 'b', 2), 'a')", true)
	mustEvalErr(t, "set({'a': 1}, 'b')")
}

func TestMap_Delete(t *testing.T) {
	mustBool(t, "hasKey(delete({'a': 1, 'b': 2}, 'a'), 'a')", false)
	mustBool(t, "hasKey(delete({'a': 1, 'b': 2}, 'a'), 'b')", true)
	// delete multiple keys
	mustInt(t, "len(keys(delete({'a': 1, 'b': 2, 'c': 3}, 'a', 'b')))", 1)
	mustEvalErr(t, "delete({'a': 1})")
}

func TestMap_Merge(t *testing.T) {
	mustBool(t, "hasKey(merge({'a': 1}, {'b': 2}), 'a')", true)
	mustBool(t, "hasKey(merge({'a': 1}, {'b': 2}), 'b')", true)
	// later map wins on conflict
	mustStr(t, "merge({'a': 1}, {'a': 2})['a']", "2")
	// three maps
	mustInt(t, "len(keys(merge({'a': 1}, {'b': 2}, {'c': 3})))", 3)
	mustEvalErr(t, "merge({'a': 1})")
}

func TestMap_Entries(t *testing.T) {
	v := mustEval(t, "entries({'a': 1})")
	if v.Len() != 1 {
		t.Errorf("entries: expected 1 entry, got %d", v.Len())
	}
	mustStr(t, "entries({'a': 1})[0][0]", "a")
	mustStr(t, "entries({'a': 1})[0][1]", "1")
	mustEvalErr(t, "entries()")
}

func TestMap_FromEntries(t *testing.T) {
	mustBool(t, "hasKey(fromEntries([['a', 1], ['b', 2]]), 'a')", true)
	mustStr(t, "fromEntries([['a', 1]])['a']", "1")
	// entries with fewer than 2 items are silently skipped
	mustInt(t, "len(keys(fromEntries([[1], ['a', 2]])))", 1)
	mustEvalErr(t, "fromEntries()")
}

func TestMap_Pick(t *testing.T) {
	mustBool(t, "hasKey(pick({'a': 1, 'b': 2, 'c': 3}, 'a', 'c'), 'a')", true)
	mustBool(t, "hasKey(pick({'a': 1, 'b': 2, 'c': 3}, 'a', 'c'), 'b')", false)
	mustBool(t, "hasKey(pick({'a': 1, 'b': 2, 'c': 3}, 'a', 'c'), 'c')", true)
	// key not in map is simply absent
	mustInt(t, "len(keys(pick({'a': 1}, 'x')))", 0)
	mustEvalErr(t, "pick({'a': 1})")
}

func TestMap_Omit(t *testing.T) {
	mustBool(t, "hasKey(omit({'a': 1, 'b': 2}, 'a'), 'a')", false)
	mustBool(t, "hasKey(omit({'a': 1, 'b': 2}, 'a'), 'b')", true)
	mustEvalErr(t, "omit({'a': 1})")
}
