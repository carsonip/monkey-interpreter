package token

type TokenType int

const (
	_ TokenType = iota
	TOKEN_ILLEGAL
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
	TOKEN_NUMBER
	TOKEN_FUNCTION
	TOKEN_LET
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_RETURN
	TOKEN_EQUAL
	TOKEN_NOTEQUAL
	TOKEN_ASSIGNMENT
	TOKEN_NOT
	TOKEN_LT
	TOKEN_GT
	TOKEN_SEMICOLON
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
	'<': TOKEN_LT,
	'>': TOKEN_GT,
	';': TOKEN_SEMICOLON,
}

var keywords = map[string]TokenType{
	"fn": TOKEN_FUNCTION,
	"let": TOKEN_LET,
	"true": TOKEN_TRUE,
	"false": TOKEN_FALSE,
	"if": TOKEN_IF,
	"else": TOKEN_ELSE,
	"return": TOKEN_RETURN,
}

type Token struct {
	Type    TokenType
	Literal string
}

func newToken(tokenType TokenType, literal string) Token {
	return Token{
		Type: tokenType,
		Literal: literal,
	}
}

type lexer struct {
	input string
	pos int
	ch byte
}

func NewLexer(input string) lexer {
	l := lexer{input: input, pos: -1}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	l.pos++
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
	}
}

func (l *lexer) peekChar() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	} else {
		return l.input[l.pos+1]
	}
}

func (l *lexer) readIdentifier() string {
	lastPos := l.pos
	for isAlpha(l.ch) {
		l.readChar()
	}
	return l.input[lastPos:l.pos]
}

func (l *lexer) readNumber() string {
	lastPos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[lastPos:l.pos]
}

func (l *lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *lexer) NextToken() Token {
	l.eatWhitespace()
	if isAlpha(l.ch) {
		str := l.readIdentifier()
		var tokenType TokenType
		if keywordType, ok := keywords[str]; ok {
			tokenType = keywordType
		} else {
			tokenType = TOKEN_IDENTIFIER
		}
		return newToken(tokenType, str)
	} else if isDigit(l.ch) {
		str := l.readNumber()
		return newToken(TOKEN_NUMBER, str)
	} else if l.ch == '=' {
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return newToken(TOKEN_EQUAL, "==")
		} else {
			return newToken(TOKEN_ASSIGNMENT, "=")
		}
	} else if l.ch == '!' {
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return newToken(TOKEN_NOTEQUAL, "!=")
		} else {
			return newToken(TOKEN_NOT, "!")
		}
	} else {
		tokenType, ok := charToToken[l.ch]
		ch := l.ch
		l.readChar()
		if !ok {
			return newToken(TOKEN_ILLEGAL, "")
		}
		return newToken(tokenType, string(ch))
	}
}

func isAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
