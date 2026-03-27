package functions

import (
	"math"
	"math/rand"

	"gowcode/value"
)

func registerNumberFuncs(r *Registry) {
	r.Register("abs", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "abs"); err != nil {
			return nil, err
		}
		if args[0].IsInt() {
			v := args[0].AsInt()
			if v < 0 {
				return value.IntVal(-v), nil
			}
			return value.IntVal(v), nil
		}
		return value.FloatVal(math.Abs(args[0].AsFloat())), nil
	})

	r.Register("floor", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "floor"); err != nil {
			return nil, err
		}
		return value.IntVal(int64(math.Floor(args[0].AsFloat()))), nil
	})

	r.Register("ceil", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "ceil"); err != nil {
			return nil, err
		}
		return value.IntVal(int64(math.Ceil(args[0].AsFloat()))), nil
	})

	r.Register("round", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "round"); err != nil {
			return nil, err
		}
		f := args[0].AsFloat()
		if len(args) == 2 {
			factor := math.Pow(10, float64(args[1].AsInt()))
			return value.FloatVal(math.Round(f*factor) / factor), nil
		}
		return value.IntVal(int64(math.Round(f))), nil
	})

	r.Register("trunc", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "trunc"); err != nil {
			return nil, err
		}
		return value.IntVal(int64(math.Trunc(args[0].AsFloat()))), nil
	})

	r.Register("sqrt", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "sqrt"); err != nil {
			return nil, err
		}
		return value.FloatVal(math.Sqrt(args[0].AsFloat())), nil
	})

	r.Register("pow", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "pow"); err != nil {
			return nil, err
		}
		return value.FloatVal(math.Pow(args[0].AsFloat(), args[1].AsFloat())), nil
	})

	r.Register("sign", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "sign"); err != nil {
			return nil, err
		}
		f := args[0].AsFloat()
		if f < 0 {
			return value.IntVal(-1), nil
		} else if f > 0 {
			return value.IntVal(1), nil
		}
		return value.IntVal(0), nil
	})

	r.Register("log", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "log"); err != nil {
			return nil, err
		}
		n := args[0].AsFloat()
		if len(args) == 2 {
			base := args[1].AsFloat()
			return value.FloatVal(math.Log(n) / math.Log(base)), nil
		}
		return value.FloatVal(math.Log(n)), nil
	})

	r.Register("log2", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "log2"); err != nil {
			return nil, err
		}
		return value.FloatVal(math.Log2(args[0].AsFloat())), nil
	})

	r.Register("log10", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "log10"); err != nil {
			return nil, err
		}
		return value.FloatVal(math.Log10(args[0].AsFloat())), nil
	})

	r.Register("pi", func(args []*value.Value) (*value.Value, error) {
		return value.FloatVal(math.Pi), nil
	})

	r.Register("e", func(args []*value.Value) (*value.Value, error) {
		return value.FloatVal(math.E), nil
	})

	r.Register("isNaN", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isNaN"); err != nil {
			return nil, err
		}
		return value.BoolVal(math.IsNaN(args[0].AsFloat())), nil
	})

	r.Register("isInf", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 1, "isInf"); err != nil {
			return nil, err
		}
		return value.BoolVal(math.IsInf(args[0].AsFloat(), 0)), nil
	})

	r.Register("min", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "min"); err != nil {
			return nil, err
		}
		nums := numericArgs(args)
		if len(nums) == 0 {
			return value.Nil(), nil
		}
		m := nums[0]
		for _, v := range nums[1:] {
			if v.Compare(m) < 0 {
				m = v
			}
		}
		return m, nil
	})

	r.Register("max", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "max"); err != nil {
			return nil, err
		}
		nums := numericArgs(args)
		if len(nums) == 0 {
			return value.Nil(), nil
		}
		m := nums[0]
		for _, v := range nums[1:] {
			if v.Compare(m) > 0 {
				m = v
			}
		}
		return m, nil
	})

	r.Register("clamp", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 3, "clamp"); err != nil {
			return nil, err
		}
		n, lo, hi := args[0].AsFloat(), args[1].AsFloat(), args[2].AsFloat()
		if n < lo {
			n = lo
		} else if n > hi {
			n = hi
		}
		if args[0].IsInt() && args[1].IsInt() && args[2].IsInt() {
			return value.IntVal(int64(n)), nil
		}
		return value.FloatVal(n), nil
	})

	r.Register("sum", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "sum"); err != nil {
			return nil, err
		}
		nums := numericArgs(args)
		allInt := true
		var isum int64
		var fsum float64
		for _, v := range nums {
			if !v.IsInt() {
				allInt = false
			}
			isum += v.AsInt()
			fsum += v.AsFloat()
		}
		if allInt {
			return value.IntVal(isum), nil
		}
		return value.FloatVal(fsum), nil
	})

	r.Register("avg", func(args []*value.Value) (*value.Value, error) {
		if err := argsMin(args, 1, "avg"); err != nil {
			return nil, err
		}
		nums := numericArgs(args)
		if len(nums) == 0 {
			return value.Nil(), nil
		}
		var fsum float64
		for _, v := range nums {
			fsum += v.AsFloat()
		}
		return value.FloatVal(fsum / float64(len(nums))), nil
	})

	r.Register("random", func(args []*value.Value) (*value.Value, error) {
		return value.FloatVal(rand.Float64()), nil
	})

	r.Register("randomInt", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "randomInt"); err != nil {
			return nil, err
		}
		lo, hi := args[0].AsInt(), args[1].AsInt()
		if lo >= hi {
			return value.IntVal(lo), nil
		}
		return value.IntVal(lo + rand.Int63n(hi-lo)), nil
	})

	r.Register("gcd", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "gcd"); err != nil {
			return nil, err
		}
		a, b := args[0].AsInt(), args[1].AsInt()
		if a < 0 {
			a = -a
		}
		if b < 0 {
			b = -b
		}
		for b != 0 {
			a, b = b, a%b
		}
		return value.IntVal(a), nil
	})

	r.Register("lcm", func(args []*value.Value) (*value.Value, error) {
		if err := argsExact(args, 2, "lcm"); err != nil {
			return nil, err
		}
		a, b := args[0].AsInt(), args[1].AsInt()
		if a == 0 || b == 0 {
			return value.IntVal(0), nil
		}
		if a < 0 {
			a = -a
		}
		if b < 0 {
			b = -b
		}
		g, tmp := a, b
		for tmp != 0 {
			g, tmp = tmp, g%tmp
		}
		return value.IntVal(a / g * b), nil
	})
}

// numericArgs expands a single list argument into its elements, or returns args as-is.
func numericArgs(args []*value.Value) []*value.Value {
	if len(args) == 1 && args[0].IsList() {
		return args[0].AsList()
	}
	return args
}
