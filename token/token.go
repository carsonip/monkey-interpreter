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
	TOKEN_IDENTIFIER
)

var charToToken = map[byte]TokenType{
	0: TOKEN_EOF,
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
}

type Token struct {
	tokenType TokenType
	literal string
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

func (l *Lexer) readIdentifier() {
	for l.ch != ' ' {
		l.readChar()
	}
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	l.eatWhitespace()
	if isAlpha(l.ch) {
		lastPos := l.pos
		l.readIdentifier()
		return Token{
			tokenType: TOKEN_IDENTIFIER,
			literal: l.input[lastPos:l.pos],
		}
	} else {
		tokenType, ok := charToToken[l.ch]
		l.readChar()
		if !ok {
			panic("unknown token")
		}
		return Token{
			tokenType: tokenType,
			literal: string(l.ch),
		}
	}
}

func isAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}
