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
	switch p.curToken.Type {
	case token.TOKEN_NUMBER:
		return p.parseNumber()
	case token.TOKEN_IDENTIFIER:
		return p.parseIdentifier()
	case token.TOKEN_FUNCTION:
	case token.TOKEN_TRUE:
	case token.TOKEN_FALSE:
	default:
		log.Panicf("expected expression, got %d %s instead", p.curToken.Type, p.curToken.Literal)
		p.next()
	}
	return nil
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

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}
