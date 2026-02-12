package lexer_2

import (
	"monk/token"
)

type Lexer struct {
	input        string
	ch           byte
	position     int
	readPosition int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar() // this is for 1st '=' as we want advance to 2nd '=' which will be consumed in the end of switch case in l.readchar() line...
			tok = token.Token{Type: token.EQ, Literal: "=="}
			// we used Literal:"==" but if we stored ch :=l.ch then we can do Literal:string(ch)+string(l.ch) why to do that becz it makes us abstract the 2chartoken function because we will create a signature function which can also be used with != without any changes !
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
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
	// additionals ==, !, !=, -, /, *, <, >
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERIK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)

	case 0:
		// tok = newToken(token.EOF, "") as "" is not a byte we will have to do...
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IdentLookUp(tok.Literal) // @
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT // @
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) { // ex foobar -> start with ch = f l.position = 0 later o,o,b,a,r, now r is a character so technically it will still run l.readChar which will be ch=0 and position will be increased by 1(thus pointing 0) then it will exit...
		l.readChar() // continuosly marching forward until boom not a char
	}
	return l.input[pos:l.position] // l.position exclusive
}

func (l *Lexer) readNumber() string { // we will be storing number as string because tok.Literal only accepts string while parsing we will convert !
	pos := l.position
	for isDigit(l.ch) { // ex foobar -> start with ch = f l.position = 0 later o,o,b,a,r, now r is a character so technically it will still run l.readChar which will be ch=0 and position will be increased by 1(thus pointing 0) then it will exit...
		l.readChar() // continuosly marching forward until boom not a char
	}
	return l.input[pos:l.position] // l.position exclusive
}

func isLetter(ch byte) bool {
	return ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch == '_' || ch == '@' // @
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
	// note we use newToken helper only for +,; etc so not worry about ch being string we explicitly handled that via tok.Literal way
}

func (l *Lexer) eatWhiteSpace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\r' || l.ch == '\t' {
		l.readChar()
	}
}
