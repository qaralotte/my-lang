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
	LBRACE    // {
	RBRACE    // }
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

	TRUE
	FALSE
	FN
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
	LBRACE:    "{",
	RBRACE:    "}",
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
	FN:    "fn",
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
	{"fn", FN},
}
