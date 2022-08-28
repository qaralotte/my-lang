package ast

import (
	"my-compiler/token"
)

// 定义方法参数
func (p *Parser) defFnArgs() (args []string) {

	// ([a, b, c])
	for p.Token().Type != token.RPAREN {
		// ([a], ...)
		name := p.require(token.IDENTITY, true)
		args = append(args, name)

		// (a[,] ...)
		// 如果是逗号，则说明后面还有参数定义
		if p.Token().Type == token.COMMA {
			p.next()
			p.require(token.IDENTITY, false)
		}
	}
	return
}

// 定义方法
func (p *Parser) defFn() {

}
