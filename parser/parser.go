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
	case token.TOKEN_EOF:
		node = nil
	case token.TOKEN_LET:
		node = p.parseLetStatement()
	case token.TOKEN_RETURN:
		node = p.parseReturnStatement()
	case token.TOKEN_IF:
		node = p.parseIfStatement()
	default:
		node = p.parseExpression()
	}
	for p.curTokenIs(token.TOKEN_SEMICOLON) {
		p.next()
	}
	return node
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	l := &ast.LetStatement{
		Token: p.curToken,
	}
	p.expectAndNext(token.TOKEN_LET)
	l.Name = p.parseIdentifier()
	p.expectAndNext(token.TOKEN_ASSIGNMENT)
	l.Value = p.parseExpression()
	return l
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	s := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.expectAndNext(token.TOKEN_RETURN)
	expr := p.parseExpression()
	s.Value = expr
	return s
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	s := &ast.IfStatement{
		Token: p.curToken,
	}
	p.expectAndNext(token.TOKEN_IF)
	expr := p.parseExpression()
	s.Condition = expr
	p.expectAndNext(token.TOKEN_LBRACE)
	for !p.curTokenIs(token.TOKEN_RBRACE) {
		s.Then = append(s.Then, p.NextNode())
	}
	p.expectAndNext(token.TOKEN_RBRACE)
	if p.curTokenIs(token.TOKEN_ELSE) {
		p.expectAndNext(token.TOKEN_ELSE)
		p.expectAndNext(token.TOKEN_LBRACE)
		for !p.curTokenIs(token.TOKEN_RBRACE) {
			s.Else = append(s.Else, p.NextNode())
		}
		p.expectAndNext(token.TOKEN_RBRACE)
	}
	return s
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseExpressionWithPrecedence(0)
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.expectAndNext(token.TOKEN_LPAREN)
	expr := p.parseExpressionWithPrecedence(0)
	p.expectAndNext(token.TOKEN_RPAREN)
	return expr
}

func (p *Parser) parseExpressionWithPrecedence(curPrecedence Precedence) ast.Expression {
	var expr ast.Expression
	switch p.curToken.Type {
	case token.TOKEN_NUMBER:
		expr = p.parseNumber()
	case token.TOKEN_IDENTIFIER:
		expr = p.parseIdentifier()
	case token.TOKEN_PLUS, token.TOKEN_MINUS, token.TOKEN_NOT:
		expr = p.parsePrefixExpression()
	case token.TOKEN_FUNCTION:
		expr = p.parseFunction()
	case token.TOKEN_TRUE, token.TOKEN_FALSE:
		expr = p.parseBoolean()
	case token.TOKEN_LPAREN:
		expr = p.parseGroupedExpression()
	default:
		log.Panicf("expected expression, got %d %s instead", p.curToken.Type, p.curToken.Literal)
		p.next()
		return nil
	}

	for {
		if p.curTokenIs(token.TOKEN_RPAREN, token.TOKEN_COMMA, token.TOKEN_SEMICOLON, token.TOKEN_EOF) {
			return expr
		}

		if p.curTokenIs(token.TOKEN_LPAREN) {
			return p.parseFunctionCall(expr)
		} else if p.curTokenIs(token.TOKEN_LBRACE, token.TOKEN_PLUS, token.TOKEN_MINUS, token.TOKEN_ASTERISK, token.TOKEN_SLASH) {
			precedence := operatorToPrecedence[p.curToken.Type]
			if precedence <= curPrecedence {
				return expr
			}

			expr = p.parseInfixExpression(expr, precedence)
		} else {
			log.Panicf("expected operator, got %d %s instead", p.curToken.Type, p.curToken.Literal)
			return nil
		}
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
	ident := &ast.Identifier{Token: p.curToken}
	p.expectAndNext(token.TOKEN_IDENTIFIER)
	return ident
}

func(p *Parser) parseFunction() *ast.Function {
	fn := &ast.Function{Token: p.curToken}
	p.expectAndNext(token.TOKEN_FUNCTION)
	p.expectAndNext(token.TOKEN_LPAREN)
	isFirst := true
	for !p.curTokenIs(token.TOKEN_RPAREN) {
		if isFirst {
			isFirst = false
		} else {
			p.expectAndNext(token.TOKEN_COMMA)
		}
		ident := p.parseIdentifier()
		fn.Params = append(fn.Params, ident)
	}
	p.expectAndNext(token.TOKEN_RPAREN)
	p.expectAndNext(token.TOKEN_LBRACE)
	for !p.curTokenIs(token.TOKEN_RBRACE) {
		node := p.NextNode()
		fn.Body = append(fn.Body, node)
	}
	p.expectAndNext(token.TOKEN_RBRACE)
	return fn
}

func (p *Parser) expectAndNext(token token.TokenType) {
	if !p.curTokenIs(token) {
		log.Panicf("expected %d, got %d %s instead", token, p.curToken.Type, p.curToken.Literal)
	}
	p.next()
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
	PRECEDENCE_PREFIX
	PRECEDENCE_CALL
)
var operatorToPrecedence = map[token.TokenType]Precedence{
	token.TOKEN_PLUS: PRECEDENCE_PLUS_MINUS,
	token.TOKEN_MINUS: PRECEDENCE_PLUS_MINUS,
	token.TOKEN_ASTERISK: PRECEDENCE_MULTIPLY_DIVIDE,
	token.TOKEN_SLASH: PRECEDENCE_MULTIPLY_DIVIDE,
	token.TOKEN_LPAREN: PRECEDENCE_CALL,
}

func (p *Parser) parseInfixExpression(left ast.Expression, curPrecedence Precedence) ast.Expression {
	expr := &ast.InfixExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.next()
	expr.Right = p.parseExpressionWithPrecedence(curPrecedence)
	return expr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token: p.curToken,
	}
	p.next()
	expr.Right = p.parseExpressionWithPrecedence(PRECEDENCE_PREFIX)
	return expr
}

func (p *Parser) parseBoolean() ast.Expression {
	val := false
	if p.curToken.Type == token.TOKEN_TRUE {
		val = true
	}
	expr := &ast.Boolean{
		Token: p.curToken,
		Value: val,
	}
	p.next()
	return expr
}

func (p *Parser) parseFunctionCall(expr ast.Expression) *ast.FunctionCall {
	fnCall := &ast.FunctionCall{Token: p.curToken, FunctionExpr: expr}
	p.expectAndNext(token.TOKEN_LPAREN)
	first := true
	for !p.curTokenIs(token.TOKEN_RPAREN) {
		if first {
			first = false
		} else {
			p.expectAndNext(token.TOKEN_COMMA)
		}
		expr := p.parseExpression()
		fnCall.Arguments = append(fnCall.Arguments, expr)
	}
	p.expectAndNext(token.TOKEN_RPAREN)
	return fnCall
}
