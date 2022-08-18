package ast

import (
	"fmt"
	"my-compiler/token"
)

type Parser struct {
	*token.Scanner // token 扫描器

	token.Token        // 当前定位的 token
	Lit         string // 当前 token 的文字标识

	Stmts   []Stmt      // 语句
	Objects *ObjectList // 对象表
}

func NewParser(scanner *token.Scanner) *Parser {
	var parser Parser

	// 初始化扫描器
	parser.init(scanner)

	// 解析语句
	parser.Stmts = parser.ParseStmts(token.EOF)

	return &parser
}

// 装载扫描器并设置全局解析器为当前解析器
func (p *Parser) init(s *token.Scanner) {
	p.Scanner = s
	p.nextToken()

	// 根对象表第一个元素必定是 nil
	p.Objects = NewObjectList(nil)
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

// ParseStmts 解析语句并整理为语句数组
func (p *Parser) ParseStmts(end token.Token) (stmts []Stmt) {
	for p.Token != end {
		switch p.Token {
		case token.SEMICOLON, token.LINEBREAK:
			// 跳过
			p.nextToken()
			continue
		case token.IDENTITY:
			// 变量 (定义 or 赋值)
			stmts = append(stmts, p.parseAssignStatement())
		case token.FN:
			// 方法定义
			p.defFn()
		default:
			// 表达式
			stmts = append(stmts, p.parseExprStatement())
		}
	}
	return
}
