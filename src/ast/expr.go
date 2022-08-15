package ast

import (
	"fmt"
	"my-compiler/token"
	"my-compiler/variable"
)

// BinaryExpr 二元表达式
type BinaryExpr struct {
	Left  Expr
	Op    int
	Right Expr
}

// LitExpr 字面量
type LitExpr struct {
	Lit string
}

// IdentityExpr 变量
type IdentityExpr struct {
	Var *variable.Variable
}

func (*BinaryExpr) expr()   {}
func (*LitExpr) expr()      {}
func (*IdentityExpr) expr() {}

// 解析 1 为何物, "str" 为何物, a 为何物
func implExpr() (Expr, variable.Type) {

	switch globalParser.Token {
	case token.LPAREN:
		// 括号 (优先计算)
		globalParser.NextToken()
		node, typ := parseExpr(0)

		globalParser.Require(token.RPAREN, false)
		return node, typ
	case token.IDENTITY:
		// 变量
		va := variable.FindVariable(globalParser.Variables, globalParser.Lit)
		if va == nil {
			// 如果变量表里没有此变量，直接报错
			panic(fmt.Sprintf("错误: 找不到变量: %s", globalParser.Lit))
		}
		return &IdentityExpr{
			Var: va,
		}, va.Type
	case token.INTLIT:
		// 整数
		var node LitExpr
		node.Lit = globalParser.Lit
		return &node, variable.INT
	case token.FLOATLIT:
		// 数字
		var node LitExpr
		node.Lit = globalParser.Lit
		return &node, variable.FLOAT
	case token.STRINGLIT:
		// 字符串
		var node LitExpr
		node.Lit = globalParser.Lit
		return &node, variable.STRING
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.String(globalParser.Token)))
}

// 表达式结尾符
func endExpr() bool {
	switch globalParser.Token {
	case token.LINEBREAK, token.SEMICOLON, token.RPAREN, token.EOF:
		return true
	}
	return false
}

// 运算符
const (
	NOP int = iota
	ADD
	SUB
	MUL
	DIV
)

func OperatorString(op int) string {
	switch op {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	}
	return "nop"
}

// 运算符优先级
func priority(op int) int {
	switch op {
	case ADD, SUB:
		return 3
	case MUL, DIV:
		return 4
	}
	return 0
}

// 根据 Token 转换成对应的 Syntax
func operator(tok token.Token) int {
	switch tok {
	case token.PLUS:
		return ADD
	case token.MINUS:
		return SUB
	case token.STAR:
		return MUL
	case token.SLASH:
		return DIV
	}
	return NOP
}

func makeBinary(left Expr, op int, right Expr) *BinaryExpr {
	return &BinaryExpr{
		Left:  left,
		Op:    op,
		Right: right,
	}
}

// ParseExpr 表达式解析
func parseExpr(currentPriority int) (Expr, variable.Type) {
	var left, right Expr
	var leftType, rightType variable.Type

	// 从左边开始解析
	// [1] + 2 + 3
	left, leftType = implExpr()
	if left == nil {
		return left, leftType
	}

	globalParser.NextToken()
	if endExpr() {
		return left, leftType
	}

	// 1 [+] 2 + 3
	op := operator(globalParser.Token)
	for priority(op) > currentPriority {
		globalParser.NextToken()

		// 1 + [2 + 3]
		right, rightType = parseExpr(currentPriority)
		if right == nil {
			panic("表达式错误")
		}

		// 向下合并类型
		variable.Merge(&leftType, &rightType)

		// 检查类型运算
		if !variable.CanCalc(leftType, rightType) {
			panic(fmt.Sprintf("错误: 类型 %s 不能与 类型 %s 计算", variable.TypeString(leftType), variable.TypeString(rightType)))
		}

		//     node
		//    /    \
		// left   right

		left = makeBinary(left, op, right)

		if endExpr() {
			return left, leftType
		}
	}

	return left, leftType
}
