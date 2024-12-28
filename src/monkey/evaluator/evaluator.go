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
		return evalStatements(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		if n.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpr(n.Operator, right)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) (result object.Object) {
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return
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
