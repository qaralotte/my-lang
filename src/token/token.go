package token

import "fmt"

type (
	Type  int
	Token struct {
		Type        // 类型
		Lit  string // 字面量
	}
)

func EmptyToken(p Type) Token {
	return Token{
		Type: p,
		Lit:  "",
	}
}

const (
	EOF Type = iota

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
	DOT       // .
	COMMA     // ,
	ASSIGN    // =
	EQ        // ==
	NOT       // !
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
	RETURN
	PRINT
	IF
	ELSE
	FOR
)

var tokens = map[Type]string{
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
	COMMA:     ",",
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

	TRUE:   "true",
	FALSE:  "false",
	RETURN: "return",
	PRINT:  "print",
	IF:     "if",
	ELSE:   "else",
	FOR:    "for",
}

func TypeString(tokType Type) string {
	return fmt.Sprintf("token(%s)", tokens[tokType])
}

type KeywordPair struct {
	Name string
	Type
}

var Keywords = []KeywordPair{
	{"_", EOF},
	{"true", TRUE},
	{"false", FALSE},
	{"return", RETURN},
	{"print", PRINT},
	{"if", IF},
	{"else", ELSE},
	{"for", FOR},
}
