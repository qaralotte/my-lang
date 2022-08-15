package ast

// ExprStmt 单表达式的语句
type ExprStmt struct {
	Expr
}

func (*ExprStmt) stmt() {}

// ParseExprStatement 表达式语句 (语句里只包含表达式)
func (p *Parser) parseExprStatement() *ExprStmt {
	expr, _ := parseExpr(0)
	return &ExprStmt{
		expr,
	}
}
