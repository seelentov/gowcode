package parser

import (
	"fmt"
	"strconv"

	"gowcode/ast"
	"gowcode/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func Parse(input string) (ast.Node, error) {
	l := lexer.New(input)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, err
	}
	p := New(tokens)
	node, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.isEOF() {
		tok := p.peek()
		return nil, fmt.Errorf("line %d:%d: unexpected token %s", tok.Line, tok.Col, tok)
	}
	return node, nil
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) peekType() lexer.TokenType {
	return p.tokens[p.pos].Type
}

func (p *Parser) advance() lexer.Token {
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *Parser) isEOF() bool {
	return p.peekType() == lexer.TokEOF
}

func (p *Parser) expect(t lexer.TokenType) (lexer.Token, error) {
	tok := p.peek()
	if tok.Type != t {
		return tok, fmt.Errorf("line %d:%d: expected %s, got %s", tok.Line, tok.Col, t, tok.Type)
	}
	return p.advance(), nil
}

// parseExpr handles ternary and pipeline (lowest precedence)
func (p *Parser) parseExpr() (ast.Node, error) {
	return p.parseTernary()
}

func (p *Parser) parseTernary() (ast.Node, error) {
	node, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.peekType() == lexer.TokQuestion {
		p.advance()
		thenNode, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.TokColon); err != nil {
			return nil, err
		}
		elseNode, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		return &ast.Ternary{Condition: node, Then: thenNode, Else: elseNode}, nil
	}
	return node, nil
}

func (p *Parser) parseOr() (ast.Node, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.peekType() == lexer.TokOr {
		op := p.advance().Value
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAnd() (ast.Node, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}
	for p.peekType() == lexer.TokAnd {
		op := p.advance().Value
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseEquality() (ast.Node, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.peekType() == lexer.TokEq || p.peekType() == lexer.TokNotEq {
		op := p.advance().Value
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseComparison() (ast.Node, error) {
	left, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
	for p.peekType() == lexer.TokLt || p.peekType() == lexer.TokLtEq ||
		p.peekType() == lexer.TokGt || p.peekType() == lexer.TokGtEq {
		op := p.advance().Value
		right, err := p.parseAddSub()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAddSub() (ast.Node, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for p.peekType() == lexer.TokPlus || p.peekType() == lexer.TokMinus {
		op := p.advance().Value
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseMulDiv() (ast.Node, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	for p.peekType() == lexer.TokStar || p.peekType() == lexer.TokSlash || p.peekType() == lexer.TokPercent {
		op := p.advance().Value
		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parsePower() (ast.Node, error) {
	base, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.peekType() == lexer.TokPower {
		p.advance()
		exp, err := p.parsePower() // right-associative
		if err != nil {
			return nil, err
		}
		return &ast.BinOp{Op: "**", Left: base, Right: exp}, nil
	}
	return base, nil
}

func (p *Parser) parseUnary() (ast.Node, error) {
	if p.peekType() == lexer.TokNot {
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOp{Op: "!", Operand: operand}, nil
	}
	if p.peekType() == lexer.TokMinus {
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOp{Op: "-", Operand: operand}, nil
	}
	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (ast.Node, error) {
	node, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		switch p.peekType() {
		case lexer.TokDot:
			p.advance()
			nameTok, err := p.expect(lexer.TokIdent)
			if err != nil {
				return nil, err
			}
			if p.peekType() == lexer.TokLParen {
				args, err := p.parseArgList()
				if err != nil {
					return nil, err
				}
				node = &ast.MethodCall{Object: node, Method: nameTok.Value, Args: args}
			} else {
				// property access as method with no args
				node = &ast.MethodCall{Object: node, Method: nameTok.Value, Args: nil}
			}
		case lexer.TokLBracket:
			p.advance()
			idx, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(lexer.TokRBracket); err != nil {
				return nil, err
			}
			node = &ast.Index{Object: node, Index: idx}
		default:
			return node, nil
		}
	}
}

func (p *Parser) parsePrimary() (ast.Node, error) {
	tok := p.peek()

	switch tok.Type {
	case lexer.TokInt:
		p.advance()
		v, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid int %q: %w", tok.Value, err)
		}
		return &ast.IntLit{Value: v}, nil

	case lexer.TokFloat:
		p.advance()
		v, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float %q: %w", tok.Value, err)
		}
		return &ast.FloatLit{Value: v}, nil

	case lexer.TokString:
		p.advance()
		return &ast.StringLit{Value: tok.Value}, nil

	case lexer.TokBool:
		p.advance()
		return &ast.BoolLit{Value: tok.Value == "true"}, nil

	case lexer.TokNull:
		p.advance()
		return &ast.NullLit{}, nil

	case lexer.TokIdent:
		p.advance()
		if p.peekType() == lexer.TokLParen {
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			return &ast.Call{Name: tok.Value, Args: args}, nil
		}
		return &ast.Ident{Name: tok.Value}, nil

	case lexer.TokLParen:
		p.advance()
		node, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.TokRParen); err != nil {
			return nil, err
		}
		return node, nil

	case lexer.TokLBracket:
		return p.parseListLit()

	case lexer.TokLBrace:
		return p.parseMapLit()
	}

	return nil, fmt.Errorf("line %d:%d: unexpected token %s", tok.Line, tok.Col, tok.Type)
}

func (p *Parser) parseArgList() ([]ast.Node, error) {
	if _, err := p.expect(lexer.TokLParen); err != nil {
		return nil, err
	}
	var args []ast.Node
	for p.peekType() != lexer.TokRParen && !p.isEOF() {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.peekType() == lexer.TokComma {
			p.advance()
		} else {
			break
		}
	}
	if _, err := p.expect(lexer.TokRParen); err != nil {
		return nil, err
	}
	return args, nil
}

func (p *Parser) parseListLit() (ast.Node, error) {
	p.advance() // [
	var items []ast.Node
	for p.peekType() != lexer.TokRBracket && !p.isEOF() {
		item, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.peekType() == lexer.TokComma {
			p.advance()
		} else {
			break
		}
	}
	if _, err := p.expect(lexer.TokRBracket); err != nil {
		return nil, err
	}
	return &ast.ListLit{Items: items}, nil
}

func (p *Parser) parseMapLit() (ast.Node, error) {
	p.advance() // {
	var entries []ast.MapEntry
	for p.peekType() != lexer.TokRBrace && !p.isEOF() {
		key, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.TokColon); err != nil {
			return nil, err
		}
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		entries = append(entries, ast.MapEntry{Key: key, Value: val})
		if p.peekType() == lexer.TokComma {
			p.advance()
		} else {
			break
		}
	}
	if _, err := p.expect(lexer.TokRBrace); err != nil {
		return nil, err
	}
	return &ast.MapLit{Entries: entries}, nil
}
