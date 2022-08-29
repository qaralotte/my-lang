package ast

import (
	"fmt"
	"my-lang/token"
)

type Parser struct {
	Tokens []token.Token // 定位 token
	Offset int           // 解析 token 的索引

	Objects *ObjectList // 对象表
}

func NewParser(toks []token.Token, objs *ObjectList) *Parser {
	toks = append(toks, token.EmptyToken(token.EOF))
	parser := Parser{
		Tokens: toks,
		Offset: 0,

		Objects: objs,
	}

	return &parser
}

func (p *Parser) Token() token.Token {
	return p.Tokens[p.Offset]
}

// 下一个 token
func (p *Parser) next() {
	p.Offset += 1
}

// 回滚一个 token
func (p *Parser) rollback() {
	p.Offset -= 1
}

// 检查传入的 token, 不符合需要的 token 就 panic
func (p *Parser) require(tokType token.Type, autoNext bool) string {
	if p.Token().Type != tokType {
		panic(fmt.Sprintf("错误: 需要的 token: %s, 实际提供的 token: %s", token.TypeString(tokType), token.TypeString(p.Token().Type)))
	}
	str := p.Token().Lit
	if autoNext {
		p.next()
	}
	return str
}

func (p *Parser) IsEnd() bool {
	return p.Token().Type == token.EOF
}

// ParseStmt 解析语句并整理为语句数组
func (p *Parser) ParseStmt() Stmt {

	if p.IsEnd() {
		return nil
	}

	switch p.Token().Type {
	case token.SEMICOLON, token.LINEBREAK:
		// 跳过
		p.next()
	case token.IDENTITY:
		// 变量声明
		name := p.Token().Lit

		// 记录此刻的 token 索引
		startOffset := p.Offset

		p.next()
		if p.Token().Type == token.ASSIGN {
			// [a = ...]
			// 变量的定义与赋值
			p.next()
			return p.parseAssignStatement(name)
		} else if p.Token().Type == token.LPAREN {
			// [a(...)]
			p.next()
			args := p.defFnArgs()

			if p.Token().Type == token.ASSIGN {
				// [a(...) = ...]
				p.next()
				p.defFn(name, args)
			} else {
				// [a(...) + 1]
				p.Offset = startOffset
				return p.parseExprStatement()
			}

		} else {
			// [a + 1]
			// 表达式
			p.Offset = startOffset
			return p.parseExprStatement()
		}
	case token.RETURN:
		// 退出作用域并返回结果
		return p.parseReturnStatement()
	case token.PRINT:
		return p.parsePrintStatement()
	case token.IF:
		return p.parseIfStatement()
	default:
		// 表达式
		return p.parseExprStatement()
	}
	return nil
}
