package ast

import (
	"fmt"
	"my-compiler/token"
	"my-compiler/variable"
)

// AssignStmt 赋值表达式
type AssignStmt struct {
	Var   *variable.Variable
	Right Expr
}

func (*AssignStmt) stmt() {}

// 解析 1 为何物, "str" 为何物, a 为何物
func identity() *variable.Variable {

	switch globalParser.Token {
	case token.IDENTITY:
		// 变量 (且必须是变量)
		va := variable.FindVariable(globalParser.Variables, globalParser.Lit)
		if va == nil {
			// 如果变量表里没有此变量，则新建一个变量
			va = variable.NewVariable(globalParser.Lit)
			globalParser.Variables = append(globalParser.Variables, va)
		}

		globalParser.NextToken()
		return va
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.String(globalParser.Token)))
}

// parseAssignStatement 赋值语句
func (p *Parser) parseAssignStatement() *AssignStmt {

	left := identity()

	op := p.Token
	p.NextToken()

	right := parseExpr(0)

	if op == token.ASSIGN {
		return &AssignStmt{
			Var:   left,
			Right: right,
		}
	} else {
		panic("错误的赋值语句")
	}

}
