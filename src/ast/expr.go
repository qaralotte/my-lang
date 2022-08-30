package ast

import (
	"fmt"
	"my-lang/token"
	"reflect"
)

type (
	// BinaryExpr 二元表达式
	BinaryExpr struct {
		Left  Expr
		Op    int
		Right Expr
	}

	// LitExpr 字面量
	LitExpr struct {
		Type
		Lit string
	}

	// IdentityExpr 变量
	IdentityExpr struct {
		Object
	}

	// BlockExpr 块状语句
	BlockExpr struct {
		Toks []token.Token
	}

	// CallFnExpr 调用方法
	CallFnExpr struct {
		Fn     *Function
		Params []Expr
	}
)

func (*BinaryExpr) expr()   {}
func (*LitExpr) expr()      {}
func (*IdentityExpr) expr() {}
func (*BlockExpr) expr()    {}
func (*CallFnExpr) expr()   {}

// 跳过方法内语句
func (p *Parser) block() (toks []token.Token) {
	p.require(token.LBRACE, true)

	level := 0 // block 层数
	for p.Token().Type != token.RBRACE || level != 0 {
		toks = append(toks, p.Token())

		if p.Token().Type == token.LBRACE {
			level += 1
		}

		if p.Token().Type == token.RBRACE {
			level -= 1
		}

		p.next()
	}

	p.require(token.RBRACE, true)

	return
}

func (p *Parser) line() (toks []token.Token) {
	for p.Token().Type != token.LINEBREAK && !p.IsEnd() {
		toks = append(toks, p.Token())
		p.next()
	}

	return
}

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
	for p.Token().Type != token.RPAREN {
		expr := p.parseExpr(0)
		params = append(params, expr)

		if p.Token().Type == token.COMMA {
			p.next()
		}
	}
	p.require(token.RPAREN, true)

	if len(params) != len(fn.Args) {
		panic(fmt.Sprintf("错误: 参数数量不一致"))
	}

	return &CallFnExpr{
		Fn:     fn,
		Params: params,
	}
}

// 解析 1 为何物, "str" 为何物, a 为何物
func (p *Parser) implExpr() (expr Expr) {

	switch p.Token().Type {
	case token.LPAREN:
		// 括号 (优先计算)
		p.next()
		expr = p.parseExpr(0)

		p.require(token.RPAREN, true)
	case token.IDENTITY:
		// 变量
		obj := p.Objects.FindObject(p.Token().Lit)
		if obj == nil {
			// 如果对象表里没有此对象，直接报错
			panic(fmt.Sprintf("错误: 找不到对象: %s", p.Token().Lit))
		}

		p.next()
		if p.Token().Type == token.LPAREN {
			// a(...)
			return p.callFn(obj)
		}

		p.rollback()
		expr = &IdentityExpr{
			Object: obj,
		}
	case token.INTLIT:
		// 整数
		expr = &LitExpr{
			Type: INT,
			Lit:  p.Token().Lit,
		}
	case token.FLOATLIT:
		// 浮点数
		expr = &LitExpr{
			Type: FLOAT,
			Lit:  p.Token().Lit,
		}
	case token.STRINGLIT:
		// 字符串
		expr = &LitExpr{
			Type: STRING,
			Lit:  p.Token().Lit,
		}
	case token.TRUE, token.FALSE:
		// 布尔值
		expr = &LitExpr{
			Type: BOOL,
			Lit:  p.Token().Lit,
		}
	case token.LBRACE:
		// 块状
		return &BlockExpr{
			Toks: p.block(),
		}
	}

	p.next()
	return
}

// 表达式结尾符
func (p *Parser) endExpr() bool {
	if p.IsEnd() {
		return true
	}

	switch p.Token().Type {
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
func operator(tokType token.Type) int {
	switch tokType {
	case token.PLUS:
		return ADD
	case token.MINUS:
		return SUB
	case token.STAR:
		return MUL
	case token.SLASH:
		return DIV
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

	if p.endExpr() {
		return left
	}
	// 1 [+] 2 + 3
	op := operator(p.Token().Type)
	for priority(op) > currentPriority {
		p.next()

		// 1 + [2 + 3]
		right = p.parseExpr(priority(op))
		if right == nil {
			panic("表达式错误")
		}

		//     node
		//    /    \
		// left   right

		left = makeBinary(left, op, right)
		op = operator(p.Token().Type)

		if p.endExpr() {
			return left
		}
	}

	return left
}
