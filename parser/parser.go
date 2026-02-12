package parser // constructing the AST we defined in ast.go
// note in HP fn + prt sc is insert key...
// we return *concretetype and interfacetype without *
import (
	"fmt"
	"monk/ast"
	"monk/lexer_2"
	"monk/token"
	"strconv"
)

//--------------------------------------------CONSTANTS---------------------------
/* => iota is inbuilt auto-incrementing counter
const (
    _           int = 0
    LOWEST          = 1
    EQUALS          = 2
    LESSGREATER     = 3
    SUM             = 4
    PRODUCT         = 5
    PREFIX          = 6 // -X or !X
    CALL            = 7 //myFunction(X)
)
*/
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

// --------------------infix----------------------

var precedences = map[token.TokenType]int{

	token.EQ:      EQUALS,
	token.NOT_EQ:  EQUALS,
	token.LT:      LESSGREATER,
	token.GT:      LESSGREATER,
	token.PLUS:    SUM,
	token.MINUS:   SUM,
	token.SLASH:   PRODUCT,
	token.ASTERIK: PRODUCT,
	token.LPAREN:  CALL,
}

func (p *Parser) peekPrecedence() int {
	prec, ok := precedences[p.peekToken.Type]
	if ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	prec, ok := precedences[p.currToken.Type]
	if ok {
		return prec
	}
	return LOWEST
}

//-----------------------------------------------------------------

type Parser struct {
	l         *lexer_2.Lexer // -> pointer to lexer from which we will call nextToken()
	currToken token.Token    // -> current token
	peekToken token.Token    // -> next token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn // key is token.TokenType and value is the prefixfn

	infixParseFns map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tT token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tT] = fn
}

func (p *Parser) registerInfix(tT token.TokenType, fn infixParseFn) {
	p.infixParseFns[tT] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) { // will called when an error "OCCURS"
	// we have 2 tokens now t(what we expect) and p.peekToken
	msg := fmt.Sprintf("Expected next token to be %s but instead got %s", t, p.peekToken.Type) // Sprintf is String printf which returns a String...instead of printing it.

	p.errors = append(p.errors, msg)
}

func New(l *lexer_2.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	// above we are storing the function and not the value so no ()
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERIK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)

	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)

	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.registerPrefix(token.IF, p.parseIfExpression)

	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(token.LPAREN, p.parseCallExpression)
	return p
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))

	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameter = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{} // this {} means not nil

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken() // consume ( reach first param
	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume ident we defined above
		p.nextToken() // consume comma
		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) { // PAREN means ( dont mistake with {
		return nil
	}

	p.nextToken() // consume (

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	p.nextToken()
	block.Statements = []ast.Statement{}

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement() // well you can see that parseStatement decides what to run based on first token that is say it will run parseLetStatement or maybe returnStatement
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	//defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}
	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression

}

func (p *Parser) parsePrefixExpression() ast.Expression {

	//defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX) // we didnt do currPredence here because its a prefix which means its only lower to func call else its executed first everytime so we just used PREFIX instead

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	//defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	//defer untrace(trace("parseIntegerLiteral"))

	lit := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as Int64", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit

}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	//defer untrace(trace("ParseProgram"))
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	//defer untrace(trace("parseStatement"))
	// if we are returning concrete type than * else not * also in case of slices like []statements we sent * due to slices e
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}
func (p *Parser) parseLetStatement() *ast.LetStatement {

	//defer untrace(trace("parseLetStatement"))
	stmt := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we can parse them

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	//defer untrace(trace("parseReturnStatement"))
	stmt := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken() // to skip "return..."
	// "Now we are in expression territory so !"

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	//defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {

	//defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.currToken.Type] // we are looking in map now to identify how to parse it...polymorphic it can be an IntegerLiteral or Identifier or a operator symnol token !

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil // as we are only parsing prefix for now so just return when done with prefix !
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { //it will always give default as LOWEST so we will chek below if a infixparse function exists or not ?
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp

}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found...", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool { // its the most common thing in all parsers...
	if p.peekTokenIs(t) {
		p.nextToken() // note we advance the currToken and peekToken here...pls note this :( you did a blunder while not considering this in parseIfExpression
		return true
	} else {
		p.peekError(t)
		return false
	}
}

type prefixParseFn /*returns a function...->*/ func() ast.Expression // just a signature...

type infixParseFn /*returns a function...->*/ func(ast.Expression) ast.Expression

// its like giving a name to this anonymous functions

/*
        _=_
      q(-_-)p
      '_) (_`
      /__/  \
    _(<_   / )_
   (__\_\_|_/__)
*/
