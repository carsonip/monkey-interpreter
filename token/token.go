package token

type TokenType int

const (
	_ = iota
	TOKEN_EOF
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_ASTERISK
	TOKEN_SLASH
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_SQUOTE
	TOKEN_DQOUTE
)

var charToToken = map[byte]TokenType{
	'(': TOKEN_LPAREN,
	')': TOKEN_RPAREN,
	',': TOKEN_COMMA,
	'.': TOKEN_DOT,
	'+': TOKEN_PLUS,
	'-': TOKEN_MINUS,
	'*': TOKEN_ASTERISK,
	'/': TOKEN_SLASH,
	'{': TOKEN_LBRACE,
	'}': TOKEN_RBRACE,
	'\'': TOKEN_SQUOTE,
	'"': TOKEN_DQOUTE,
	0: TOKEN_EOF,
}

type Lexer struct {
	input string
	pos int
	ch byte
}

func NewLexer(input string) Lexer {
	l := Lexer{input: input, pos: -1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	l.pos++
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
	}
}

func (l *Lexer) NextToken() TokenType {
	tokenType, ok := charToToken[l.ch]
	l.readChar()
	if !ok {
		panic("unknown token")
	}
	return tokenType
}
