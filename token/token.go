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

func (l *Lexer) peekChar() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	} else {
		return l.input[l.pos+1]
	}
}

func (l *Lexer) readIdentifier() string {
	lastPos := l.pos
	for isAlpha(l.ch) {
		l.readChar()
	}
	return l.input[lastPos:l.pos]
}

func (l *Lexer) readNumber() string {
	lastPos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[lastPos:l.pos]
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	l.eatWhitespace()
	if isAlpha(l.ch) {
		str := l.readIdentifier()
		var tokenType TokenType
		if keywordType, ok := keywords[str]; ok {
			tokenType = keywordType
		} else {
			tokenType = TOKEN_IDENTIFIER
		}
		return Token{
			Type:    tokenType,
			Literal: str,
		}
	} else if isDigit(l.ch) {
		str := l.readNumber()
		return Token{
			Type:    TOKEN_NUMBER,
			Literal: str,
		}
	} else if l.ch == '=' {
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return Token{
				Type:    TOKEN_EQUAL,
				Literal: "==",
			}
		} else {
			return Token{
				Type:    TOKEN_ASSIGNMENT,
				Literal: "=",
			}
		}
	} else if l.ch == '!' {
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return Token{
				Type:    TOKEN_NOTEQUAL,
				Literal: "!=",
			}
		} else {
			return Token{
				Type:    TOKEN_NOT,
				Literal: "!",
			}
		}
	} else {
		tokenType, ok := charToToken[l.ch]
		ch := l.ch
		l.readChar()
		if !ok {
			return Token{
				Type:    TOKEN_ILLEGAL,
				Literal: "",
			}
		}
		return Token{
			Type:    tokenType,
			Literal: string(ch),
		}
	}
}

func isAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
