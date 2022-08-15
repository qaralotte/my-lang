package token

import (
	"os"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	src []byte

	ch           rune // 当前读取的字符 (utf-8)
	offset       int  // 当前偏移位置
	inlineOffset int  // 在行内偏移位置
	lineOffset   int  // 偏移行位置
}

// end of file
const eof = -1

func NewScanner(path string) (s *Scanner) {
	var scanner Scanner

	// 读取文件内容并存进 src
	byt, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	scanner.src = byt

	scanner.next()
	return &scanner
}

// 是否 offset 到达 EOF
func (s *Scanner) isEOF() bool {
	return s.offset >= len(s.src)
}

// 读取下一个字符并偏移 offset
func (s *Scanner) next() {
	if !s.isEOF() {
		// 切片转码成 utf8
		r, w := utf8.DecodeRune(s.src[s.offset:])
		s.ch = r
		s.offset += w
		s.inlineOffset += w
	} else {
		s.ch = eof
	}
}

// 更新行偏移数据, 即行偏移 += 1, 行内偏移重新 = 0
func (s *Scanner) updateLn() {
	s.lineOffset += 1
	s.inlineOffset = 0
}

// 跳过空格 or 换行 or 制表符
func (s *Scanner) skipSpace() {
	for unicode.IsSpace(s.ch) {
		// 换行符是 token.LINEBREAK, 不能跳过
		if s.ch == '\n' {
			return
		}
		s.next()
	}
}

func (s *Scanner) isKeyword(identity string) Token {
	for _, keyword := range Keywords {
		if identity == keyword {
			return KEYWORD
		}
	}
	return IDENTITY
}

func (s *Scanner) scanIdentity() (tok Token, lit string) {
	for unicode.IsLetter(s.ch) || s.ch == '_' {
		// 目前仅十进制
		lit += string(s.ch)
		s.next()
	}
	tok = s.isKeyword(lit)
	return
}

// 扫描当前字符
func (s *Scanner) scanNumber() (tok Token, lit string) {
	tok = INTEGER
	for unicode.IsNumber(s.ch) {
		// 目前仅十进制
		lit += string(s.ch)
		s.next()
	}
	return
}

// ScanNext 扫描当前字符返回对应的 Token, 并且偏移 offset 至下一个字符
func (s *Scanner) ScanNext() (tok Token, lit string) {

	s.skipSpace()

	lit = string(s.ch)
	switch s.ch {
	case eof:
		tok = EOF
	case '\n':
		tok = LINEBREAK
		lit = `\n`
	case '+':
		tok = PLUS
	case '-':
		tok = MINUS
	case '*':
		tok = STAR
	case '/':
		tok = SLASH
	case ';':
		tok = SEMICOLON
	case '(':
		tok = LPAREN
	case ')':
		tok = RPAREN
	case '=':
		tok = ASSIGN
	default:
		if unicode.IsNumber(s.ch) {
			// 如果是数字
			tok, lit = s.scanNumber()
			return
		} else if unicode.IsLetter(s.ch) || s.ch == '_' {
			// 如果是字母 (变量名 or 关键词)
			tok, lit = s.scanIdentity()
			return
		}
	}
	s.next()
	return
}