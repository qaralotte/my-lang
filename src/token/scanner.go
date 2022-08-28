package token

import (
	"os"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	// todo 文件名
	src []byte // 源码

	offset   int  // 当前偏移位置
	ch       rune // 当前读取的字符 (utf-8)
	nearlyCh byte // 下一个紧挨着的字符 (必定是 ascii)
}

// end of file
const eof = -1

func NewScanner(path string) (s *Scanner) {
	var scanner Scanner

	// 读取文件内容并存进 src
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	scanner.src = bytes

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

		// 下一个紧挨着的字符
		if s.offset <= len(s.src)-1 {
			s.nearlyCh = s.src[s.offset]
		} else {
			s.nearlyCh = 0
		}
	} else {
		s.ch = eof
	}
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

// 判断是否是关键词
func (s *Scanner) isKeyword(identity string) Type {
	for _, keyword := range Keywords {
		if identity == keyword.Name {
			return keyword.Type
		}
	}
	return IDENTITY
}

func (s *Scanner) nextNearlyChar(ch byte) bool {
	if s.nearlyCh == ch {
		s.next()
		return true
	}
	return false
}

// 扫描变量或者关键词字面量
func (s *Scanner) scanIdentity() (tok Token) {
	for unicode.IsLetter(s.ch) || s.ch == '_' {
		// 目前仅十进制
		tok.Lit += string(s.ch)
		s.next()
	}
	tok.Type = s.isKeyword(tok.Lit)
	return
}

func (s *Scanner) scanFloat(decimal string) (tok Token) {
	tok.Type = FLOATLIT
	for unicode.IsNumber(s.ch) {
		// 目前仅十进制
		tok.Lit += string(s.ch)
		s.next()
	}
	tok.Lit = decimal + "." + tok.Lit
	return
}

// 扫描数字字面量
func (s *Scanner) scanNumber() (tok Token) {
	tok.Type = INTLIT
	for unicode.IsNumber(s.ch) {
		// 目前仅十进制
		tok.Lit += string(s.ch)
		s.next()

		// 如果有小数点，则应该是浮点数
		if s.ch == '.' {
			s.next()
			return s.scanFloat(tok.Lit)
		}

	}
	return
}

// 扫描字符串字面量
func (s *Scanner) scanString(end rune) (tok Token) {
	tok.Type = STRINGLIT
	// [']xxx'
	s.next()

	// '[xxx]'
	for s.ch != end {
		tok.Lit += string(s.ch)
		s.next()
	}

	// 'xxx[']
	s.next()

	return
}

// ScanNext 扫描当前字符返回对应的 Token, 并且偏移 offset 至下一个字符
func (s *Scanner) scanNext() (tok Token) {

	s.skipSpace()

	switch s.ch {
	case eof:
		tok.Type = EOF
	case '\n':
		tok.Type = LINEBREAK
		tok.Lit = `\n`
	case '+':
		tok.Type = PLUS
	case '-':
		tok.Type = MINUS
	case '*':
		tok.Type = STAR
	case '/':
		tok.Type = SLASH
	case ';':
		tok.Type = SEMICOLON
	case '(':
		tok.Type = LPAREN
	case ')':
		tok.Type = RPAREN
	case '{':
		tok.Type = LBRACE
	case '}':
		tok.Type = RBRACE
	case '.':
		tok.Type = DOT
	case ',':
		tok.Type = COMMA
	case '\'':
		tok = s.scanString('\'')
		return
	case '"':
		tok = s.scanString('"')
		return
	case '=':
		tok.Type = ASSIGN
		if s.nextNearlyChar('=') {
			tok.Type = EQ
		}
	case '!':
		tok.Type = NOT
		if s.nextNearlyChar('=') {
			tok.Type = NQ
		}
	case '>':
		tok.Type = GT
		if s.nextNearlyChar('=') {
			tok.Type = GE
		}
	case '<':
		tok.Type = LT
		if s.nextNearlyChar('=') {
			tok.Type = LE
		}
	default:
		if unicode.IsNumber(s.ch) {
			// 如果是数字
			tok = s.scanNumber()
			return
		} else if unicode.IsLetter(s.ch) || s.ch == '_' {
			// 如果是字母 (变量名 or 关键词)
			tok = s.scanIdentity()
			return
		}
	}
	s.next()
	return
}

func (s *Scanner) ScanTokens() (toks []Token) {
	for tok := s.scanNext(); tok.Type != EOF; tok = s.scanNext() {
		toks = append(toks, tok)
	}
	return
}
