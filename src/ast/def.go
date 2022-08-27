package ast

import (
	"fmt"
	"my-compiler/token"
)

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
	// [fn] name(...)
	p.require(token.FN, true)

	// fn [name](...)
	name := p.require(token.IDENTITY, true)
	obj := p.Objects.findObject(name)
	if obj != nil {
		panic(fmt.Sprintf("错误: 重复定义的名称 %s", name))
	}

	// 添加方法进对象表
	fn := NewFunction(name)
	p.Objects.Add(fn)

	// fn name[(...)]
	if p.Token == token.LPAREN {
		p.nextToken()
		fn.Args = p.defFnArgs()
		p.require(token.RPAREN, true)
	}

	// 跳过块状方法
	fn.Block, _ = p.SkipBlock()
}
