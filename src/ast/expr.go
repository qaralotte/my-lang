package ast

import (
	"fmt"
	"my-compiler/token"
)

// BinaryExpr 二元表达式
type BinaryExpr struct {
	Left  Expr
	Op    int
	Right Expr
}

// LitExpr 字面量
type LitExpr struct {
	Type
	Lit string
}

// IdentityExpr 变量
type IdentityExpr struct {
	Object
}

// CallFnExpr 调用方法
type CallFnExpr struct {
	Fn     *Function
	Params []Object
}

func (*BinaryExpr) expr()   {}
func (*LitExpr) expr()      {}
func (*IdentityExpr) expr() {}
func (*CallFnExpr) expr()   {}

// 解析 1 为何物, "str" 为何物, a 为何物
func (p *Parser) implExpr() (Expr, Type) {

	switch p.Token {
	case token.LPAREN:
		// 括号 (优先计算)
		p.nextToken()
		node, typ := p.parseExpr(0)

		p.require(token.RPAREN, false)
		return node, typ
	case token.IDENTITY:
		// 变量
		obj := p.Objects.findObject(p.Lit)
		if obj == nil {
			// 如果变量表里没有此变量，直接报错
			panic(fmt.Sprintf("错误: 找不到变量: %s", p.Lit))
		}

		typ, isExsit := getObjectField(obj, "Type")
		if !isExsit {
			panic(fmt.Sprintf("错误: 无法获取变量 %s 的类型", p.Lit))
		}

		return &IdentityExpr{
			Object: obj,
		}, typ.Interface().(Type)
	case token.INTLIT:
		// 整数
		return &LitExpr{
			Type: INT,
			Lit:  p.Lit,
		}, INT
	case token.FLOATLIT:
		// 浮点数
		return &LitExpr{
			Type: FLOAT,
			Lit:  p.Lit,
		}, FLOAT
	case token.STRINGLIT:
		// 字符串
		return &LitExpr{
			Type: STRING,
			Lit:  p.Lit,
		}, STRING
	case token.TRUE, token.FALSE:
		// 布尔值
		return &LitExpr{
			Type: BOOL,
			Lit:  p.Lit,
		}, BOOL
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.String(p.Token)))
}

// 表达式结尾符
func (p *Parser) endExpr() bool {
	switch p.Token {
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
	EQ
	NQ
	GT
	GE
	LT
	LE
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
	case EQ:
		return "=="
	case NQ:
		return "!="
	case GT:
		return ">"
	case GE:
		return ">="
	case LT:
		return "<"
	case LE:
		return "<="
	}
	return "nop"
}

// 运算符优先级
func priority(op int) int {
	switch op {
	case EQ, NQ, GT, GE, LT, LE:
		return 2
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
	case token.EQ:
		return EQ
	case token.NQ:
		return NQ
	case token.GT:
		return GT
	case token.GE:
		return GE
	case token.LT:
		return LT
	case token.LE:
		return LE
	}
	return NOP
}

func isComparedOperator(op int) bool {
	switch op {
	case EQ, NQ, GT, GE, LT, LE:
		return true
	}
	return false
}

func makeBinary(left Expr, op int, right Expr) *BinaryExpr {
	return &BinaryExpr{
		Left:  left,
		Op:    op,
		Right: right,
	}
}

// ParseExpr 表达式解析
func (p *Parser) parseExpr(currentPriority int) (Expr, Type) {
	var left, right Expr
	var leftType, rightType Type

	// 从左边开始解析
	// [1] + 2 + 3
	left, leftType = p.implExpr()
	if left == nil {
		return left, leftType
	}

	p.nextToken()
	if p.endExpr() {
		return left, leftType
	}

	// 1 [+] 2 + 3
	op := operator(p.Token)
	for priority(op) > currentPriority {
		p.nextToken()

		// 1 + [2 + 3]
		right, rightType = p.parseExpr(currentPriority)
		if right == nil {
			panic("表达式错误")
		}

		// 向下合并类型
		tryMerge(&leftType, &rightType)

		// 检查类型运算
		if !canCalc(op, leftType, rightType) {
			panic(fmt.Sprintf("错误: 不合法的计算 %s %s %s", TypeString(leftType), OperatorString(op), TypeString(rightType)))
		}

		// 如果是比较运算符, 结果应该是bool类型
		if isComparedOperator(op) {
			leftType = BOOL
		}

		//     node
		//    /    \
		// left   right

		left = makeBinary(left, op, right)

		if p.endExpr() {
			return left, leftType
		}
	}

	return left, leftType
}
