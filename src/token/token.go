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
	NOT       // !
	EQ        // ==
	NQ        // !=
	GT        // >
	GE        // >=
	LT        // <
	LE        // <=

	IDENTITY  // abc
	INTLIT    // 123
	FLOATLIT  // 123.456
	STRINGLIT // "xx", 'xx'

	TRUE  // true
	FALSE // false
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
	DOT:       ".",
	ASSIGN:    "=",
	EQ:        "==",
	NOT:       "!",
	NQ:        "!=",
	GT:        ">",
	GE:        ">=",
	LT:        "<",
	LE:        "<=",

	IDENTITY:  "IDENTITY",
	INTLIT:    "INTLIT",
	FLOATLIT:  "FLOATLIT",
	STRINGLIT: "STRINGLIT",

	TRUE:  "true",
	FALSE: "false",
}

func String(token Token) string {
	return fmt.Sprintf("token(%s)", tokens[token])
}

type KeywordPair struct {
	Name string
	Token
}

var Keywords = []KeywordPair{
	{"_", EOF},
	{"true", TRUE},
	{"false", FALSE},
}
