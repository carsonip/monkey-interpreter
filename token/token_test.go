package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	var tok Token
	l := NewLexer(",. foo let 0123 == != = !foo")
	assert.Equal(t, TOKEN_COMMA, l.NextToken().Type)
	assert.Equal(t, TOKEN_DOT, l.NextToken().Type)
	tok = l.NextToken()
	assert.Equal(t, TOKEN_IDENTIFIER, tok.Type)
	assert.Equal(t, "foo", tok.Literal)
	tok = l.NextToken()
	assert.Equal(t, TOKEN_LET, tok.Type)
	assert.Equal(t, "let", tok.Literal)
	tok = l.NextToken()
	assert.Equal(t, TOKEN_NUMBER, tok.Type)
	assert.Equal(t, "0123", tok.Literal)
	assert.Equal(t, TOKEN_EQUAL, l.NextToken().Type)
	assert.Equal(t, TOKEN_NOTEQUAL, l.NextToken().Type)
	assert.Equal(t, TOKEN_ASSIGNMENT, l.NextToken().Type)
	assert.Equal(t, TOKEN_NOT, l.NextToken().Type)
	assert.Equal(t, TOKEN_IDENTIFIER, l.NextToken().Type)
	assert.Equal(t, TOKEN_EOF, l.NextToken().Type)
}
