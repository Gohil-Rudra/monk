package lexer

import (
	"monk/token"
)

// #####################################Lexing is done here and passed above########
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input} // returns address
	l.readChar()              // now our l.ch is ready else it would be still empty !
	return l
}

func (l *Lexer) readChar() { // func (L *lexer) means its a function of class(struct) lexer
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ascii for NUL in Go printing s[0] for s = "ABC" produces 65
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition // position - is position of character present in l.ch
	l.readPosition += 1         // readPosition is position of char next to l.ch
}

func (l *Lexer) NextToken() token.Token { // token is our previous program
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		// tok = newToken(token.EOF, "") as "" is not a byt we will have to do...
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		tok = newToken(token.ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

// helper function of nextToken
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
