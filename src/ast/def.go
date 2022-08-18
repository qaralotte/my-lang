package ast

import (
	"fmt"
	"my-compiler/token"
)

func (p *Parser) defFn() {
	// [fn] name(...) {...}
	p.require(token.FN, true)

	// fn [name](...) {...}
	name := p.require(token.IDENTITY, true)
	obj := p.Objects.findObject(name)
	if obj != nil {
		panic(fmt.Sprintf("错误: 重复定义的名称 %s", name))
	}

	// fn name[(...)] {...}
	args := make([]Object, 0) // todo DEFAULT_SIZE
	if p.Token == token.LPAREN {
		// todo 函数参数定义
		p.nextToken()
		p.require(token.RPAREN, true)
	}

	// 添加方法进对象表
	fn := NewFunction(name, args, p.Objects)
	p.Objects.add(fn)

	// fn name(...) [{]...}
	p.require(token.LBRACE, true)
	// 设置当前对象表为方法内局部对象表
	p.Objects = fn.Objects
	// todo GET AST
	fn.Stmts = p.ParseStmts(token.RBRACE)
	// 设置当前对象表为方法外全局对象表
	p.Objects = p.Objects.get(0).(*Channel).Next
	// fn name(...) {...[}]
	p.require(token.RBRACE, true)
}
