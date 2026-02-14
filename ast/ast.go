package ast

/*

let x = 10;
let y = 15;
let add = fn(a, b) {
	return a + b;
};


*/

import (
	"bytes"
	"monk/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string // only for debugging no logic
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

//--------------Let statement : let x = 5---------------------

type LetStatement struct {
	Name  *Identifier // a node
	Value Expression  // 5 a node
	Token token.Token // we have this only for debugging

}

func (ls *LetStatement) StatementNode() {

}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

//-------------------------Identifier----------------------------

type Identifier struct {
	Token token.Token
	Value string // identifier in future do produce a value...so to use it then we defined this value
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) ExpressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

//---------------------Return Statement--------------------------

type ReturnStatement struct {
	Token       token.Token // we have this only for debugging
	ReturnValue Expression  // 5 a node
}

func (rs *ReturnStatement) StatementNode() {

}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

//--------------------Expression---------------------------------

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) StatementNode() {}

func (es *ExpressionStatement) String() string {

	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

//------------------------IntegerLiteral------------------------

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) ExpressionNode() {}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

//-------------------------prefixExpression----------------------

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	// (-5)
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

//--------------------Infix Expression---------------------

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) ExpressionNode() {}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()

}

// -----------------------Boolean Expression-------------------
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) ExpressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

//--------------------------IF EXPRESSION---------------

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) ExpressionNode() {}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // {
	Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) StatementNode() {}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

//--------------------FUNCTION LITERAL----------

type FunctionLiteral struct {
	Token     token.Token
	Parameter []*Identifier
	Body      *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	params := []string{}
	for _, p := range fl.Parameter {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())
	return out.String()
}

func (fl *FunctionLiteral) ExpressionNode() {}

//-;----------------Function CALL EXPRESSION-----------

type CallExpression struct {
	Token token.Token // (

	Function Expression // add i.e identifier expression or fn add(x,y){x+y;}

	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) ExpressionNode() {}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Function.String())
	out.WriteString("(")

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")
	return out.String()
}

//--------------------------String-----------------------------

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) ExpressionNode() {}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Value
}
