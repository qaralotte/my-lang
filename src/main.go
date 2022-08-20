package main

import (
	"my-compiler/ast"
	"my-compiler/runtime"
	"my-compiler/token"
	"os"
	"path/filepath"
)

func main() {

	s := token.NewScanner(filepath.Join(os.Getenv("GOPATH"), "sample", "test2.m"))
	p := ast.NewParser(s)
	e := runtime.NewExec(p)
	e.Run()

	// printer.PrintTokens(s)
	// printer.PrintStmts(0, p.Stmts)
	// printer.PrintObjectList(0, p.Objects)

}
