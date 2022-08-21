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
func (p *Parser) identity() Object {

	switch p.Token {
	case token.IDENTITY:
		// 变量 (且必须是变量)
		va := p.Objects.findObject(p.Lit)
		if va == nil {
			// 如果变量表里没有此变量，则新建一个变量
			va = NewVariable(p.Lit)
			p.Objects.add(va)
		}

		p.nextToken()
		return va
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.String(p.Token)))
}

func (p *Parser) assign() {
	va := p.identity().(*Variable)

	op := p.Token
	p.nextToken()

	if op == token.EOF || op == token.LINEBREAK {
		panic("非法的变量定义语法")
	}

	if op == token.ASSIGN {
		expr, typ := p.parseExpr(0)
		va.Type = typ
		va.Value = expr
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

// 赋值语句
func (p *Parser) parseAssignStatement() *AssignStmt {

	left := p.identity().(*Variable)

	op := p.Token
	p.nextToken()

	if op == token.EOF || op == token.LINEBREAK {
		panic("非法的变量定义语法")
	}

	right, typ := p.parseExpr(0)
	left.Type = typ

	if op == token.ASSIGN {
		return &AssignStmt{
			Object: left,
			Right:  right,
		}
	} else {
		panic(fmt.Sprintf("非法符号: %s", token.String(op)))
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
