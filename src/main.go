package main

import (
	"my-compiler/ast"
	"my-compiler/printer"
	"my-compiler/token"
	"os"
	"path/filepath"
)

// token -> stmt(token[]) -> function(stmt[]) -> module(function[]) -> package(module[]) -> project(module[])

func main() {

	s := token.NewScanner(filepath.Join(os.Getenv("GOPATH"), "sample", "test1.m"))
	p := ast.NewParser(s)

	// printer.PrintTokens(s)
	// printer.PrintStmts(0, p.Stmts)
	printer.PrintObjectList(0, p.Objects)

}
