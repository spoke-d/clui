package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	LOWEST = iota
	ASSIGNMENT
	CALL
)

var precedence = map[TokenType]int{
	ASSIGN: ASSIGNMENT,
	LPAREN: CALL,
}

type Parser struct {
	lex *Lexer

	errors []string

	currentToken Token
	peekToken    Token

	prefix map[TokenType]PrefixFunc
	infix  map[TokenType]InfixFunc
}

type PrefixFunc func() Expression
type InfixFunc func(Expression) Expression

// NewParser creates a parser for consuming a lexer tokens.
func NewParser(lex *Lexer) *Parser {
	p := &Parser{
		lex: lex,
	}
	p.prefix = map[TokenType]PrefixFunc{
		IDENT:  p.parseIdentifier,
		INT:    p.parseInteger,
		FLOAT:  p.parseFloat,
		STRING: p.parseString,
		LPAREN: p.parseGroup,
		TRUE:   p.parseBool,
		FALSE:  p.parseBool,
	}
	p.infix = map[TokenType]InfixFunc{
		ASSIGN: p.parseInfixExpression,
		LPAREN: p.parseCall,
	}
	p.nextToken()
	p.nextToken()
	return p
}

// Run the parser to the end, which is either an EOF or an error.
func (p *Parser) Run() (*QueryExpression, error) {
	var exp QueryExpression
	for p.currentToken.Type != EOF {
		exp.Expressions = append(exp.Expressions, p.parseExpressionStatement())
		p.nextToken()
	}
	var err error
	if len(p.errors) > 0 {
		err = errors.Errorf(strings.Join(p.errors, "\n"))
		return nil, errors.WithStack(err)
	}
	return &exp, nil
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{
		Token: p.currentToken,
	}
}

func (p *Parser) parseString() Expression {
	return &String{
		Token: p.currentToken,
	}
}

func (p *Parser) parseBool() Expression {
	value, err := strconv.ParseBool(p.currentToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v could not parse %q as bool", p.currentToken.Pos, p.currentToken.Literal)
		p.errors = append(p.errors, msg)
	}
	return &Bool{
		Token: p.currentToken,
		Value: value,
	}
}

func (p *Parser) parseInteger() Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v could not parse %q as integer", p.currentToken.Pos, p.currentToken.Literal)
		p.errors = append(p.errors, msg)
	}
	return &Integer{
		Token: p.currentToken,
		Value: value,
	}
}

func (p *Parser) parseFloat() Expression {
	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v could not parse %q as float", p.currentToken.Pos, p.currentToken.Literal)
		p.errors = append(p.errors, msg)
	}
	return &Float{
		Token: p.currentToken,
		Value: value,
	}
}

func (p *Parser) parseExpressionStatement() Expression {
	stmt := &ExpressionStatement{
		Token: p.currentToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.isPeekToken(SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefix[p.currentToken.Type]
	if prefix == nil {
		if p.currentToken.Type != EOF {
			msg := fmt.Sprintf("Syntax Error:%v invalid character '%s' found", p.currentToken.Pos, p.currentToken.Type)
			p.errors = append(p.errors, msg)
		}
		return nil
	}
	leftExp := prefix()

	// Run the infix function until the next token has
	// a higher precedence.
	for !p.isPeekToken(SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infix[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroup() Expression {
	p.nextToken()
	if p.currentToken.Type == LPAREN && p.isCurrentToken(RPAREN) {
		// This is an empty group, not sure what we should do here.
		return &Empty{
			Token: p.currentToken,
		}
	}

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseCall(left Expression) Expression {
	if p.isPeekToken(RPAREN) {
		currentToken := p.currentToken
		p.nextToken()
		return &CallExpression{
			Token: currentToken,
			Name:  left,
		}
	}

	p.nextToken()

	arguments := []Expression{
		p.parseExpression(LOWEST),
	}
	for p.isPeekToken(COMMA) {
		p.nextToken()
		p.nextToken()
		arguments = append(arguments, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(RPAREN) {
		return nil
	}

	return &CallExpression{
		Token:     p.currentToken,
		Name:      left,
		Arguments: arguments,
	}
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedence[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) isPeekToken(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) isCurrentToken(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.isPeekToken(t) {
		p.nextToken()
		return true
	}
	msg := fmt.Sprintf("Syntax Error: %v expected token to be %s, got %s instead", p.currentToken.Pos, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
	return false
}
