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

	// AssignStmt 赋值表达式 deprecated
	AssignStmt struct {
		Object
		Right Expr
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

func (p *Parser) assign() {
	obj, name := p.identity()

	op := p.Token
	p.nextToken()

	if op == token.EOF || op == token.LINEBREAK {
		panic("非法的变量定义语法")
	}

	if op == token.ASSIGN {
		expr, typ := p.parseExpr(0)
		if obj == nil {
			// 如果对象表没有对象，则新建一个
			obj = &Variable{
				Type:  typ,
				Name:  name,
				Value: expr,
			}
			p.Objects.add(obj)
		} else {
			// 如果对象表有对象，则修改
			va := obj.(*Variable)
			va.Type = typ
			va.Value = expr
		}

	} else {
		panic(fmt.Sprintf("非法符号: %s", token.String(op)))
	}
}

// 表达式语句 (语句里只包含表达式)
func (p *Parser) parseExprStatement() *ExprStmt {
	expr, _ := p.parseExpr(0)
	return &ExprStmt{
		expr,
	}
}

func (p *Parser) parsePrintStatement() *PrintStmt {

	p.require(token.PRINT, true)

	expr, _ := p.parseExpr(0)
	return &PrintStmt{
		Expr: expr,
	}

}

func (p *Parser) parseReturnStatement() *ReturnStmt {

	p.require(token.RETURN, true)

	parent := p.Objects.getParentObject()
	if parent == nil {
		panic(fmt.Sprintf("错误: 不合法的 return 使用场合"))
	}

	expr, _ := p.parseExpr(0)
	return &ReturnStmt{
		Expr: expr,
	}
}
