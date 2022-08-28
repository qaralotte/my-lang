package main

import (
	"my-compiler/ast"
	"my-compiler/rt"
	"my-compiler/token"
	"os"
	"path/filepath"
)

func main() {

	scanner := token.NewScanner(filepath.Join(os.Getenv("GOPATH"), "sample", "test1.m"))
	toks := scanner.ScanTokens()
	// printer.PrintTokens(toks)

	p := ast.NewParser(toks)
	e := rt.NewExec(p)
	e.Run()

}
