package token

import "testing"

func TestLexer_NextToken(t *testing.T) {
	var tok Token
	l := NewLexer(",. foo let 0123 == != = !foo")
	if l.NextToken().Type != TOKEN_COMMA {
		t.Error()
	}
	if l.NextToken().Type != TOKEN_DOT {
		t.Error()
	}
	tok = l.NextToken()
	if tok.Type != TOKEN_IDENTIFIER {
		t.Error()
	}
	if tok.Literal != "foo" {
		t.Error()
	}
	tok = l.NextToken()
	if tok.Type != TOKEN_LET {
		t.Error()
	}
	if tok.Literal != "let" {
		t.Error()
	}
	tok = l.NextToken()
	if tok.Type != TOKEN_NUMBER {
		t.Error()
	}
	if tok.Literal != "0123" {
		t.Error()
	}
	if l.NextToken().Type != TOKEN_EQUAL {
		t.Error()
	}
	if l.NextToken().Type != TOKEN_NOTEQUAL {
		t.Error()
	}
	if l.NextToken().Type != TOKEN_ASSIGNMENT {
		t.Error()
	}
	if l.NextToken().Type != TOKEN_NOT {
		t.Error()
	}
	if l.NextToken().Type != TOKEN_IDENTIFIER {
		t.Error()
	}
	if l.NextToken().Type != TOKEN_EOF {
		t.Error()
	}
}
