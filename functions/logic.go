package functions

import "gowcode/value"

func registerLogicFuncs(r *Registry) {
	r.Register("not", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "not"); err != nil {
			return nil, err
		}
		return value.BoolVal(!args[0].AsBool()), nil
	})

	r.Register("and", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "and"); err != nil {
			return nil, err
		}
		for _, a := range args {
			if !a.AsBool() {
				return value.BoolVal(false), nil
			}
		}
		return value.BoolVal(true), nil
	})

	r.Register("or", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "or"); err != nil {
			return nil, err
		}
		for _, a := range args {
			if a.AsBool() {
				return value.BoolVal(true), nil
			}
		}
		return value.BoolVal(false), nil
	})

	r.Register("xor", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "xor"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].AsBool() != args[1].AsBool()), nil
	})

	// if(cond, then, else?) — note: both branches are eagerly evaluated
	r.Register("if", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 2, "if"); err != nil {
			return nil, err
		}
		if args[0].AsBool() {
			return args[1], nil
		}
		if len(args) == 3 {
			return args[2], nil
		}
		return value.Nil(), nil
	})

	// coalesce returns the first non-null value
	r.Register("coalesce", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "coalesce"); err != nil {
			return nil, err
		}
		for _, a := range args {
			if !a.IsNull() {
				return a, nil
			}
		}
		return value.Nil(), nil
	})

	// defaultTo returns the second arg if the first is null
	r.Register("defaultTo", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "defaultTo"); err != nil {
			return nil, err
		}
		if args[0].IsNull() {
			return args[1], nil
		}
		return args[0], nil
	})

	r.Register("isTruthy", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isTruthy"); err != nil {
			return nil, err
		}
		return value.BoolVal(args[0].AsBool()), nil
	})

	r.Register("isFalsy", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isFalsy"); err != nil {
			return nil, err
		}
		return value.BoolVal(!args[0].AsBool()), nil
	})
}
