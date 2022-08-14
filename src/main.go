package main

import (
	"my-compiler/ast"
	"my-compiler/printer"
	"my-compiler/scanner"
	"os"
	"path/filepath"
)

func main() {

	s := scanner.NewScanner(filepath.Join(os.Getenv("GOPATH"), "test.m"))
	// printer.PrintTokens(s)

	p := ast.NewParser(s)
	for _, stmt := range p.Stmts {
		printer.PrintStmt(0, stmt)
	}

}
