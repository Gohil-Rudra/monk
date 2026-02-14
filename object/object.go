package object

import (
	"bytes"
	"fmt"
	"monk/ast"
	"strings"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
)

type ObjectType string

type Object interface { // similar to Token by differs that this is an interface and that was a struct...to not cram boolean and Int together
	Type() ObjectType
	Inspect() string
}

//-----------------------Integer-----------------------------------

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

//--------------------------Boolean-----------------------------

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

//---------------------------Null----------------------------------

type Null struct {
	Value string
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

//------------------------Return Value--------------------------

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

//------------------------Error---------------------------------

type Error struct {
	Message string
}

func (er *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (er *Error) Inspect() string {
	return "ERROR :" + er.Message
}

// --------------------------Function-------------------------
type Function struct {
	Env       *Environment // closure !
	Body      *ast.BlockStatement
	Parameter []*ast.Identifier
}

func (fn *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (fn *Function) Inspect() string {

	var out bytes.Buffer

	params := []string{}

	for _, p := range fn.Parameter {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString("\n {")
	out.WriteString(fn.Body.String())
	out.WriteString("}")

	return out.String()
}

//-------------------------Strings----------------------------

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}
