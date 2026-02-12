package lexer

import (
	"monk/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `=+,;(){}`

	// table-driven testing
	// input -> output
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// in same order as above input !
		{token.ASSIGN, "="}, // token.assign here is the constant we defined in package token
		{token.PLUS, "+"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input) // lexer is like a pointer to a token object
	// we initialized that pointer here and read first character in input
	// and stored it in l.ch now we need token so we will ne tokenizing that l.ch

	for i, tt := range tests {
		tok := l.NextToken() // make a token from l.ch
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
