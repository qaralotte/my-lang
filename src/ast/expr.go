package ast

import (
	"fmt"
	"my-compiler/token"
	"reflect"
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
	Params []Expr
}

func (*BinaryExpr) expr()   {}
func (*LitExpr) expr()      {}
func (*IdentityExpr) expr() {}
func (*CallFnExpr) expr()   {}

// 调用方法
func (p *Parser) callFn(obj Object) *CallFnExpr {

	name, _ := getObjectField(obj, "Name")

	// 如果调用的对象不是方法
	if reflect.TypeOf(obj).String() != "*ast.Function" {
		panic(fmt.Sprintf("错误: 无法调用方法 %s, 因为 %s 不是方法", name, name))
	}
	fn := obj.(*Function)

	params := make([]Expr, 0)

	p.require(token.LPAREN, true)
	for p.Token != token.RPAREN {
		expr := p.parseExpr(0)
		params = append(params, expr)

		if p.Token == token.COMMA {
			p.nextToken()
		}
	}

	if len(params) != len(fn.Args) {
		panic(fmt.Sprintf("错误: 参数数量不一致"))
	}

	return &CallFnExpr{
		Fn:     fn,
		Params: params,
	}
}

// 解析 1 为何物, "str" 为何物, a 为何物
func (p *Parser) implExpr() Expr {

	switch p.Token {
	case token.LPAREN:
		// 括号 (优先计算)
		p.nextToken()
		node := p.parseExpr(0)

		p.require(token.RPAREN, false)
		return node
	case token.IDENTITY:
		// 变量
		obj := p.Objects.findObject(p.Lit)
		if obj == nil {
			// 如果对象表里没有此对象，直接报错
			panic(fmt.Sprintf("错误: 找不到对象: %s", p.Lit))
		}

		copyParser := p.Copy()

		p.nextToken()
		if p.Token == token.LPAREN {
			// a(...)
			return p.callFn(obj)
		}

		p.Load(copyParser)

		return &IdentityExpr{
			Object: obj,
		}
	case token.INTLIT:
		// 整数
		return &LitExpr{
			Type: INT,
			Lit:  p.Lit,
		}
	case token.FLOATLIT:
		// 浮点数
		return &LitExpr{
			Type: FLOAT,
			Lit:  p.Lit,
		}
	case token.STRINGLIT:
		// 字符串
		return &LitExpr{
			Type: STRING,
			Lit:  p.Lit,
		}
	case token.TRUE, token.FALSE:
		// 布尔值
		return &LitExpr{
			Type: BOOL,
			Lit:  p.Lit,
		}
	}

	panic(fmt.Sprintf("错误: 表达式未知的 token: %s", token.String(p.Token)))
}

// 表达式结尾符
func (p *Parser) endExpr() bool {
	switch p.Token {
	case token.LINEBREAK, token.SEMICOLON, token.RPAREN, token.EOF, token.COMMA:
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
func (p *Parser) parseExpr(currentPriority int) Expr {
	var left, right Expr

	// 从左边开始解析
	// [1] + 2 + 3
	left = p.implExpr()
	if left == nil {
		return left
	}

	p.nextToken()
	if p.endExpr() {
		return left
	}

	// 1 [+] 2 + 3
	op := operator(p.Token)
	for priority(op) > currentPriority {
		p.nextToken()

		// 1 + [2 + 3]
		right = p.parseExpr(currentPriority)
		if right == nil {
			panic("表达式错误")
		}

		// 向下合并类型
		// tryMerge(&leftType, &rightType)

		// 检查类型运算
		// if !canCalc(op, leftType, rightType) {
		// 	panic(fmt.Sprintf("错误: 不合法的计算 %s %s %s", TypeString(leftType), OperatorString(op), TypeString(rightType)))
		// }

		// 如果是比较运算符, 结果应该是bool类型
		// if isComparedOperator(op) {
		// 	leftType = BOOL
		// }

		//     node
		//    /    \
		// left   right

		left = makeBinary(left, op, right)

		if p.endExpr() {
			return left
		}
	}

	return left
}
