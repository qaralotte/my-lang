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

	// Def define
	Def interface {
		Node
		def()
	}
)

// ----- Statement -----

// ----- Define -----

// DefVar 定义变量
type DefVar struct {
	Def
	Name  string
	Value Expr
}

//

func (*DefVar) def() {}
