package evaluator

import (
	"monk/lexer_2"
	"monk/object"
	"monk/parser"
	"testing"
)

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x){x}; identity(5)", 5},
		{"let identity = fn{return x;x*x};identity(5)", 5},
		{"let double = fn{return x*2};double(4)", 8},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctions(t *testing.T) {
	input := "fn(x) {x + 2;};"
	eval := testEval(input)

	fn, ok := eval.(*object.Function)

	if !ok {
		t.Fatalf("object is not function ! Got : %T", eval)
	}

	if len(fn.Parameter) != 1 {
		t.Fatalf("function has wrong parameters...,Got %v", len(fn.Parameter))
	}

	if fn.Parameter[0].String() != "x" {
		t.Fatalf("parameter is not x got : %q", fn.Parameter[0])
	}

	if fn.Body.String() != "(x+2)" {
		t.Fatalf("body is not %q , got : %q", "(x+2)", fn.Body.String())
	}

}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5 ; a", 5},
		{"let a = 5*5; a", 25},
		{"let a = 5; let b = a; b", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)

	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "Type Mismatch : INTEGER + BOOLEAN"},

		{"5 + true ; 5", "Type Mismatch : INTEGER + BOOLEAN"},

		{"-true", "Unknown Operator : -BOOLEAN"},

		{"true + false", "Unknown Operator : BOOLEAN + BOOLEAN"},

		{"if (10>1){ true + false ; }", "Unknown Operator : BOOLEAN + BOOLEAN"},

		{"foobar", "Identifier Not Found : foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errorObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No error Object Returned. got %T %+v", evaluated, evaluated)
			continue
		}
		if errorObj.Message != tt.expectedMessage {
			t.Errorf("Wrong Error Message : expected = %q , got %q", tt.expectedMessage, errorObj.Message)
		}
	}
}

func TestReturnObject(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10; ", 10},
		{"return 10; 9; ", 10},
		{"9;return 2*5; 9 ", 10},
		{
			`if (10 > 1){
				return 10;
			}
			return 1;
			`, 10, // multiline -> needs comma
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1<2) { 10 }", 10},
		{"if (1>2) { 10 }", nil},
		{"if (1>2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}

	}

}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("obj is not null. got %T", obj)
		return false
	} else {
		return true
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerLiteral(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"5+5+5", 15},
		{"4*5-10", 10},
		{"5*5*5", 125},
		{"20 + 2*-10", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.output)
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input string
		ouput bool
	}{
		{"true", true},
		{"false", false},
		{"5>5", false},
		{"5==5", true},
		{"5<5", false},
		{"5>4", true},
		{"4<5", true},
		{"5!=5", false},
		{"true == true", true},
		{"true != false", true},
		{"(1 < 2) == true", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.ouput)
	}
}

// This is the Heart !!
func testEval(input string) object.Object {
	l := lexer_2.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Expected Integer object got %T ", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value...want : %d , got : %d", expected, result.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Expected Boolean Object...got :%T", obj)
		return false
	}
	if expected != result.Value {
		t.Errorf("object has wrong value , got : %t , Expected : %t", result.Value, expected)
		return false
	}
	return true
}
