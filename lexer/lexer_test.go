package lexer

import (
	"testing"
)

// helper: tokenize and return only non-EOF tokens
func tokenize(t *testing.T, input string) []Token {
	t.Helper()
	l := New(input)
	toks, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize(%q) unexpected error: %v", input, err)
	}
	// strip trailing EOF
	if len(toks) > 0 && toks[len(toks)-1].Type == TokEOF {
		return toks[:len(toks)-1]
	}
	return toks
}

func tokenizeErr(t *testing.T, input string) {
	t.Helper()
	l := New(input)
	_, err := l.Tokenize()
	if err == nil {
		t.Fatalf("Tokenize(%q): expected error, got nil", input)
	}
}

func assertToken(t *testing.T, tok Token, wantType TokenType, wantVal string) {
	t.Helper()
	if tok.Type != wantType {
		t.Errorf("token type = %v, want %v (value=%q)", tok.Type, wantType, tok.Value)
	}
	if tok.Value != wantVal {
		t.Errorf("token value = %q, want %q", tok.Value, wantVal)
	}
}

// --- Integer literals ---

func TestLexer_IntLiteral(t *testing.T) {
	toks := tokenize(t, "42")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokInt, "42")
}

func TestLexer_NegativeInt(t *testing.T) {
	// negative int at start of input
	toks := tokenize(t, "-5")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d: %v", len(toks), toks)
	}
	assertToken(t, toks[0], TokInt, "-5")
}

func TestLexer_IntWithUnderscores(t *testing.T) {
	toks := tokenize(t, "1_000_000")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokInt, "1000000")
}

// --- Float literals ---

func TestLexer_FloatLiteral(t *testing.T) {
	toks := tokenize(t, "3.14")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokFloat, "3.14")
}

func TestLexer_NegativeFloat(t *testing.T) {
	toks := tokenize(t, "-3.14")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d: %v", len(toks), toks)
	}
	assertToken(t, toks[0], TokFloat, "-3.14")
}

func TestLexer_ScientificNotation(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"1e10", "1e10"},
		{"1.5E3", "1.5E3"},
		{"1E-3", "1E-3"},
		{"2.5e+2", "2.5e+2"},
	}
	for _, c := range cases {
		toks := tokenize(t, c.input)
		if len(toks) != 1 {
			t.Errorf("%s: expected 1 token, got %d", c.input, len(toks))
			continue
		}
		assertToken(t, toks[0], TokFloat, c.want)
	}
}

// --- String literals ---

func TestLexer_DoubleQuoteString(t *testing.T) {
	toks := tokenize(t, `"hello"`)
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokString, "hello")
}

func TestLexer_SingleQuoteString(t *testing.T) {
	toks := tokenize(t, "'world'")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokString, "world")
}

func TestLexer_BacktickString(t *testing.T) {
	toks := tokenize(t, "`raw`")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokString, "raw")
}

func TestLexer_StringEscapes(t *testing.T) {
	toks := tokenize(t, `"a\nb\tc\r\\"`)
	assertToken(t, toks[0], TokString, "a\nb\tc\r\\")
}

func TestLexer_StringEscapeQuotes(t *testing.T) {
	toks := tokenize(t, `"say \"hi\""`)
	assertToken(t, toks[0], TokString, `say "hi"`)
}

func TestLexer_UnterminatedString(t *testing.T) {
	tokenizeErr(t, `"unterminated`)
}

// --- Boolean and null ---

func TestLexer_BoolTrue(t *testing.T) {
	toks := tokenize(t, "true")
	assertToken(t, toks[0], TokBool, "true")
}

func TestLexer_BoolFalse(t *testing.T) {
	toks := tokenize(t, "false")
	assertToken(t, toks[0], TokBool, "false")
}

func TestLexer_Null(t *testing.T) {
	toks := tokenize(t, "null")
	assertToken(t, toks[0], TokNull, "null")
}

func TestLexer_Nil(t *testing.T) {
	toks := tokenize(t, "nil")
	assertToken(t, toks[0], TokNull, "nil")
}

// --- Identifiers ---

func TestLexer_Identifier(t *testing.T) {
	toks := tokenize(t, "myVar")
	assertToken(t, toks[0], TokIdent, "myVar")
}

func TestLexer_IdentifierWithUnderscore(t *testing.T) {
	toks := tokenize(t, "_my_var_2")
	assertToken(t, toks[0], TokIdent, "_my_var_2")
}

// --- Variable references ---

func TestLexer_VarRef(t *testing.T) {
	toks := tokenize(t, "{myVar}")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokVar, "myVar")
}

func TestLexer_VarRefWithSpaces(t *testing.T) {
	toks := tokenize(t, "{ my_var }")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	assertToken(t, toks[0], TokVar, "my_var")
}

func TestLexer_BraceNotVar(t *testing.T) {
	// { followed by non-letter → map literal token (LBrace)
	toks := tokenize(t, "{'key': 1}")
	if toks[0].Type != TokLBrace {
		t.Errorf("expected LBrace, got %v", toks[0])
	}
}

// --- Operators ---

func TestLexer_ArithmeticOps(t *testing.T) {
	cases := []struct {
		input string
		typ   TokenType
		val   string
	}{
		{"+", TokPlus, "+"},
		{"-", TokMinus, "-"},
		{"*", TokStar, "*"},
		{"/", TokSlash, "/"},
		{"%", TokPercent, "%"},
		{"**", TokPower, "**"},
	}
	for _, c := range cases {
		toks := tokenize(t, c.input)
		assertToken(t, toks[0], c.typ, c.val)
	}
}

func TestLexer_ComparisonOps(t *testing.T) {
	cases := []struct {
		input string
		typ   TokenType
	}{
		{"==", TokEq},
		{"!=", TokNotEq},
		{"<", TokLt},
		{"<=", TokLtEq},
		{">", TokGt},
		{">=", TokGtEq},
	}
	for _, c := range cases {
		toks := tokenize(t, c.input)
		if toks[0].Type != c.typ {
			t.Errorf("input %q: got %v, want %v", c.input, toks[0].Type, c.typ)
		}
	}
}

func TestLexer_LogicalOps(t *testing.T) {
	cases := []struct {
		input string
		typ   TokenType
	}{
		{"&&", TokAnd},
		{"||", TokOr},
		{"!", TokNot},
	}
	for _, c := range cases {
		toks := tokenize(t, c.input)
		if toks[0].Type != c.typ {
			t.Errorf("input %q: got %v, want %v", c.input, toks[0].Type, c.typ)
		}
	}
}

func TestLexer_Punctuation(t *testing.T) {
	input := "()[]{},.;?:."
	toks := tokenize(t, input)
	expected := []TokenType{
		TokLParen, TokRParen,
		TokLBracket, TokRBracket,
		TokLBrace, TokRBrace,
		TokComma, TokDot, TokSemicolon,
		TokQuestion, TokColon, TokDot,
	}
	if len(toks) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(toks), toks)
	}
	for i, e := range expected {
		if toks[i].Type != e {
			t.Errorf("token[%d] = %v, want %v", i, toks[i].Type, e)
		}
	}
}

// --- Comments ---

func TestLexer_LineComment(t *testing.T) {
	toks := tokenize(t, "42 // this is a comment")
	if len(toks) != 1 {
		t.Fatalf("expected 1 token after comment skip, got %d", len(toks))
	}
	assertToken(t, toks[0], TokInt, "42")
}

func TestLexer_CommentOnly(t *testing.T) {
	toks := tokenize(t, "// just a comment")
	if len(toks) != 0 {
		t.Errorf("expected 0 tokens for comment-only input, got %d", len(toks))
	}
}

// --- Whitespace ---

func TestLexer_IgnoresWhitespace(t *testing.T) {
	toks := tokenize(t, "  1   +   2  ")
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(toks))
	}
	assertToken(t, toks[0], TokInt, "1")
	assertToken(t, toks[1], TokPlus, "+")
	assertToken(t, toks[2], TokInt, "2")
}

// --- EOF token ---

func TestLexer_EOFToken(t *testing.T) {
	l := New("")
	toks, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 1 || toks[0].Type != TokEOF {
		t.Errorf("empty input should produce single EOF token, got %v", toks)
	}
}

// --- Complex expression ---

func TestLexer_ComplexExpression(t *testing.T) {
	toks := tokenize(t, "upper({name}) + '!'")
	expected := []struct {
		typ TokenType
		val string
	}{
		{TokIdent, "upper"},
		{TokLParen, "("},
		{TokVar, "name"},
		{TokRParen, ")"},
		{TokPlus, "+"},
		{TokString, "!"},
	}
	if len(toks) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(toks), toks)
	}
	for i, e := range expected {
		assertToken(t, toks[i], e.typ, e.val)
	}
}

// --- Negative int after operator ---

func TestLexer_NegativeAfterOperator(t *testing.T) {
	// 5 + -3: the -3 should be a negative int, not minus + int
	toks := tokenize(t, "5 + -3")
	// Could be [5] [+] [-3] or [5] [+] [-] [3] depending on context
	// Per the lexer: '-' followed by digit after operator → negative number
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d: %v", len(toks), toks)
	}
	assertToken(t, toks[0], TokInt, "5")
	assertToken(t, toks[1], TokPlus, "+")
	assertToken(t, toks[2], TokInt, "-3")
}

// --- Unexpected character ---

func TestLexer_UnexpectedCharacter(t *testing.T) {
	tokenizeErr(t, "@invalid")
}

// --- Line/Col tracking ---

func TestLexer_LineColTracking(t *testing.T) {
	l := New("1\n2")
	toks, _ := l.Tokenize()
	// toks[0] = 1 at line 1, toks[1] = 2 at line 2
	if toks[0].Line != 1 {
		t.Errorf("token 0 line = %d, want 1", toks[0].Line)
	}
	if toks[1].Line != 2 {
		t.Errorf("token 1 line = %d, want 2", toks[1].Line)
	}
}

// --- TokenType.String() ---

func TestTokenType_String(t *testing.T) {
	if TokPlus.String() != "+" {
		t.Errorf("TokPlus.String() = %q, want '+'", TokPlus.String())
	}
	if TokEOF.String() != "EOF" {
		t.Errorf("TokEOF.String() = %q, want 'EOF'", TokEOF.String())
	}
	if TokVar.String() != "VAR" {
		t.Errorf("TokVar.String() = %q, want 'VAR'", TokVar.String())
	}
}

// --- Token.String() ---

func TestToken_String(t *testing.T) {
	tok := Token{Type: TokInt, Value: "42", Line: 1, Col: 1}
	s := tok.String()
	if s == "" {
		t.Error("Token.String() should not be empty")
	}
}
