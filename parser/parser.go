package parser

import (
	"github.com/carsonip/monkey-interpreter/ast"
	"github.com/carsonip/monkey-interpreter/token"
	"log"
	"strconv"
)

type Parser struct {
	lexer *token.Lexer
	curToken token.Token
}

func NewParser(l *token.Lexer) Parser {
	p := Parser{lexer: l}
	p.next()
	return p
}

func (p *Parser) next() {
	p.curToken = p.lexer.NextToken()
}

func (p *Parser) NextNode() ast.Node {
	var node ast.Node
	switch p.curToken.Type {
	case token.TOKEN_LET:
		node = p.parseLetStatement()
	default:
		node = p.parseExpression()
	}
	return node
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	l := &ast.LetStatement{
		Token: p.curToken,
	}
	p.next()
	l.Name = p.parseIdentifier()
	if !p.curTokenIs(token.TOKEN_ASSIGNMENT) {
		log.Panicf("expected =, got %d instead", p.curToken.Type)
	}
	p.next()
	l.Value = p.parseExpression()
	return l
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseExpressionWithPrecedence(0)
}

func (p *Parser) parseExpressionWithPrecedence(curPrecedence Precedence) ast.Expression {
	var exp ast.Expression

	switch p.curToken.Type {
	case token.TOKEN_NUMBER:
		exp = p.parseNumber()
	case token.TOKEN_IDENTIFIER:
		exp = p.parseIdentifier()
	case token.TOKEN_FUNCTION:
	case token.TOKEN_TRUE:
	case token.TOKEN_FALSE:
		log.Panicf("not implemented")
		p.next()
		return nil
	default:
		log.Panicf("expected expression, got %d %s instead", p.curToken.Type, p.curToken.Literal)
		p.next()
		return nil
	}

	for {
		if p.curTokenIs(token.TOKEN_SEMICOLON) || p.curTokenIs(token.TOKEN_EOF) {
			return exp
		}

		if !p.curTokenIs(token.TOKEN_PLUS, token.TOKEN_MINUS, token.TOKEN_ASTERISK, token.TOKEN_SLASH) {
			log.Panicf("expected operator, got %d %s instead", p.curToken.Type, p.curToken.Literal)
			return nil
		}

		precedence := operatorToPrecedence[p.curToken.Type]
		if precedence <= curPrecedence {
			return exp
		}

		exp = p.parseInfixExpression(exp, precedence)
	}
}

func (p *Parser) parseNumber() *ast.NumberLiteral {
	if !p.curTokenIs(token.TOKEN_NUMBER) {
		log.Panicf("expected number, got %d %s instead", p.curToken.Type, p.curToken.Literal)
	}
	if val, err := strconv.Atoi(p.curToken.Literal); err != nil {
		log.Panicf("bad number %s", p.curToken.Literal)
		p.next()
		return nil
	} else {
		lit := &ast.NumberLiteral{
			Token: p.curToken,
			Value: val,
		}
		p.next()
		return lit
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	if !p.curTokenIs(token.TOKEN_IDENTIFIER) {
		log.Panicf("expected identifier, got %d %s instead", p.curToken.Type, p.curToken.Literal)
	}
	lit := &ast.Identifier{Token: p.curToken}
	p.next()
	return lit
}

func (p *Parser) curTokenIs(tokenTypes ...token.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.curToken.Type == tokenType {
			return true
		}
	}
	return false
}

type Precedence int
const (
	_ Precedence = iota
	PRECEDENCE_PLUS_MINUS
	PRECEDENCE_MULTIPLY_DIVIDE
)
var operatorToPrecedence = map[token.TokenType]Precedence{
	token.TOKEN_PLUS: PRECEDENCE_PLUS_MINUS,
	token.TOKEN_MINUS: PRECEDENCE_PLUS_MINUS,
	token.TOKEN_ASTERISK: PRECEDENCE_MULTIPLY_DIVIDE,
	token.TOKEN_SLASH: PRECEDENCE_MULTIPLY_DIVIDE,
}

func (p *Parser) parseInfixExpression(left ast.Expression, curPrecedence Precedence) ast.Expression {
	exp := &ast.InfixExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.next()
	exp.Right = p.parseExpressionWithPrecedence(curPrecedence)
	return exp
}
