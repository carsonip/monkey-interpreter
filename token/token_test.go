package token

import "testing"

func TestLexer_NextToken(t *testing.T) {
	l := NewLexer(",.")
	if l.NextToken() != TOKEN_COMMA {
		t.Error()
	}
	if l.NextToken() != TOKEN_DOT {
		t.Error()
	}
	if l.NextToken() != TOKEN_EOF {
		t.Error()
	}
}
