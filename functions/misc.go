package functions

import (
	"crypto/rand"
	"fmt"
	"time"

	"gowcode/value"
)

func registerMiscFuncs(r *Registry) {
	// print(v, ...) — prints arguments to stdout, returns null
	r.Register("print", func(args []*value.Value) (*value.Value, error) {
		for i, a := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(a.AsString())
		}
		fmt.Println()
		return value.Nil(), nil
	})

	// error(msg) — returns an error with the given message
	r.Register("error", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "error"); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", args[0].AsString())
	})

	// uuid() — generates a random UUID v4 string
	r.Register("uuid", func(args []*value.Value) (*value.Value, error) {
		var b [16]byte
		_, _ = rand.Read(b[:]) // crypto/rand.Read never fails
		b[6] = (b[6] & 0x0f) | 0x40 // version 4
		b[8] = (b[8] & 0x3f) | 0x80 // variant bits
		u := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		return value.StringVal(u), nil
	})

	// now() — current Unix timestamp in seconds
	r.Register("now", func(args []*value.Value) (*value.Value, error) {
		return value.IntVal(time.Now().Unix()), nil
	})

	// nowMs() — current Unix timestamp in milliseconds
	r.Register("nowMs", func(args []*value.Value) (*value.Value, error) {
		return value.IntVal(time.Now().UnixMilli()), nil
	})

	// dateFormat(unixSec, layout) — formats a Unix timestamp using Go time layout
	r.Register("dateFormat", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "dateFormat"); err != nil {
			return nil, err
		}
		t := time.Unix(args[0].AsInt(), 0)
		return value.StringVal(t.Format(args[1].AsString())), nil
	})

	// dateParse(layout, str) — parses a date string, returns Unix timestamp
	r.Register("dateParse", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "dateParse"); err != nil {
			return nil, err
		}
		t, err := time.Parse(args[0].AsString(), args[1].AsString())
		if err != nil {
			return nil, fmt.Errorf("dateParse: %w", err)
		}
		return value.IntVal(t.Unix()), nil
	})

	// sleep(ms) — pauses execution for the given number of milliseconds
	r.Register("sleep", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "sleep"); err != nil {
			return nil, err
		}
		time.Sleep(time.Duration(args[0].AsInt()) * time.Millisecond)
		return value.Nil(), nil
	})
}
