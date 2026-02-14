/*

The subset of the Monk language we’re going to lex in our first step looks like this:
let five = 5;
let ten = 10;
let add = fn(x, y) {
	x + y;
};
let result = add(five, ten);

*/

package token

type TokenType string // TokenType is not about limiting value but limiting meaning

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// IDENTIFIERS/VARIABLES + LITERALS
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	//OPERATORS
	ASSIGN  = "="
	PLUS    = "+"
	MINUS   = "-"
	BANG    = "!"
	ASTERIK = "*"
	SLASH   = "/"
	LT      = "<"
	GT      = ">"
	EQ      = "=="
	NOT_EQ  = "!="

	//DELIMETERS
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	LET      = "LET"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var Keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func IdentLookUp(ident string) TokenType {
	if tokType, OkFlag := Keywords[ident]; OkFlag {
		return tokType
	} //else {
	return IDENT
	//}
}
