package eval

import (
	"fmt"
	"math"

	"gowcode/ast"
	"gowcode/functions"
	"gowcode/parser"
	"gowcode/value"
)

type Evaluator struct {
	vars     map[string]*value.Value
	registry *functions.Registry
}

func NewEvaluator(vars map[string]*value.Value) *Evaluator {
	return &Evaluator{
		vars:     vars,
		registry: functions.NewRegistry(),
	}
}

func NewEvaluatorWithRegistry(vars map[string]*value.Value, registry *functions.Registry) *Evaluator {

	return &Evaluator{
		vars:     vars,
		registry: registry,
	}
}

// Eval parses and evaluates an expression string.
func Eval(expr string, vars map[string]*value.Value) (*value.Value, error) {
	node, err := parser.Parse(expr)
	if err != nil {
		return nil, err
	}
	return NewEvaluator(vars).eval(node)
}

func (e *Evaluator) eval(node ast.Node) (*value.Value, error) {
	switch n := node.(type) {
	case *ast.IntLit:
		return value.IntVal(n.Value), nil
	case *ast.FloatLit:
		return value.FloatVal(n.Value), nil
	case *ast.StringLit:
		return value.StringVal(n.Value), nil
	case *ast.BoolLit:
		return value.BoolVal(n.Value), nil
	case *ast.NullLit:
		return value.Nil(), nil

	case *ast.Ident:
		if v, ok := e.vars[n.Name]; ok {
			return v, nil
		}
		return nil, fmt.Errorf("undefined variable: %q", n.Name)

	case *ast.UnaryOp:
		operand, err := e.eval(n.Operand)
		if err != nil {
			return nil, err
		}
		switch n.Op {
		case "-":
			if operand.IsInt() {
				return value.IntVal(-operand.AsInt()), nil
			}
			return value.FloatVal(-operand.AsFloat()), nil
		case "!":
			return value.BoolVal(!operand.AsBool()), nil
		}

	case *ast.BinOp:
		left, err := e.eval(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.eval(n.Right)
		if err != nil {
			return nil, err
		}
		return evalBinOp(n.Op, left, right)

	case *ast.Ternary:
		cond, err := e.eval(n.Condition)
		if err != nil {
			return nil, err
		}
		if cond.AsBool() {
			return e.eval(n.Then)
		}
		return e.eval(n.Else)

	case *ast.Call:
		args := make([]*value.Value, len(n.Args))
		for i, a := range n.Args {
			v, err := e.eval(a)
			if err != nil {
				return nil, err
			}
			args[i] = v
		}
		return e.registry.Call(n.Name, args)

	case *ast.ListLit:
		items := make([]*value.Value, len(n.Items))
		for i, item := range n.Items {
			v, err := e.eval(item)
			if err != nil {
				return nil, err
			}
			items[i] = v
		}
		return value.ListVal(items...), nil

	case *ast.MapLit:
		m := value.MapVal()
		for _, entry := range n.Entries {
			k, err := e.eval(entry.Key)
			if err != nil {
				return nil, err
			}
			v, err := e.eval(entry.Value)
			if err != nil {
				return nil, err
			}
			m.Set(k.AsString(), v)
		}
		return m, nil

	case *ast.Index:
		obj, err := e.eval(n.Object)
		if err != nil {
			return nil, err
		}
		idx, err := e.eval(n.Index)
		if err != nil {
			return nil, err
		}
		if obj.IsList() {
			return obj.Get(int(idx.AsInt())), nil
		}
		if obj.IsMap() {
			return obj.GetKey(idx.AsString()), nil
		}
		return nil, fmt.Errorf("index on non-indexable type")
	}

	return nil, fmt.Errorf("unsupported node: %T", node)
}

func evalBinOp(op string, left, right *value.Value) (*value.Value, error) {
	switch op {
	case "==":
		return value.BoolVal(left.Equal(right)), nil
	case "!=":
		return value.BoolVal(!left.Equal(right)), nil
	case "&&":
		return value.BoolVal(left.AsBool() && right.AsBool()), nil
	case "||":
		return value.BoolVal(left.AsBool() || right.AsBool()), nil
	case "<":
		return value.BoolVal(left.Compare(right) < 0), nil
	case "<=":
		return value.BoolVal(left.Compare(right) <= 0), nil
	case ">":
		return value.BoolVal(left.Compare(right) > 0), nil
	case ">=":
		return value.BoolVal(left.Compare(right) >= 0), nil
	}

	// Arithmetic
	if left.IsInt() && right.IsInt() {
		a, b := left.AsInt(), right.AsInt()
		switch op {
		case "+":
			return value.IntVal(a + b), nil
		case "-":
			return value.IntVal(a - b), nil
		case "*":
			return value.IntVal(a * b), nil
		case "/":
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return value.IntVal(a / b), nil
		case "%":
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return value.IntVal(a % b), nil
		case "**":
			return value.FloatVal(math.Pow(float64(a), float64(b))), nil
		}
	}

	if left.IsNumber() && right.IsNumber() {
		a, b := left.AsFloat(), right.AsFloat()
		switch op {
		case "+":
			return value.FloatVal(a + b), nil
		case "-":
			return value.FloatVal(a - b), nil
		case "*":
			return value.FloatVal(a * b), nil
		case "/":
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return value.FloatVal(a / b), nil
		case "**":
			return value.FloatVal(math.Pow(a, b)), nil
		}
	}

	if op == "+" && left.IsString() {
		return value.StringVal(left.AsString() + right.AsString()), nil
	}

	return nil, fmt.Errorf("operator %q not supported for types", op)
}
