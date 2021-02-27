package token

import "testing"

func TestLexer_NextToken(t *testing.T) {
	l := NewLexer(",. let .")
	if l.NextToken().tokenType != TOKEN_COMMA {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_DOT {
		t.Error()
	}
	tok := l.NextToken()
	if tok.tokenType != TOKEN_IDENTIFIER {
		t.Error()
	}
	if tok.literal != "let" {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_DOT {
		t.Error()
	}
	if l.NextToken().tokenType != TOKEN_EOF {
		t.Error()
	}
}
