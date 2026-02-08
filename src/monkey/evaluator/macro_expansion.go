package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}
	for i, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			definitions = append(definitions, i)
		}
	}
	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		defIdx := definitions[i]
		program.Statements = append(program.Statements[:defIdx], program.Statements[defIdx+1:]...)
	}
}

func isMacroDefinition(statement ast.Statement) bool {
	letStmt, ok := statement.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letStmt.Value.(*ast.MacroLiteral)
	return ok
}

func addMacro(statement ast.Statement, env *object.Environment) {
	letStmt, _ := statement.(*ast.LetStatement)
	macroLiteral, _ := letStmt.Value.(*ast.MacroLiteral)
	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}
	env.Set(letStmt.Name.Value, macro)
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(n ast.Node) ast.Node {
		callExpr, ok := n.(*ast.CallExpression)
		if !ok {
			return n
		}
		macro, ok := isMacroCall(callExpr, env)
		if !ok {
			return n
		}
		args := quoteArgs(callExpr)
		evalEnv := extendMacroEnv(macro, args)
		evaluated := Eval(macro.Body, evalEnv)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("Only returning AST-nodes from macros is supported")
		}
		return quote.Node
	})
}

func isMacroCall(callExpr *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := callExpr.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}
	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}
	return macro, true
}

func quoteArgs(callExpr *ast.CallExpression) (args []*object.Quote) {
	for _, a := range callExpr.Arguments {
		args = append(args, &object.Quote{Node: a})
	}
	return
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)
	for paramIdx, param := range macro.Parameters {
		extended.Set(param.Value, args[paramIdx])
	}
	return &extended
}
