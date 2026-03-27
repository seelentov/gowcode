package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	// Literals
	TokInt    TokenType = iota
	TokFloat
	TokString
	TokBool
	TokNull
	TokIdent

	// Operators
	TokPlus
	TokMinus
	TokStar
	TokSlash
	TokPercent
	TokPower    // **
	TokEq       // ==
	TokNotEq    // !=
	TokLt       // <
	TokLtEq     // <=
	TokGt       // >
	TokGtEq     // >=
	TokAnd      // &&
	TokOr       // ||
	TokNot      // !
	TokDot      // .
	TokQuestion // ?
	TokColon    // :

	// Punctuation
	TokLParen   // (
	TokRParen   // )
	TokLBracket // [
	TokRBracket // ]
	TokLBrace   // {
	TokRBrace   // }
	TokComma    // ,
	TokSemicolon // ;

	TokEOF
)

var tokenNames = map[TokenType]string{
	TokInt: "INT", TokFloat: "FLOAT", TokString: "STRING",
	TokBool: "BOOL", TokNull: "NULL", TokIdent: "IDENT",
	TokPlus: "+", TokMinus: "-", TokStar: "*", TokSlash: "/",
	TokPercent: "%", TokPower: "**", TokEq: "==", TokNotEq: "!=",
	TokLt: "<", TokLtEq: "<=", TokGt: ">", TokGtEq: ">=",
	TokAnd: "&&", TokOr: "||", TokNot: "!", TokDot: ".",
	TokQuestion: "?", TokColon: ":",
	TokLParen: "(", TokRParen: ")", TokLBracket: "[", TokRBracket: "]",
	TokLBrace: "{", TokRBrace: "}", TokComma: ",", TokSemicolon: ";",
	TokEOF: "EOF",
}

func (t TokenType) String() string {
	if s, ok := tokenNames[t]; ok {
		return s
	}
	return fmt.Sprintf("tok(%d)", int(t))
}

type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Col     int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)", t.Type, t.Value)
}

type Lexer struct {
	input  []rune
	pos    int
	line   int
	col    int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input), pos: 0, line: 1, col: 1}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekAt(offset int) rune {
	i := l.pos + offset
	if i >= len(l.input) {
		return 0
	}
	return l.input[i]
}

func (l *Lexer) advance() rune {
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(l.peek()) {
		l.advance()
	}
}

func (l *Lexer) skipLineComment() {
	for l.pos < len(l.input) && l.peek() != '\n' {
		l.advance()
	}
}

func (l *Lexer) readString(quote rune) (Token, error) {
	line, col := l.line, l.col
	l.advance() // consume opening quote
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == quote {
			l.advance()
			return Token{TokString, sb.String(), line, col}, nil
		}
		if ch == '\\' {
			l.advance()
			esc := l.advance()
			switch esc {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '\\':
				sb.WriteByte('\\')
			case '"':
				sb.WriteByte('"')
			case '\'':
				sb.WriteByte('\'')
			default:
				sb.WriteByte('\\')
				sb.WriteRune(esc)
			}
			continue
		}
		sb.WriteRune(l.advance())
	}
	return Token{}, fmt.Errorf("line %d:%d: unterminated string", line, col)
}

func (l *Lexer) readNumber() Token {
	line, col := l.line, l.col
	var sb strings.Builder
	isFloat := false
	for l.pos < len(l.input) && (unicode.IsDigit(l.peek()) || l.peek() == '_') {
		ch := l.advance()
		if ch != '_' {
			sb.WriteRune(ch)
		}
	}
	if l.pos < len(l.input) && l.peek() == '.' && unicode.IsDigit(l.peekAt(1)) {
		isFloat = true
		sb.WriteRune(l.advance()) // .
		for l.pos < len(l.input) && unicode.IsDigit(l.peek()) {
			sb.WriteRune(l.advance())
		}
	}
	// scientific notation
	if l.pos < len(l.input) && (l.peek() == 'e' || l.peek() == 'E') {
		isFloat = true
		sb.WriteRune(l.advance())
		if l.pos < len(l.input) && (l.peek() == '+' || l.peek() == '-') {
			sb.WriteRune(l.advance())
		}
		for l.pos < len(l.input) && unicode.IsDigit(l.peek()) {
			sb.WriteRune(l.advance())
		}
	}
	if isFloat {
		return Token{TokFloat, sb.String(), line, col}
	}
	return Token{TokInt, sb.String(), line, col}
}

func (l *Lexer) readIdent() Token {
	line, col := l.line, l.col
	var sb strings.Builder
	for l.pos < len(l.input) && (unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_') {
		sb.WriteRune(l.advance())
	}
	word := sb.String()
	switch word {
	case "true", "false":
		return Token{TokBool, word, line, col}
	case "null", "nil":
		return Token{TokNull, word, line, col}
	}
	return Token{TokIdent, word, line, col}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			tokens = append(tokens, Token{TokEOF, "", l.line, l.col})
			break
		}

		line, col := l.line, l.col
		ch := l.peek()

		// Line comment
		if ch == '/' && l.peekAt(1) == '/' {
			l.skipLineComment()
			continue
		}

		// String literals
		if ch == '"' || ch == '\'' || ch == '`' {
			tok, err := l.readString(ch)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
			continue
		}

		// Numbers
		if unicode.IsDigit(ch) || (ch == '-' && unicode.IsDigit(l.peekAt(1)) && (len(tokens) == 0 || isOperatorOrParen(tokens[len(tokens)-1].Type))) {
			if ch == '-' {
				l.advance()
				tok := l.readNumber()
				tok.Value = "-" + tok.Value
				tok.Line, tok.Col = line, col
				tokens = append(tokens, tok)
				continue
			}
			tokens = append(tokens, l.readNumber())
			continue
		}

		// Identifiers and keywords
		if unicode.IsLetter(ch) || ch == '_' {
			tokens = append(tokens, l.readIdent())
			continue
		}

		// Two-char operators
		next := l.peekAt(1)
		switch {
		case ch == '*' && next == '*':
			l.advance(); l.advance()
			tokens = append(tokens, Token{TokPower, "**", line, col})
		case ch == '=' && next == '=':
			l.advance(); l.advance()
			tokens = append(tokens, Token{TokEq, "==", line, col})
		case ch == '!' && next == '=':
			l.advance(); l.advance()
			tokens = append(tokens, Token{TokNotEq, "!=", line, col})
		case ch == '<' && next == '=':
			l.advance(); l.advance()
			tokens = append(tokens, Token{TokLtEq, "<=", line, col})
		case ch == '>' && next == '=':
			l.advance(); l.advance()
			tokens = append(tokens, Token{TokGtEq, ">=", line, col})
		case ch == '&' && next == '&':
			l.advance(); l.advance()
			tokens = append(tokens, Token{TokAnd, "&&", line, col})
		case ch == '|' && next == '|':
			l.advance(); l.advance()
			tokens = append(tokens, Token{TokOr, "||", line, col})
		default:
			l.advance()
			switch ch {
			case '+':
				tokens = append(tokens, Token{TokPlus, "+", line, col})
			case '-':
				tokens = append(tokens, Token{TokMinus, "-", line, col})
			case '*':
				tokens = append(tokens, Token{TokStar, "*", line, col})
			case '/':
				tokens = append(tokens, Token{TokSlash, "/", line, col})
			case '%':
				tokens = append(tokens, Token{TokPercent, "%", line, col})
			case '<':
				tokens = append(tokens, Token{TokLt, "<", line, col})
			case '>':
				tokens = append(tokens, Token{TokGt, ">", line, col})
			case '!':
				tokens = append(tokens, Token{TokNot, "!", line, col})
			case '.':
				tokens = append(tokens, Token{TokDot, ".", line, col})
			case '?':
				tokens = append(tokens, Token{TokQuestion, "?", line, col})
			case ':':
				tokens = append(tokens, Token{TokColon, ":", line, col})
			case '(':
				tokens = append(tokens, Token{TokLParen, "(", line, col})
			case ')':
				tokens = append(tokens, Token{TokRParen, ")", line, col})
			case '[':
				tokens = append(tokens, Token{TokLBracket, "[", line, col})
			case ']':
				tokens = append(tokens, Token{TokRBracket, "]", line, col})
			case '{':
				tokens = append(tokens, Token{TokLBrace, "{", line, col})
			case '}':
				tokens = append(tokens, Token{TokRBrace, "}", line, col})
			case ',':
				tokens = append(tokens, Token{TokComma, ",", line, col})
			case ';':
				tokens = append(tokens, Token{TokSemicolon, ";", line, col})
			default:
				return nil, fmt.Errorf("line %d:%d: unexpected character %q", line, col, ch)
			}
		}
	}
	return tokens, nil
}

func isOperatorOrParen(t TokenType) bool {
	switch t {
	case TokPlus, TokMinus, TokStar, TokSlash, TokPercent, TokPower,
		TokEq, TokNotEq, TokLt, TokLtEq, TokGt, TokGtEq,
		TokAnd, TokOr, TokNot, TokLParen, TokLBracket, TokComma,
		TokColon, TokQuestion:
		return true
	}
	return false
}
