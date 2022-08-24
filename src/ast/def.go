package ast

import (
	"fmt"
	"my-compiler/token"
)

// 跳过方法内语句
func (p *Parser) skipBlock() {
	level := 0 // block 层数
	for p.Token != token.RBRACE || level != 0 {
		if p.Token == token.LPAREN {
			level += 1
		}

		if p.Token == token.RPAREN {
			level -= 1
		}

		p.nextToken()
	}
}

// 定义方法参数
func (p *Parser) defFnArgs() (args []string) {

	// ([a, b, c])
	for p.Token != token.RPAREN {
		// ([a], ...)
		name := p.require(token.IDENTITY, true)
		args = append(args, name)

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
	p.Objects.Add(fn)

	// fn name[(...)] {...}
	if p.Token == token.LPAREN {
		p.nextToken()
		fn.Args = p.defFnArgs()
		p.require(token.RPAREN, true)
	}

	// fn name(...) [{]...}
	p.require(token.LBRACE, true)

	fn.Parser = p.Copy()
	p.skipBlock()
	// fn name(...) {...[}]
	p.require(token.RBRACE, true)
}

// Deprecated
func (p *Parser) defVar(name string) {
	expr := p.parseExpr(0)

	p.Objects.Add(&Variable{
		Name:  name,
		Value: expr,
	})
}
