package ast

import "fmt"

// Node is the base interface for all AST nodes
type Node interface {
	nodeType() string
	String() string
}

// --- Literals ---

type IntLit struct{ Value int64 }
type FloatLit struct{ Value float64 }
type StringLit struct{ Value string }
type BoolLit struct{ Value bool }
type NullLit struct{}

func (n *IntLit) nodeType() string    { return "IntLit" }
func (n *FloatLit) nodeType() string  { return "FloatLit" }
func (n *StringLit) nodeType() string { return "StringLit" }
func (n *BoolLit) nodeType() string   { return "BoolLit" }
func (n *NullLit) nodeType() string   { return "NullLit" }

func (n *IntLit) String() string    { return fmt.Sprintf("%d", n.Value) }
func (n *FloatLit) String() string  { return fmt.Sprintf("%g", n.Value) }
func (n *StringLit) String() string { return fmt.Sprintf("%q", n.Value) }
func (n *BoolLit) String() string {
	if n.Value {
		return "true"
	}
	return "false"
}
func (n *NullLit) String() string { return "null" }

// --- List literal [a, b, c] ---

type ListLit struct {
	Items []Node
}

func (n *ListLit) nodeType() string { return "ListLit" }
func (n *ListLit) String() string   { return "List[...]" }

// --- Map literal {key: val, ...} ---

type MapEntry struct {
	Key   Node
	Value Node
}

type MapLit struct {
	Entries []MapEntry
}

func (n *MapLit) nodeType() string { return "MapLit" }
func (n *MapLit) String() string   { return "Map{...}" }

// --- Identifier ---

type Ident struct{ Name string }

func (n *Ident) nodeType() string { return "Ident" }
func (n *Ident) String() string   { return n.Name }

// --- Binary operation ---

type BinOp struct {
	Op    string
	Left  Node
	Right Node
}

func (n *BinOp) nodeType() string { return "BinOp" }
func (n *BinOp) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left, n.Op, n.Right)
}

// --- Unary operation ---

type UnaryOp struct {
	Op      string
	Operand Node
}

func (n *UnaryOp) nodeType() string { return "UnaryOp" }
func (n *UnaryOp) String() string {
	return fmt.Sprintf("(%s%s)", n.Op, n.Operand)
}

// --- Function call ---

type Call struct {
	Name string
	Args []Node
}

func (n *Call) nodeType() string { return "Call" }
func (n *Call) String() string   { return fmt.Sprintf("%s(...)", n.Name) }

// --- Method call: value.method(args) ---

type MethodCall struct {
	Object Node
	Method string
	Args   []Node
}

func (n *MethodCall) nodeType() string { return "MethodCall" }
func (n *MethodCall) String() string {
	return fmt.Sprintf("%s.%s(...)", n.Object, n.Method)
}

// --- Index access: value[index] ---

type Index struct {
	Object Node
	Index  Node
}

func (n *Index) nodeType() string { return "Index" }
func (n *Index) String() string {
	return fmt.Sprintf("%s[%s]", n.Object, n.Index)
}

// --- Ternary: condition ? then : else ---

type Ternary struct {
	Condition Node
	Then      Node
	Else      Node
}

func (n *Ternary) nodeType() string { return "Ternary" }
func (n *Ternary) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", n.Condition, n.Then, n.Else)
}

// --- Pipeline: a | b | c (b and c are function calls) ---

type Pipeline struct {
	Steps []Node
}

func (n *Pipeline) nodeType() string { return "Pipeline" }
func (n *Pipeline) String() string   { return "Pipeline(...)" }
