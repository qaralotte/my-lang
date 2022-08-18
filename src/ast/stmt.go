package ast

import (
	"fmt"
	"my-compiler/token"
)

// AssignStmt 赋值表达式
type AssignStmt struct {
	Object
	Right Expr
}

func (*AssignStmt) stmt() {}

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

// ParseExprStatement 表达式语句 (语句里只包含表达式)
func (p *Parser) parseExprStatement() *ExprStmt {
	expr, _ := p.parseExpr(0)
	return &ExprStmt{
		expr,
	}
}

// parseAssignStatement 赋值语句
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
