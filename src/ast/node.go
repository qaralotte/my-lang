package ast

type (
	Node interface{}

	// Expr Expression
	Expr interface {
		Node
		expr()
	}

	// Stmt statement
	Stmt interface {
		Node
		stmt()
	}
)
