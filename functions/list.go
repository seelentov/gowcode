package functions

import (
	"fmt"
	"sort"
	"strings"

	"gowcode/value"
)

func registerListFuncs(r *Registry) {
	r.Register("append", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "append"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		newList := make([]*value.Value, len(list)+1)
		copy(newList, list)
		newList[len(list)] = args[1]
		return value.ListVal(newList...), nil
	})

	r.Register("prepend", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "prepend"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		newList := make([]*value.Value, len(list)+1)
		newList[0] = args[1]
		copy(newList[1:], list)
		return value.ListVal(newList...), nil
	})

	r.Register("concat", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "concat"); err != nil {
			return nil, err
		}
		var result []*value.Value
		for _, a := range args {
			result = append(result, a.AsList()...)
		}
		return value.ListVal(result...), nil
	})

	r.Register("first", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "first"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		if len(args) == 1 {
			if len(list) == 0 {
				return value.Nil(), nil
			}
			return list[0], nil
		}
		n := int(args[1].AsInt())
		if n > len(list) {
			n = len(list)
		}
		return value.ListVal(list[:n]...), nil
	})

	r.Register("last", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "last"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		if len(args) == 1 {
			if len(list) == 0 {
				return value.Nil(), nil
			}
			return list[len(list)-1], nil
		}
		n := int(args[1].AsInt())
		if n > len(list) {
			n = len(list)
		}
		return value.ListVal(list[len(list)-n:]...), nil
	})

	r.Register("nth", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "nth"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		idx := int(args[1].AsInt())
		if idx < 0 {
			idx = len(list) + idx
		}
		if idx < 0 || idx >= len(list) {
			return value.Nil(), nil
		}
		return list[idx], nil
	})

	r.Register("slice", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "slice"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		start := int(args[1].AsInt())
		if start < 0 {
			start = len(list) + start
		}
		if start < 0 {
			start = 0
		}
		if start > len(list) {
			return value.ListVal(), nil
		}
		if len(args) == 3 {
			end := int(args[2].AsInt())
			if end < 0 {
				end = len(list) + end
			}
			if end > len(list) {
				end = len(list)
			}
			if end < start {
				return value.ListVal(), nil
			}
			return value.ListVal(list[start:end]...), nil
		}
		return value.ListVal(list[start:]...), nil
	})

	r.Register("take", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "take"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		n := int(args[1].AsInt())
		if n < 0 {
			n = 0
		}
		if n > len(list) {
			n = len(list)
		}
		return value.ListVal(list[:n]...), nil
	})

	r.Register("drop", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "drop"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		n := int(args[1].AsInt())
		if n < 0 {
			n = 0
		}
		if n > len(list) {
			n = len(list)
		}
		return value.ListVal(list[n:]...), nil
	})

	// Override contains from strings.go to handle both strings and lists
	r.Register("contains", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "contains"); err != nil {
			return nil, err
		}
		if args[0].IsList() {
			for _, item := range args[0].AsList() {
				if item.Equal(args[1]) {
					return value.BoolVal(true), nil
				}
			}
			return value.BoolVal(false), nil
		}
		return value.BoolVal(strings.Contains(args[0].AsString(), args[1].AsString())), nil
	})

	// Override indexOf from strings.go to handle both strings and lists
	r.Register("indexOf", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "indexOf"); err != nil {
			return nil, err
		}
		if args[0].IsList() {
			for i, item := range args[0].AsList() {
				if item.Equal(args[1]) {
					return value.IntVal(int64(i)), nil
				}
			}
			return value.IntVal(-1), nil
		}
		return value.IntVal(int64(strings.Index(args[0].AsString(), args[1].AsString()))), nil
	})

	// Override lastIndexOf from strings.go to handle both strings and lists
	r.Register("lastIndexOf", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "lastIndexOf"); err != nil {
			return nil, err
		}
		if args[0].IsList() {
			list := args[0].AsList()
			for i := len(list) - 1; i >= 0; i-- {
				if list[i].Equal(args[1]) {
					return value.IntVal(int64(i)), nil
				}
			}
			return value.IntVal(-1), nil
		}
		return value.IntVal(int64(strings.LastIndex(args[0].AsString(), args[1].AsString()))), nil
	})

	r.Register("flatten", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "flatten"); err != nil {
			return nil, err
		}
		var result []*value.Value
		for _, item := range args[0].AsList() {
			if item.IsList() {
				result = append(result, item.AsList()...)
			} else {
				result = append(result, item)
			}
		}
		return value.ListVal(result...), nil
	})

	r.Register("flattenAll", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "flattenAll"); err != nil {
			return nil, err
		}
		return value.ListVal(flattenDeep(args[0].AsList())...), nil
	})

	r.Register("unique", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "unique"); err != nil {
			return nil, err
		}
		var result []*value.Value
		for _, item := range args[0].AsList() {
			found := false
			for _, existing := range result {
				if existing.Equal(item) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, item)
			}
		}
		return value.ListVal(result...), nil
	})

	r.Register("sort", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "sort"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		sorted := make([]*value.Value, len(list))
		copy(sorted, list)
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].Compare(sorted[j]) < 0
		})
		return value.ListVal(sorted...), nil
	})

	r.Register("sortDesc", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "sortDesc"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		sorted := make([]*value.Value, len(list))
		copy(sorted, list)
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].Compare(sorted[j]) > 0
		})
		return value.ListVal(sorted...), nil
	})

	// range(end) or range(start, end) or range(start, end, step)
	r.Register("range", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "range"); err != nil {
			return nil, err
		}
		var start, end, step int64
		switch len(args) {
		case 1:
			start, end, step = 0, args[0].AsInt(), 1
		case 2:
			start, end, step = args[0].AsInt(), args[1].AsInt(), 1
		default:
			start, end, step = args[0].AsInt(), args[1].AsInt(), args[2].AsInt()
		}
		if step == 0 {
			return nil, fmt.Errorf("range: step cannot be zero")
		}
		var result []*value.Value
		if step > 0 {
			for i := start; i < end; i += step {
				result = append(result, value.IntVal(i))
			}
		} else {
			for i := start; i > end; i += step {
				result = append(result, value.IntVal(i))
			}
		}
		return value.ListVal(result...), nil
	})

	r.Register("chunk", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "chunk"); err != nil {
			return nil, err
		}
		list := args[0].AsList()
		size := int(args[1].AsInt())
		if size <= 0 {
			return nil, fmt.Errorf("chunk: size must be positive")
		}
		var result []*value.Value
		for i := 0; i < len(list); i += size {
			end := i + size
			if end > len(list) {
				end = len(list)
			}
			result = append(result, value.ListVal(list[i:end]...))
		}
		return value.ListVal(result...), nil
	})

	r.Register("zip", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "zip"); err != nil {
			return nil, err
		}
		a, b := args[0].AsList(), args[1].AsList()
		n := len(a)
		if len(b) < n {
			n = len(b)
		}
		result := make([]*value.Value, n)
		for i := 0; i < n; i++ {
			result[i] = value.ListVal(a[i], b[i])
		}
		return value.ListVal(result...), nil
	})

	// without(list, item1, item2, ...) — returns list without specified items
	r.Register("without", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "without"); err != nil {
			return nil, err
		}
		exclude := args[1:]
		var result []*value.Value
		for _, item := range args[0].AsList() {
			found := false
			for _, ex := range exclude {
				if item.Equal(ex) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, item)
			}
		}
		return value.ListVal(result...), nil
	})

	r.Register("count", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "count"); err != nil {
			return nil, err
		}
		return value.IntVal(int64(args[0].Len())), nil
	})
}

func flattenDeep(list []*value.Value) []*value.Value {
	var result []*value.Value
	for _, item := range list {
		if item.IsList() {
			result = append(result, flattenDeep(item.AsList())...)
		} else {
			result = append(result, item)
		}
	}
	return result
}
