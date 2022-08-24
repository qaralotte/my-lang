package ast

import (
	"fmt"
	"my-compiler/token"
)

type (

	// ExprStmt 单表达式的语句
	ExprStmt struct {
		Expr
	}

	// AssignStmt 赋值表达式
	AssignStmt struct {
		Name  string
		Value Expr
	}

	// PrintStmt 打印 (暂时) deprecated
	PrintStmt struct {
		Expr
	}

	ReturnStmt struct {
		Expr
	}
)

func (*ExprStmt) stmt()   {}
func (*AssignStmt) stmt() {}
func (*PrintStmt) stmt()  {}
func (*ReturnStmt) stmt() {}

// 获取当前token的identity
func (p *Parser) identity() (obj Object, name string) {
	switch p.Token {
	case token.IDENTITY:
		// 变量 (且必须是变量)
		obj = p.Objects.findObject(p.Lit)
		name = p.Lit
		p.nextToken()
		return
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.String(p.Token)))
}

// 表达式语句 (语句里只包含表达式)
func (p *Parser) parseExprStatement() *ExprStmt {
	expr := p.parseExpr(0)
	return &ExprStmt{
		expr,
	}
}

// 打印语句
func (p *Parser) parsePrintStatement() *PrintStmt {

	p.require(token.PRINT, true)

	expr := p.parseExpr(0)
	return &PrintStmt{
		Expr: expr,
	}

}

// 赋值语句
func (p *Parser) parseAssignStatement(name string) *AssignStmt {
	expr := p.parseExpr(0)

	return &AssignStmt{
		Name:  name,
		Value: expr,
	}
}

// 退出方法并返回值
func (p *Parser) parseReturnStatement() *ReturnStmt {

	p.require(token.RETURN, true)

	parent := p.Objects.getParentObject()
	if parent == nil {
		panic(fmt.Sprintf("错误: 不合法的 return 使用场合"))
	}

	expr := p.parseExpr(0)
	return &ReturnStmt{
		Expr: expr,
	}
}
