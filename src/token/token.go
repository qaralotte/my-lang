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

	IDENTITY // abc
	KEYWORD  // _
	INTEGER  // 123
	FLOAT    // 123.456
	STRING   // "xx", 'xx'
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

	IDENTITY: "IDENTITY",
	KEYWORD:  "KEYWORD",
	INTEGER:  "INTEGER",
	FLOAT:    "FLOAT",
	STRING:   "STRING",
}

func String(token Token) string {
	return fmt.Sprintf("token(%s)", tokens[token])
}

var Keywords = []string{
	"_",
}
