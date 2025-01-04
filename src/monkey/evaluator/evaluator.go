package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return toBooleanObject(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpr(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)
		return evalInfixExpr(n.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStmt(n)
	case *ast.IfExpression:
		return evalIfExpression(n)
	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue)
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func toBooleanObject(val bool) object.Object {
	if val {
		return TRUE
	}
	return FALSE
}

func evalProgram(stmts []ast.Statement) (result object.Object) {
	for _, stmt := range stmts {
		result = Eval(stmt)
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return
}

func evalBlockStmt(block *ast.BlockStatement) (result object.Object) {
	for _, stmt := range block.Statements {
		result = Eval(stmt)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}
	return result
}

func evalPrefixExpr(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalMinusOperatorExpr(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpr(right object.Object) object.Object {
	switch right {
	case TRUE, NULL:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalInfixExpr(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfix(operator, left, right)
	case operator == "==":
		return toBooleanObject(left == right)
	case operator == "!=":
		return toBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIntegerInfix(operator string, left, right object.Object) object.Object {
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
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return toBooleanObject(leftVal < rightVal)
	case ">":
		return toBooleanObject(leftVal > rightVal)
	case "==":
		return toBooleanObject(leftVal == rightVal)
	case "!=":
		return toBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
