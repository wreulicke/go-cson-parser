package main

import (
	"fmt"
)

type Parser struct {
	l         *Lexer
	errors    []error
	curToken  Token
	peekToken Token
	indent    int
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() Value {
	return p.parseValue(0)
}

func (p *Parser) parseValue(indent int) Value {
	switch p.curToken.Type {
	case IDENTIFIER:
		return p.parseObject(indent)
	case STRING:
		return p.parseString()
	case NUMBER:
		return p.parseNumber()
	case INDENT:
		p.nextToken()
		p.indent++
		v := p.parseValue(indent)
		return v
	default:
		panic("unexpected token. token:" + TokenName(p.curToken.Type))
	}
}

func (p *Parser) parseObject(indent int) Value {
	object := ObjectValue{}
	identifier := p.parseIdentifier()
	for {
		if !p.expectPeek(COLON) {
			return nil
		}
		p.nextToken()
		v := p.parseValue(p.indent + 1)
		if v == nil {
			return nil
		}
		pair := Pair{
			Key:   &Key{Identifier: identifier},
			Value: v,
		}
		object.Pair = append(object.Pair, pair)
		for p.peekTokenIs(DEDENT) {
			p.nextToken()
			p.indent--
		}
		for p.peekTokenIs(INDENT) {
			p.nextToken()
			p.indent++
		}
		if p.indent != indent {
			break
		}
		if p.peekTokenIs(EOF) {
			break
		}
		if !p.expectPeek(IDENTIFIER) {
			return nil
		}
		identifier = p.parseIdentifier()
	}
	return &object
}

func (p *Parser) parseIdentifier() *Identifier {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseNumber() Value {
	return &NumberValue{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseString() Value {
	return &StringValue{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t TokenType) {
	p.errors = append(p.errors, fmt.Errorf("expected next token to be %s, got %s instead", TokenName(t), TokenName(p.peekToken.Type)))
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	p.errors = append(p.errors, fmt.Errorf("no prefix parse function for %s found", TokenName(t)))
}
