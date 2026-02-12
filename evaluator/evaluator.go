package evaluator

import (
	"fmt"
	"monk/ast"
	"monk/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{Value: "null"}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	//expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:

		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env) // return ni andar nu expression decode kariye
		if isError(val) {                  // to avoid bubbling it further down
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// we need to bound node.Name with val using env

		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{Env: env, Body: node.Body, Parameter: node.Parameter}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env) // we got the arguments
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	function, ok := fn.(*object.Function)

	if !ok {
		return newError("not a function : %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)

	evaluated := Eval(function.Body, extendedEnv)

	return unWrapReturnValue(evaluated)
	// we unWrap otherwise the return object will just bubble up
}

func unWrapReturnValue(obj object.Object) object.Object {
	return_value, ok := obj.(*object.ReturnValue)
	if ok {
		return return_value.Value
	}
	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {

	env := object.NewEnclosedEnvironment(fn.Env)

	// lets add args now...
	for ParamIdx, Param := range fn.Parameter {
		env.Set(Param.Value, args[ParamIdx])
		// kemke parameter Identifier hata aney Identifier ma Value hoye je String hoye
	}

	return env
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {

	var result []object.Object

	for _, e := range exps { // left to right arguments evaluation a nature of monkey !

		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated} // return Null i.e array has only 1 object hence we could use args[0] above...
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}
	return newError("Identifier Not Found : %s", node.Value)
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {

	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right) // look its pointer comparison as left and right use same TRUE and FALSE pointeeeers !also thus it is faster then Integer where we compare value !

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("Type Mismatch : %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("Unknown Operator : %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}

	case "-":
		return &object.Integer{Value: leftVal - rightVal}

	case "*":
		return &object.Integer{Value: leftVal * rightVal}

	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return newError("Unknown Operator : %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown Operator : %s %s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown Operator : -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)
		switch result := result.(type) {

		case *object.ReturnValue:
			return result.Value // here its Value is an Object like our result as return is just a warapper around an object , so check if the object as return wrapper and if it does unwrap it and return that object as result.
		case *object.Error:
			return result
		}
	}
	return result // returning last result !
}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {

		case *object.ReturnValue:
			return result // let the error bubble up baby !
		case *object.Error:
			return result
		}

	}
	return result // exactly similar to evalProgram but this propogates the return object upward so it can also be used by outer scope  whereas evalProgram unwraps it and returns it
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
} // ... is variadic parameter which accepts any number of values of any type(any type because its interface type )

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// one of the best lines :

// Closure works because :
//Function literals capture the environment where they are defined, not where they are called.

//And when called:

//We extend that captured environment — not the current one.
