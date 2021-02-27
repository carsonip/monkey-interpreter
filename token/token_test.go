package token

import "testing"

func TestLexer_NextToken(t *testing.T) {
	var tok Token
	l := NewLexer(",. foo let 0123 == != = !foo")
	if l.NextToken().tokenType != TOKEN_COMMA {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_DOT {
		t.Error()
	}
	tok = l.NextToken()
	if tok.tokenType != TOKEN_IDENTIFIER {
		t.Error()
	}
	if tok.literal != "foo" {
		t.Error()
	}
	tok = l.NextToken()
	if tok.tokenType != TOKEN_LET {
		t.Error()
	}
	if tok.literal != "let" {
		t.Error()
	}
	tok = l.NextToken()
	if tok.tokenType != TOKEN_NUMBER {
		t.Error()
	}
	if tok.literal != "0123" {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_EQUAL {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_NOTEQUAL {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_ASSIGNMENT {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_NOT {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_IDENTIFIER {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_EOF {
		t.Error()
	}
}
