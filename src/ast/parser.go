package ast

import (
	"fmt"
	"my-compiler/data"
	"my-compiler/token"
)

type Parser struct {
	*token.Scanner             // token 扫描器
	token.Token                // 当前定位 token
	Lit            string      // 当前文字标识
	Objects        *ObjectList // 当前对象表
	EndTokens      *data.Stack // token 结束符
}

func NewParser(scanner *token.Scanner) *Parser {
	parser := Parser{
		Scanner:   scanner,
		Objects:   NewObjectList(nil),
		EndTokens: data.NewStack(),
	}

	parser.nextToken()
	parser.EndTokens.Push(token.EOF)

	return &parser
}

// 扫描下一个 token
func (p *Parser) nextToken() {
	p.Token, p.Lit = p.Scanner.ScanNext()
}

// 检查传入的 token, 不符合需要的 token 就 panic
func (p *Parser) require(tok token.Token, autoNext bool) string {
	if p.Token != tok {
		panic(fmt.Sprintf("错误: 需要的 token: %s, 实际提供的 token: %s", token.String(tok), token.String(p.Token)))
	}
	str := p.Lit
	if autoNext {
		p.nextToken()
	}
	return str
}

// ParseStmt 解析语句并整理为语句数组
func (p *Parser) ParseStmt() (Stmt, bool) {
	if p.Token == p.EndTokens.Top() {
		p.EndTokens.Pop()
		return nil, true
	}
	switch p.Token {
	case token.SEMICOLON, token.LINEBREAK:
		// 跳过
		p.nextToken()
	case token.IDENTITY:
		oldParser := p.Copy()
		name := p.Lit

		p.nextToken()
		if p.Token == token.ASSIGN {
			// 如果是等于，则该变量定义且赋值
			p.nextToken()
			return p.parseAssignStatement(name), false
		} else {
			// 如果是别的，则只是表达式
			p.Load(oldParser)
			return p.parseExprStatement(), false
		}
	case token.FN:
		// 方法定义
		p.defFn()
	case token.RETURN:
		// 退出作用域并返回结果
		return p.parseReturnStatement(), false
	case token.PRINT:
		return p.parsePrintStatement(), false
	default:
		// 表达式
		return p.parseExprStatement(), false
	}
	return nil, false
}

func (p *Parser) Copy() Parser {
	return Parser{
		Scanner:   p.Scanner.Copy(),
		Token:     p.Token,
		Lit:       p.Lit,
		Objects:   p.Objects.Copy(),
		EndTokens: p.EndTokens.Copy(),
	}
}

func (p *Parser) Load(np Parser) {
	p.Scanner.Load(np.Scanner)
	p.Token = np.Token
	p.Lit = np.Lit
	p.Objects.Load(np.Objects)
	p.EndTokens.Load(np.EndTokens)
}
