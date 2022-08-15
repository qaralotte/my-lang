package token

import "fmt"

type Token int

const (
	EOF Token = iota

	PLUS      // +
	MINUS     // -
	STAR      // *
	SLASH     // /
	LINEBREAK // \n
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	ASSIGN    // =
	DOT       // .

	IDENTITY  // abc
	KEYWORD   // _
	INTLIT    // 123
	FLOATLIT  // 123.456
	STRINGLIT // "xx", 'xx'
)

var tokens = map[Token]string{
	EOF: "EOF",

	PLUS:      "+",
	MINUS:     "-",
	STAR:      "*",
	SLASH:     "/",
	LINEBREAK: "LINE-BREAK",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	ASSIGN:    "=",
	DOT:       ".",

	IDENTITY:  "IDENTITY",
	KEYWORD:   "KEYWORD",
	INTLIT:    "INTLIT",
	FLOATLIT:  "FLOATLIT",
	STRINGLIT: "STRINGLIT",
}

func String(token Token) string {
	return fmt.Sprintf("token(%s)", tokens[token])
}

var Keywords = []string{
	"_",
}
