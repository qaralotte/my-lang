package main

import (
	"my-compiler/ast"
	"my-compiler/rt"
	"my-compiler/token"
	"os"
	"path/filepath"
)

func main() {

	s := token.NewScanner(filepath.Join(os.Getenv("GOPATH"), "sample", "test1.m"))
	// printer.PrintTokens(s)

	p := ast.NewParser(s)
	e := rt.NewExec(p)
	e.Run()

}
