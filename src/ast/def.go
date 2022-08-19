package ast

import (
	"fmt"
	"my-compiler/token"
)

// 定义方法参数
func (p *Parser) defFnArgs() (args []Object) {

	// ([a, b, c])
	for p.Token != token.RPAREN {
		// ([a], ...)
		// 对于每一个局部变量定义，都应该创建新的
		name := p.require(token.IDENTITY, true)
		va := NewVariable(name)
		args = append(args, va)

		// (a[,] ...)
		// 如果是逗号，则说明后面还有参数定义
		if p.Token == token.COMMA {
			p.nextToken()
			p.require(token.IDENTITY, false)
		}
	}
	return
}

// 定义方法
func (p *Parser) defFn() {
	// [fn] name(...) {...}
	p.require(token.FN, true)

	// fn [name](...) {...}
	name := p.require(token.IDENTITY, true)
	obj := p.Objects.findObject(name)
	if obj != nil {
		panic(fmt.Sprintf("错误: 重复定义的名称 %s", name))
	}

	// 添加方法进对象表
	fn := NewFunction(name, p.Objects)
	p.Objects.add(fn)

	// fn name[(...)] {...}
	args := make([]Object, 0)
	if p.Token == token.LPAREN {
		p.nextToken()
		args = p.defFnArgs()
		p.require(token.RPAREN, true)
	}
	fn.Objects.addBatch(args)

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
