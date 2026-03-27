package functions

import "gowcode/value"

func registerMapFuncs(r *Registry) {
	r.Register("keys", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "keys"); err != nil {
			return nil, err
		}
		ks := args[0].Keys()
		result := make([]*value.Value, len(ks))
		for i, k := range ks {
			result[i] = value.StringVal(k)
		}
		return value.ListVal(result...), nil
	})

	r.Register("values", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "values"); err != nil {
			return nil, err
		}
		ks := args[0].Keys()
		result := make([]*value.Value, len(ks))
		for i, k := range ks {
			result[i] = args[0].GetKey(k)
		}
		return value.ListVal(result...), nil
	})

	r.Register("hasKey", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "hasKey"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].HasKey(args[1].AsString())), nil
	})

	// get(map, key, default?) — returns default if key is missing
	r.Register("get", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "get"); err != nil {
			return nil, err
		}
		v := args[0].GetKey(args[1].AsString())
		if v.IsNull() && len(args) == 3 {
			return args[2], nil
		}
		return v, nil
	})

	// set(map, key, value) — returns a new map with the key set
	r.Register("set", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 3, "set"); err != nil {
			return nil, err
		}
		result := value.MapVal()
		for _, k := range args[0].Keys() {
			result.Set(k, args[0].GetKey(k))
		}
		result.Set(args[1].AsString(), args[2])
		return result, nil
	})

	// delete(map, key1, key2, ...) — returns a new map without the given keys
	r.Register("delete", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "delete"); err != nil {
			return nil, err
		}
		toDelete := make(map[string]bool, len(args)-1)
		for _, a := range args[1:] {
			toDelete[a.AsString()] = true
		}
		result := value.MapVal()
		for _, k := range args[0].Keys() {
			if !toDelete[k] {
				result.Set(k, args[0].GetKey(k))
			}
		}
		return result, nil
	})

	// merge(map1, map2, ...) — later keys win
	r.Register("merge", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "merge"); err != nil {
			return nil, err
		}
		result := value.MapVal()
		for _, a := range args {
			for _, k := range a.Keys() {
				result.Set(k, a.GetKey(k))
			}
		}
		return result, nil
	})

	// entries(map) — returns [[key, value], ...]
	r.Register("entries", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "entries"); err != nil {
			return nil, err
		}
		ks := args[0].Keys()
		result := make([]*value.Value, len(ks))
		for i, k := range ks {
			result[i] = value.ListVal(value.StringVal(k), args[0].GetKey(k))
		}
		return value.ListVal(result...), nil
	})

	// fromEntries([[key, value], ...]) — builds a map from entry pairs
	r.Register("fromEntries", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "fromEntries"); err != nil {
			return nil, err
		}
		result := value.MapVal()
		for _, entry := range args[0].AsList() {
			items := entry.AsList()
			if len(items) >= 2 {
				result.Set(items[0].AsString(), items[1])
			}
		}
		return result, nil
	})

	// pick(map, key1, key2, ...) — returns a new map with only specified keys
	r.Register("pick", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "pick"); err != nil {
			return nil, err
		}
		result := value.MapVal()
		for _, a := range args[1:] {
			k := a.AsString()
			if args[0].HasKey(k) {
				result.Set(k, args[0].GetKey(k))
			}
		}
		return result, nil
	})

	// omit(map, key1, key2, ...) — returns a new map without specified keys
	r.Register("omit", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "omit"); err != nil {
			return nil, err
		}
		toOmit := make(map[string]bool, len(args)-1)
		for _, a := range args[1:] {
			toOmit[a.AsString()] = true
		}
		result := value.MapVal()
		for _, k := range args[0].Keys() {
			if !toOmit[k] {
				result.Set(k, args[0].GetKey(k))
			}
		}
		return result, nil
	})
}
