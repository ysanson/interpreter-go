package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modFunc ModifierFunc) Node {
	switch n := node.(type) {
	case *Program:
		for i, statement := range n.Statements {
			n.Statements[i], _ = Modify(statement, modFunc).(Statement)
		}
	case *ExpressionStatement:
		n.Expression, _ = Modify(n.Expression, modFunc).(Expression)
	case *InfixExpression:
		n.Left, _ = Modify(n.Left, modFunc).(Expression)
		n.Right, _ = Modify(n.Right, modFunc).(Expression)
	case *PrefixExpression:
		n.Right, _ = Modify(n.Right, modFunc).(Expression)
	case *IndexExpression:
		n.Left, _ = Modify(n.Left, modFunc).(Expression)
		n.Index, _ = Modify(n.Index, modFunc).(Expression)
	case *IfExpression:
		n.Condition, _ = Modify(n.Condition, modFunc).(Expression)
		n.Consequence, _ = Modify(n.Consequence, modFunc).(*BlockStatement)
		if n.Alternative != nil {
			n.Alternative, _ = Modify(n.Alternative, modFunc).(*BlockStatement)
		}
	case *BlockStatement:
		for i, stmt := range n.Statements {
			n.Statements[i], _ = Modify(stmt, modFunc).(Statement)
		}
	case *ReturnStatement:
		n.ReturnValue, _ = Modify(n.ReturnValue, modFunc).(Expression)
	case *LetStatement:
		n.Value, _ = Modify(n.Value, modFunc).(Expression)
	case *FunctionLiteral:
		for i, param := range n.Parameters {
			n.Parameters[i], _ = Modify(param, modFunc).(*Identifier)
		}
		n.Body, _ = Modify(n.Body, modFunc).(*BlockStatement)
	case *ArrayLiteral:
		for i, elt := range n.Elements {
			n.Elements[i], _ = Modify(elt, modFunc).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, val := range n.Pairs {
			key, _ = Modify(key, modFunc).(Expression)
			val, _ = Modify(val, modFunc).(Expression)
			newPairs[key] = val
		}
		n.Pairs = newPairs
	}
	return modFunc(node)
}
