package parser

import (
	"fmt"
	"monk/ast"
	"monk/lexer_2"
	"testing"
)

// Let Statement Testing
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5", "x", 5}, // ; is optional and that saved me from a big bug :)
		{"let y = true", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer_2.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		CheckParseErrors(p, t)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got = %d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			// we will print why it failed in testLetStatement
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

// Return Statement Testing
func TestReturnStatements(t *testing.T) {
	input := `

	return 5;
	return 9347982;

	`

	l := lexer_2.New(input)
	p := New(l)

	program := p.ParseProgram()
	CheckParseErrors(p, t)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 2 { // also included ";"
		t.Fatalf("program.Statements does not contain 2 statements.got = %d", len(program.Statements))
	}
	// well we dont have much to check except that it is a ReturnStatement
	for _, stmt := range program.Statements {
		if !testReturnStatement(t, stmt) {
			continue
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer_2.New(input)

	p := New(l)

	program := p.ParseProgram()

	CheckParseErrors(p, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements . got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("Expression is not an*ast.Identifier . got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident value not %s. got %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident Token.Literal not %s. got %s", "foobar", ident.TokenLiteral())
	}

}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"
	l := lexer_2.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParseErrors(p, t)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements . got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement. got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Expression is not an *ast.IntegerLiteral . got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("ident value not %d. got %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("ident Token.Literal not %s. got %s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{} // its will now be a pair {32,int64} , {"abc",strign} , {true,bool} so we can type check in Literal expression !
	}{
		{input: "!5;", operator: "!", value: 5},
		{input: "-15;", operator: "-", value: 15},
		{input: "!true", operator: "!", value: true},
		{input: "!false", operator: "!", value: false},
	}

	for _, tt := range prefixTests {
		l := lexer_2.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParseErrors(p, t)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements don't contain %d statements. got %d", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression Statement. got %T .", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got = %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got %s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}

}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"true == true", true, "==", true},
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	for _, tt := range infixTests {
		l := lexer_2.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParseErrors(p, t)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements don't contain %d statements. got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{input: "-1*2+3", output: "(((-1)*2)+3)"},
		{input: "!-a", output: "(!(-a))"},
		{input: "a+b+c", output: "((a+b)+c)"},
		{input: "a+b-c", output: "((a+b)-c)"},
		{input: "a*b*c", output: "((a*b)*c)"},
		{input: "a*b/c", output: "((a*b)/c)"},
		{input: "a + b / c", output: "(a+(b/c))"},
		{input: "a + b * c + d / e - f", output: "(((a+(b*c))+(d/e))-f)"},
		{input: "true", output: "true"},
		{input: "false", output: "false"},
		{input: "4<5 == true", output: "((4<5)==true)"},
		{input: "5<4 == false", output: "((5<4)==false)"},

		// grouped expressions !
		{input: "5+(1+2)+3", output: "((5+(1+2))+3)"},
		{input: "!(true==true)", output: "(!(true==true))"},

		{input: "3 + add(4,5) * 2", output: "(3+(add(4,5)*2))"},
		{
			"add(a,b,1,2*3,4+5,add(6,7*8))",
			"add(a,b,1,(2*3),(4+5),add(6,(7*8)))",
		},
	}

	for _, tt := range tests {
		l := lexer_2.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParseErrors(p, t)

		if program.String() != tt.output {
			t.Errorf("expected %q got %q...", tt.output, program.String())
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer_2.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParseErrors(p, t)
		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements . got %d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement . got %T", program.Statements[0])
		}
		if !testBoolean(t, stmt.Expression, tt.output) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if ( x < y ) { x }"
	l := lexer_2.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParseErrors(p, t)
	if len(program.Statements) != 1 {
		t.Fatalf(" program.Statements does not contain %d Statetments got %d \n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. got : %T", program.Statements[0])
	}

	ifexp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not an ast.IfExpression. got %T", stmt.Expression)
	}

	if !testInfixExpression(t, ifexp.Condition, "x", "<", "y") {
		return
	}

	if len(ifexp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement got , %d", len(ifexp.Consequence.Statements))
	}

	consequence, ok := ifexp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not an ast.ExpressionStatement , got %T", ifexp.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, consequence.Expression, "x") {
		return
	}

	if ifexp.Alternative != nil {
		t.Errorf("ifexp.Alternative.Statements is not nil got %+v", ifexp.Alternative)
	}
	//%+v- include struct field names and %v is print it in human readable {} form.
}

func TestIfElseExpression(t *testing.T) {
	input := "if ( x < y ) { x } else { y }"
	l := lexer_2.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParseErrors(p, t)
	if len(program.Statements) != 1 {
		t.Fatalf(" program.Statements does not contain %d Statetments got %d \n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. got : %T", program.Statements[0])
	}

	ifexp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not an ast.IfExpression. got %T", stmt.Expression)
	}

	if !testInfixExpression(t, ifexp.Condition, "x", "<", "y") {
		return
	}

	if len(ifexp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement got , %d", len(ifexp.Consequence.Statements))
	}

	consequence, ok := ifexp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not an ast.ExpressionStatement , got %T", ifexp.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, consequence.Expression, "x") {
		return
	}

	if len(ifexp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statement got , %d", len(ifexp.Consequence.Statements))
	}

	alternative, ok := ifexp.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not an ast.ExpressionStatement , got %T", ifexp.Alternative.Statements[0])
	}

	if !testLiteralExpression(t, alternative.Expression, "y") {
		return
	}

}

func TestFunctionExpression(t *testing.T) {
	input := "fn(x,y) {x + y;}"

	l := lexer_2.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("program.body contains %d statements , got : %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got : %T", program.Statements[0])
	}

	fnExp, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral , got : %T", stmt.Expression)
	}

	if len(fnExp.Parameter) != 2 {
		t.Fatalf("function literal parameter wrong. want 2 got %d", len(fnExp.Parameter))
	}

	testLiteralExpression(t, fnExp.Parameter[0], "x")
	testLiteralExpression(t, fnExp.Parameter[1], "y")

	if len(fnExp.Body.Statements) != 1 {
		t.Fatalf("function.Body.statements has not 1 statements , got : %d", len(fnExp.Body.Statements))
	}

	bodyStmt, ok := fnExp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not an Expression Statement , got :%T", fnExp.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameter(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn(){};", expectedParams: []string{}},
		{input: "fn(x){};", expectedParams: []string{"x"}},
		{input: "fn(x,y){};", expectedParams: []string{"x", "y"}},
	}

	for _, tt := range tests {
		l := lexer_2.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParseErrors(p, t)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameter) != len(tt.expectedParams) {
			t.Fatalf("parameters length is wrong , got  : %d expected :%d", len(function.Parameter), len(tt.expectedParams))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameter[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {

	input := "add(1,2*3,4+5);"

	l := lexer_2.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParseErrors(p, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements , got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement got : %T", program.Statements[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not of type CallExpression instead got %T", stmt.Expression)
	}

	if !testIdentifier(t, callExp.Function, "add") {
		return
	}

	if len(callExp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got :%d", len(callExp.Arguments))
	}

	testLiteralExpression(t, callExp.Arguments[0], 1)
	testInfixExpression(t, callExp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExp.Arguments[2], 4, "+", 5)
}

// -------------------Helper functions interrelated---------------------

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let' got %s", s.TokenLiteral()) // errorf doesnt stop the execution
		return false
	}

	letstmt, ok := s.(*ast.LetStatement) // checks does this interface points to Letstatement struct ?

	if !ok {
		t.Errorf("s not *ast.LetStatement , got = %T", s)
		return false
	} // i still have doubt why we did this...

	if letstmt.Name.Value != name {
		t.Errorf("letstmt.Name.Value not '%s'. got=%s", name, letstmt.Name.Value)
		return false
	}

	if letstmt.Name.TokenLiteral() != name {
		t.Errorf("letstmt.Name.Token.Literal not '%s'. got=%s", name, letstmt.Name.TokenLiteral())
		return false
	}
	return true
}

func testReturnStatement(t *testing.T, s ast.Statement) bool {
	if s.TokenLiteral() != "return" {
		t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
			s.TokenLiteral())
		return false
	}
	// value, ok := interfaceValue.(ConcreteType)

	_, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("stmt is not *ast.ReturnStatement instead got %T", s)
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral . got %T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("il value not %d . got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d .got %s", value, integer.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, il ast.Expression, value string) bool {
	ident, ok := il.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.identifier. got %T", il)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s got %s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s got %s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testBoolean(t *testing.T, bl ast.Expression, value bool) bool {
	bo, ok := bl.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not ast.Boolean . got %T", bl)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.value not %t. got %t", bo.Value, value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t , got %s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}
	t.Errorf("type of exp not handled. got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	Opexp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got = %T(%s)", exp, exp) //here %s will take the String() method defined for exp to spit the operator...
		return false
	}

	if !testLiteralExpression(t, Opexp.Left, left) {
		return false
	}

	if Opexp.Operator != operator {
		t.Errorf("exp.operator is not %s , got %s ", operator, Opexp.Operator)
		return false
	}

	if !testLiteralExpression(t, Opexp.Right, right) {
		return false
	}
	return true
}

// ------------------------Error Handling-----------------------

func CheckParseErrors(p *Parser, t *testing.T) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf(" Parse Error : %q ", msg)
	}
	t.FailNow()
}
