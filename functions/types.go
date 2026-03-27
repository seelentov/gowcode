package functions

import "gowcode/value"

func registerTypeFuncs(r *Registry) {
	r.Register("toString", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "toString"); err != nil {
			return nil, err
		}
		return value.StringVal(args[0].AsString()), nil
	})

	r.Register("toInt", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "toInt"); err != nil {
			return nil, err
		}
		return value.IntVal(args[0].AsInt()), nil
	})

	r.Register("toFloat", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "toFloat"); err != nil {
			return nil, err
		}
		return value.FloatVal(args[0].AsFloat()), nil
	})

	r.Register("toBool", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "toBool"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].AsBool()), nil
	})

	// toList wraps a non-list value in a list; null becomes an empty list
	r.Register("toList", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "toList"); err != nil {
			return nil, err
		}
		if args[0].IsList() {
			return args[0], nil
		}
		if args[0].IsNull() {
			return value.ListVal(), nil
		}
		return value.ListVal(args[0]), nil
	})

	// typeOf returns the type name as a string
	typeNames := map[value.Type]string{
		value.Null:   "null",
		value.Bool:   "bool",
		value.Int:    "int",
		value.Float:  "float",
		value.String: "string",
		value.List:   "list",
		value.Map:    "map",
	}
	r.Register("typeOf", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "typeOf"); err != nil {
			return nil, err
		}
		if name, ok := typeNames[args[0].Type()]; ok {
			return value.StringVal(name), nil
		}
		return value.StringVal("unknown"), nil
	})

	r.Register("isString", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isString"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsString()), nil
	})

	r.Register("isInt", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isInt"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsInt()), nil
	})

	r.Register("isFloat", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isFloat"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsFloat()), nil
	})

	r.Register("isBool", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isBool"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsBool()), nil
	})

	r.Register("isList", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isList"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsList()), nil
	})

	r.Register("isMap", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isMap"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsMap()), nil
	})

	r.Register("isNumber", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isNumber"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsNumber()), nil
	})

	// isNull overrides the one in logic.go (types is registered last)
	r.Register("isNull", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isNull"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].IsNull()), nil
	})
}
