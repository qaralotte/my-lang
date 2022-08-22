package main

import (
	"my-compiler/ast"
	"my-compiler/printer"
	"my-compiler/token"
	"os"
	"path/filepath"
)

func main() {

	s := token.NewScanner(filepath.Join(os.Getenv("GOPATH"), "sample", "test3.m"))
	p := ast.NewParser(s)

	// printer.PrintTokens(s)
	// printer.PrintStmts(0, p.Stmts)
	printer.PrintObjectList(0, p.Objects)

	// e := runtime.NewExec(p)
	// e.Run()

}
