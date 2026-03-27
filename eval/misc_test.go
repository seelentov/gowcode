package eval

import (
	"strings"
	"testing"
	"time"
)

func TestMisc_Print(t *testing.T) {
	// print always returns null regardless of args
	v := mustEval(t, "print('hello')")
	if !v.IsNull() {
		t.Errorf("print should return null, got %v", v.AsString())
	}
	v = mustEval(t, "print('a', 'b', 'c')")
	if !v.IsNull() {
		t.Errorf("print(multi) should return null")
	}
	v = mustEval(t, "print()")
	if !v.IsNull() {
		t.Errorf("print() should return null")
	}
}

func TestMisc_Error(t *testing.T) {
	_, err := Eval("error('something went wrong')", nil)
	if err == nil {
		t.Fatal("error() should return an error")
	}
	if !strings.Contains(err.Error(), "something went wrong") {
		t.Errorf("error message mismatch: %v", err)
	}
	mustEvalErr(t, "error()")
}

func TestMisc_UUID(t *testing.T) {
	v := mustEval(t, "uuid()")
	uid := v.AsString()
	if len(uid) != 36 {
		t.Errorf("uuid() length = %d, want 36; got %q", len(uid), uid)
	}
	// format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	parts := strings.Split(uid, "-")
	if len(parts) != 5 {
		t.Errorf("uuid() should have 5 parts separated by '-', got %q", uid)
	}
}

func TestMisc_Now(t *testing.T) {
	before := time.Now().Unix()
	v := mustEval(t, "now()")
	after := time.Now().Unix()
	ts := v.AsInt()
	if ts < before || ts > after {
		t.Errorf("now() = %d, want [%d, %d]", ts, before, after)
	}
}

func TestMisc_NowMs(t *testing.T) {
	before := time.Now().UnixMilli()
	v := mustEval(t, "nowMs()")
	after := time.Now().UnixMilli()
	ts := v.AsInt()
	if ts < before || ts > after {
		t.Errorf("nowMs() = %d, want [%d, %d]", ts, before, after)
	}
}

func TestMisc_DateFormat(t *testing.T) {
	// Unix epoch in UTC: 1970-01-01
	mustStr(t, "dateFormat(0, '2006-01-02')", "1970-01-01")
	mustEvalErr(t, "dateFormat(0)")
}

func TestMisc_DateParse(t *testing.T) {
	// Parsing the epoch date should yield 0 (UTC)
	v := mustEval(t, "dateParse('2006-01-02', '1970-01-01')")
	// time.Parse uses UTC, so this should be exactly 0
	if v.AsInt() != 0 {
		t.Errorf("dateParse epoch = %d, want 0", v.AsInt())
	}
	// parse error
	mustEvalErr(t, "dateParse('2006-01-02', 'not-a-date')")
	mustEvalErr(t, "dateParse('2006-01-02')")
}

func TestMisc_Sleep(t *testing.T) {
	v := mustEval(t, "sleep(0)")
	if !v.IsNull() {
		t.Errorf("sleep(0) should return null, got %v", v.AsString())
	}
	mustEvalErr(t, "sleep()")
}
