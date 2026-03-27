package functions

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"gowcode/value"
)

func registerStringFuncs(r *Registry) {
	r.Register("upper", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "upper"); err != nil {
			return nil, err
		}
		return value.StringVal(strings.ToUpper(args[0].AsString())), nil
	})

	r.Register("lower", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "lower"); err != nil {
			return nil, err
		}
		return value.StringVal(strings.ToLower(args[0].AsString())), nil
	})

	r.Register("trim", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "trim"); err != nil {
			return nil, err
		}
		s := args[0].AsString()
		if len(args) == 2 {
			return value.StringVal(strings.Trim(s, args[1].AsString())), nil
		}
		return value.StringVal(strings.TrimSpace(s)), nil
	})

	r.Register("trimLeft", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "trimLeft"); err != nil {
			return nil, err
		}
		s := args[0].AsString()
		if len(args) == 2 {
			return value.StringVal(strings.TrimLeft(s, args[1].AsString())), nil
		}
		return value.StringVal(strings.TrimLeft(s, " \t\n\r")), nil
	})

	r.Register("trimRight", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "trimRight"); err != nil {
			return nil, err
		}
		s := args[0].AsString()
		if len(args) == 2 {
			return value.StringVal(strings.TrimRight(s, args[1].AsString())), nil
		}
		return value.StringVal(strings.TrimRight(s, " \t\n\r")), nil
	})

	r.Register("trimPrefix", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "trimPrefix"); err != nil {
			return nil, err
		}
		return value.StringVal(strings.TrimPrefix(args[0].AsString(), args[1].AsString())), nil
	})

	r.Register("trimSuffix", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "trimSuffix"); err != nil {
			return nil, err
		}
		return value.StringVal(strings.TrimSuffix(args[0].AsString(), args[1].AsString())), nil
	})

	r.Register("replace", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 3, "replace"); err != nil {
			return nil, err
		}
		n := -1
		if len(args) == 4 {
			n = int(args[3].AsInt())
		}
		return value.StringVal(strings.Replace(args[0].AsString(), args[1].AsString(), args[2].AsString(), n)), nil
	})

	r.Register("replaceAll", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 3, "replaceAll"); err != nil {
			return nil, err
		}
		return value.StringVal(strings.ReplaceAll(args[0].AsString(), args[1].AsString(), args[2].AsString())), nil
	})

	r.Register("startsWith", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "startsWith"); err != nil {
			return nil, err
		}
		return value.BoolVal(strings.HasPrefix(args[0].AsString(), args[1].AsString())), nil
	})

	r.Register("endsWith", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "endsWith"); err != nil {
			return nil, err
		}
		return value.BoolVal(strings.HasSuffix(args[0].AsString(), args[1].AsString())), nil
	})

	r.Register("split", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "split"); err != nil {
			return nil, err
		}
		parts := strings.Split(args[0].AsString(), args[1].AsString())
		items := make([]*value.Value, len(parts))
		for i, p := range parts {
			items[i] = value.StringVal(p)
		}
		return value.ListVal(items...), nil
	})

	r.Register("join", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "join"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		sep := args[1].AsString()
		parts := make([]string, len(list))
		for i, v := range list {
			parts[i] = v.AsString()
		}
		return value.StringVal(strings.Join(parts, sep)), nil
	})

	r.Register("len", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "len"); err != nil {
			return nil, err
		}
		return value.IntVal(int64(args[0].Len())), nil
	})

	r.Register("substr", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "substr"); err != nil {
			return nil, err
		}
		runes := []rune(args[0].AsString())
		start := int(args[1].AsInt())
		if start < 0 {
			start = len(runes) + start
		}
		if start < 0 {
			start = 0
		}
		if start > len(runes) {
			return value.StringVal(""), nil
		}
		if len(args) == 3 {
			end := int(args[2].AsInt())
			if end < 0 {
				end = len(runes) + end
			}
			if end > len(runes) {
				end = len(runes)
			}
			if end < start {
				return value.StringVal(""), nil
			}
			return value.StringVal(string(runes[start:end])), nil
		}
		return value.StringVal(string(runes[start:])), nil
	})

	r.Register("repeat", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "repeat"); err != nil {
			return nil, err
		}
		return value.StringVal(strings.Repeat(args[0].AsString(), int(args[1].AsInt()))), nil
	})

	r.Register("padLeft", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "padLeft"); err != nil {
			return nil, err
		}
		s := args[0].AsString()
		width := int(args[1].AsInt())
		pad := " "
		if len(args) == 3 {
			pad = args[2].AsString()
		}
		runes := []rune(s)
		for utf8.RuneCountInString(s) < width {
			s = pad + s
			_ = runes
		}
		return value.StringVal(s), nil
	})

	r.Register("padRight", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "padRight"); err != nil {
			return nil, err
		}
		s := args[0].AsString()
		width := int(args[1].AsInt())
		pad := " "
		if len(args) == 3 {
			pad = args[2].AsString()
		}
		for utf8.RuneCountInString(s) < width {
			s = s + pad
		}
		return value.StringVal(s), nil
	})

	r.Register("format", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "format"); err != nil {
			return nil, err
		}
		template := args[0].AsString()
		fmtArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			switch {
			case a.IsInt():
				fmtArgs[i] = a.AsInt()
			case a.IsFloat():
				fmtArgs[i] = a.AsFloat()
			case a.IsBool():
				fmtArgs[i] = a.AsBool()
			default:
				fmtArgs[i] = a.AsString()
			}
		}
		return value.StringVal(fmt.Sprintf(template, fmtArgs...)), nil
	})

	r.Register("charAt", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "charAt"); err != nil {
			return nil, err
		}
		runes := []rune(args[0].AsString())
		idx := int(args[1].AsInt())
		if idx < 0 {
			idx = len(runes) + idx
		}
		if idx < 0 || idx >= len(runes) {
			return value.StringVal(""), nil
		}
		return value.StringVal(string(runes[idx])), nil
	})

	r.Register("reverse", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "reverse"); err != nil {
			return nil, err
		}
		if args[0].IsString() {
			runes := []rune(args[0].AsString())
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return value.StringVal(string(runes)), nil
		}
		if args[0].IsList() {
			list := args[0].AsList()
			reversed := make([]*value.Value, len(list))
			for i, v := range list {
				reversed[len(list)-1-i] = v
			}
			return value.ListVal(reversed...), nil
		}
		return args[0], nil
	})
}
