package evaluator

import (
	"github.com/Skarlso/went/ast"
	"github.com/Skarlso/went/object"
)

// could this be Generalized?
func Eval(node ast.Node) object.Object {
	// TODO: this walks the AST but why? It just doesn't care what it encounters
	// as long as we get Int?
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

// It's looking for the thing?
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		// This is a recursion here!
		result = Eval(statement)
	}

	return result
}
