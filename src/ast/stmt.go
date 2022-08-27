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

	// IfStmt 选择分支语句
	IfStmt struct {
		Cond  Expr
		True  *Parser
		False *Parser
		End   *Parser
	}
)

func (*ExprStmt) stmt()   {}
func (*AssignStmt) stmt() {}
func (*PrintStmt) stmt()  {}
func (*ReturnStmt) stmt() {}
func (*IfStmt) stmt()     {}

// SkipBlock 跳过块状语句
func (p *Parser) SkipBlock() (*Parser, *Parser) {
	p.require(token.LBRACE, true)
	start := p.Copy()

	level := 0 // block 层数
	for p.Token != token.RBRACE || level != 0 {
		if p.Token == token.LBRACE {
			level += 1
		}

		if p.Token == token.RBRACE {
			level -= 1
		}

		p.nextToken()
	}

	p.require(token.RBRACE, true)
	end := p.Copy()

	return start, end
}

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

func (p *Parser) parseIfStatement() *IfStmt {

	p.require(token.IF, true)

	cond := p.parseExpr(0)

	var trueBlock, falseBlock, endPos *Parser = nil, nil, nil
	trueBlock, endPos = p.SkipBlock()

	if p.Token == token.ELSE {
		p.nextToken()
		falseBlock, endPos = p.SkipBlock()
	}

	return &IfStmt{
		Cond:  cond,
		True:  trueBlock,
		False: falseBlock,
		End:   endPos,
	}
}
