package eval

import "testing"

func TestList_Append(t *testing.T) {
	mustStr(t, "append([1, 2], 3)[2]", "3")
	mustInt(t, "len(append([1, 2], 3))", 3)
	mustEvalErr(t, "append([1])")
}

func TestList_Prepend(t *testing.T) {
	mustStr(t, "prepend([2, 3], 1)[0]", "1")
	mustInt(t, "len(prepend([2, 3], 1))", 3)
	mustEvalErr(t, "prepend([1])")
}

func TestList_Concat(t *testing.T) {
	mustInt(t, "len(concat([1, 2], [3, 4]))", 4)
	mustStr(t, "concat([1, 2], [3, 4])[3]", "4")
	mustInt(t, "len(concat([1], [2], [3]))", 3)
	mustEvalErr(t, "concat([1])")
}

func TestList_First(t *testing.T) {
	mustStr(t, "first([10, 20, 30])", "10")
	v := mustEval(t, "first([])")
	if !v.IsNull() {
		t.Errorf("first([]) should be null")
	}
	mustInt(t, "len(first([1, 2, 3], 2))", 2)
	mustStr(t, "first([1, 2, 3], 2)[0]", "1")
	mustInt(t, "len(first([1, 2], 5))", 2)
	mustEvalErr(t, "first()")
}

func TestList_Last(t *testing.T) {
	mustStr(t, "last([10, 20, 30])", "30")
	v := mustEval(t, "last([])")
	if !v.IsNull() {
		t.Errorf("last([]) should be null")
	}
	mustInt(t, "len(last([1, 2, 3], 2))", 2)
	mustStr(t, "last([1, 2, 3], 2)[1]", "3")
	mustInt(t, "len(last([1, 2], 5))", 2)
	mustEvalErr(t, "last()")
}

func TestList_Nth(t *testing.T) {
	mustStr(t, "nth([10, 20, 30], 1)", "20")
	mustStr(t, "nth([10, 20, 30], -1)", "30")
	v := mustEval(t, "nth([1, 2], 5)")
	if !v.IsNull() {
		t.Errorf("nth out-of-bounds should be null")
	}
	v = mustEval(t, "nth([1, 2], -5)")
	if !v.IsNull() {
		t.Errorf("nth negative out-of-bounds should be null")
	}
	mustEvalErr(t, "nth([1])")
}

func TestList_Slice(t *testing.T) {
	mustStr(t, "slice([1, 2, 3, 4], 1)[0]", "2")
	mustInt(t, "len(slice([1, 2, 3, 4], 1))", 3)
	mustInt(t, "len(slice([1, 2, 3, 4], 1, 3))", 2)
	mustInt(t, "len(slice([1, 2, 3], 5))", 0)
	mustStr(t, "slice([1, 2, 3], -2)[0]", "2")
	mustInt(t, "len(slice([1, 2, 3], -100))", 3)
	mustStr(t, "slice([1, 2, 3], 0, -1)[0]", "1")
	mustInt(t, "len(slice([1, 2, 3], 0, -1))", 2)
	mustInt(t, "len(slice([1, 2, 3], 2, 1))", 0)
	mustInt(t, "len(slice([1, 2, 3], 1, 100))", 2)
	mustEvalErr(t, "slice([1])")
}

func TestList_Take(t *testing.T) {
	mustInt(t, "len(take([1, 2, 3], 2))", 2)
	mustStr(t, "take([1, 2, 3], 2)[0]", "1")
	mustInt(t, "len(take([1, 2, 3], 5))", 3)
	mustInt(t, "len(take([1, 2, 3], -1))", 0)
	mustEvalErr(t, "take([1])")
}

func TestList_Drop(t *testing.T) {
	mustInt(t, "len(drop([1, 2, 3], 1))", 2)
	mustStr(t, "drop([1, 2, 3], 1)[0]", "2")
	mustInt(t, "len(drop([1, 2, 3], 5))", 0)
	mustInt(t, "len(drop([1, 2, 3], -1))", 3)
	mustEvalErr(t, "drop([1])")
}

func TestList_Contains(t *testing.T) {
	mustBool(t, "contains([1, 2, 3], 2)", true)
	mustBool(t, "contains([1, 2, 3], 5)", false)
	// string path (overridden to handle both)
	mustBool(t, "contains('hello', 'ell')", true)
	mustBool(t, "contains('hello', 'xyz')", false)
	mustEvalErr(t, "contains([1])")
}

func TestList_IndexOf(t *testing.T) {
	mustInt(t, "indexOf([10, 20, 30], 20)", 1)
	mustInt(t, "indexOf([10, 20, 30], 99)", -1)
	mustInt(t, "indexOf('hello', 'l')", 2)
	mustInt(t, "indexOf('hello', 'z')", -1)
	mustEvalErr(t, "indexOf([1])")
}

func TestList_LastIndexOf(t *testing.T) {
	mustInt(t, "lastIndexOf([1, 2, 1], 1)", 2)
	mustInt(t, "lastIndexOf([1, 2, 3], 99)", -1)
	mustInt(t, "lastIndexOf('hello', 'l')", 3)
	mustInt(t, "lastIndexOf('hello', 'z')", -1)
	mustEvalErr(t, "lastIndexOf([1])")
}

func TestList_Flatten(t *testing.T) {
	mustInt(t, "len(flatten([[1, 2], [3, 4]]))", 4)
	mustStr(t, "flatten([[1, 2], [3, 4]])[0]", "1")
	// non-list items pass through
	mustInt(t, "len(flatten([1, [2, 3], 4]))", 4)
	mustStr(t, "flatten([1, [2, 3], 4])[0]", "1")
	mustEvalErr(t, "flatten()")
}

func TestList_FlattenAll(t *testing.T) {
	mustInt(t, "len(flattenAll([1, [2, [3, 4]]]))", 4)
	mustStr(t, "flattenAll([1, [2, [3, 4]]])[3]", "4")
	// non-list leaf items
	mustStr(t, "flattenAll([1, 2])[0]", "1")
	mustEvalErr(t, "flattenAll()")
}

func TestList_Unique(t *testing.T) {
	mustInt(t, "len(unique([1, 2, 1, 3, 2]))", 3)
	mustStr(t, "unique([1, 2, 1, 3, 2])[0]", "1")
	mustInt(t, "len(unique([]))", 0)
	mustEvalErr(t, "unique()")
}

func TestList_Sort(t *testing.T) {
	mustStr(t, "sort([3, 1, 2])[0]", "1")
	mustStr(t, "sort([3, 1, 2])[2]", "3")
	mustStr(t, "sort(['c', 'a', 'b'])[0]", "a")
	mustEvalErr(t, "sort()")
}

func TestList_SortDesc(t *testing.T) {
	mustStr(t, "sortDesc([1, 3, 2])[0]", "3")
	mustStr(t, "sortDesc([1, 3, 2])[2]", "1")
	mustStr(t, "sortDesc(['c', 'a', 'b'])[0]", "c")
	mustEvalErr(t, "sortDesc()")
}

func TestList_Range(t *testing.T) {
	mustInt(t, "len(range(5))", 5)
	mustStr(t, "range(5)[0]", "0")
	mustStr(t, "range(5)[4]", "4")

	mustInt(t, "len(range(2, 5))", 3)
	mustStr(t, "range(2, 5)[0]", "2")

	mustInt(t, "len(range(0, 6, 2))", 3)
	mustStr(t, "range(0, 6, 2)[1]", "2")

	mustInt(t, "len(range(5, 2, -1))", 3)
	mustStr(t, "range(5, 2, -1)[0]", "5")
	mustStr(t, "range(5, 2, -1)[2]", "3")

	mustEvalErr(t, "range(1, 5, 0)")
	mustEvalErr(t, "range()")
}

func TestList_Chunk(t *testing.T) {
	v := mustEval(t, "chunk([1, 2, 3, 4, 5], 2)")
	if v.Len() != 3 {
		t.Errorf("chunk: expected 3 chunks, got %d", v.Len())
	}
	mustInt(t, "len(chunk([1, 2, 3, 4, 5], 2)[0])", 2)
	mustInt(t, "len(chunk([1, 2, 3, 4, 5], 2)[2])", 1)
	mustEvalErr(t, "chunk([1, 2], 0)")
	mustEvalErr(t, "chunk([1])")
}

func TestList_Zip(t *testing.T) {
	v := mustEval(t, "zip([1, 2, 3], ['a', 'b'])")
	if v.Len() != 2 {
		t.Errorf("zip: expected 2 pairs, got %d", v.Len())
	}
	mustStr(t, "zip([1, 2], ['a', 'b', 'c'])[0][0]", "1")
	mustStr(t, "zip([1, 2], ['a', 'b', 'c'])[0][1]", "a")
	mustEvalErr(t, "zip([1])")
}

func TestList_Without(t *testing.T) {
	mustInt(t, "len(without([1, 2, 3, 2], 2))", 2)
	mustStr(t, "without([1, 2, 3, 2], 2)[0]", "1")
	mustInt(t, "len(without([1, 2, 3], 2, 3))", 1)
	mustStr(t, "without([1, 2, 3], 2, 3)[0]", "1")
	mustEvalErr(t, "without([1])")
}

func TestList_Count(t *testing.T) {
	mustInt(t, "count([1, 2, 3])", 3)
	mustInt(t, "count([])", 0)
	mustEvalErr(t, "count()")
}
