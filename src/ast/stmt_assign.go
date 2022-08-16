package ast

import (
	"fmt"
	"my-compiler/object"
	"my-compiler/token"
)

// AssignStmt 赋值表达式
type AssignStmt struct {
	Object object.Object
	Right  Expr
}

func (*AssignStmt) stmt() {}

// 解析 1 为何物, "str" 为何物, a 为何物
func identity() object.Object {

	switch globalParser.Token {
	case token.IDENTITY:
		// 变量 (且必须是变量)
		va := object.FindObject(globalParser.Objects, globalParser.Lit)
		if va == nil {
			// 如果变量表里没有此变量，则新建一个变量
			va = object.NewVariable(globalParser.Lit)
			globalParser.Objects = append(globalParser.Objects, va)
		}

		globalParser.NextToken()
		return va
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.String(globalParser.Token)))
}

// parseAssignStatement 赋值语句
func (p *Parser) parseAssignStatement() *AssignStmt {

	left := identity().(*object.Variable)

	op := p.Token
	p.NextToken()

	if op == token.EOF || op == token.LINEBREAK {
		panic("非法的变量定义语法")
	}

	right, typ := parseExpr(0)
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
