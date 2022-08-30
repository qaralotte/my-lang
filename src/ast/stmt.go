package ast

import (
	"fmt"
	"my-lang/token"
)

type (

	// ExprStmt 单表达式的语句
	ExprStmt struct {
		Expr
	}

	// AssignStmt 赋值语句
	AssignStmt struct {
		Name  string
		Value Expr
	}

	// PrintStmt 打印 (暂时) deprecated
	PrintStmt struct {
		Expr
	}

	// ReturnStmt 返回语句
	ReturnStmt struct {
		Expr
	}

	// IfStmt 选择语句
	IfStmt struct {
		Cond      Expr
		TrueBody  []token.Token
		FalseBody []token.Token
	}

	// ForStmt 循环语句
	ForStmt struct {
		Cond Expr
		Body []token.Token
	}
)

func (*ExprStmt) stmt()   {}
func (*AssignStmt) stmt() {}
func (*PrintStmt) stmt()  {}
func (*ReturnStmt) stmt() {}
func (*IfStmt) stmt()     {}
func (*ForStmt) stmt()    {}

// 获取当前token的identity
func (p *Parser) identity() (obj Object) {
	switch p.Token().Type {
	case token.IDENTITY:
		// 变量 (且必须是变量)
		obj = p.Objects.FindObject(p.Token().Lit)
		p.next()
		return
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.TypeString(p.Token().Type)))
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

	expr := p.parseExpr(0)
	return &ReturnStmt{
		Expr: expr,
	}
}

// 分支选择语句
func (p *Parser) parseIfStatement() *IfStmt {

	p.require(token.IF, true)

	// 条件
	cond := p.parseExpr(0)

	// 如果条件成立则进入
	trueBody := p.block()

	// 如果条件非成立则进入
	var falseBody []token.Token = nil
	if p.Token().Type == token.ELSE {
		p.next()
		falseBody = p.block()
	}

	return &IfStmt{
		Cond:      cond,
		TrueBody:  trueBody,
		FalseBody: falseBody,
	}
}

func (p *Parser) parseForStatement() *ForStmt {
	p.require(token.FOR, true)

	// 条件
	cond := p.parseExpr(0)

	// 如果条件成立则进入
	body := p.block()

	return &ForStmt{
		Cond: cond,
		Body: body,
	}
}
