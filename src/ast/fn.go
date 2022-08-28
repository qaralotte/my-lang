package ast

import "my-lang/token"

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
	p.require(token.RPAREN, true)

	return
}

// 定义方法
func (p *Parser) defFn(name string, args []string) {

	var body []token.Token
	if p.Token().Type == token.LBRACE {
		body = p.block()
		p.require(token.RBRACE, true)
	} else {
		// 偷偷加个return
		body = append(body, token.EmptyToken(token.RETURN))
		body = append(body, p.line()...)
		// 行格式不用加载下一个 token (避免 EOF)
	}

	fn := &Function{
		Name: name,
		Args: args,
		Body: body,
	}

	p.Objects.Add(fn)

}
