package main

import (
	"my-compiler/ast"
	"my-compiler/rt"
	"my-compiler/token"
	"os"
	"path/filepath"
)

func main() {

	// 新建扫描器
	scanner := token.NewScanner(filepath.Join(os.Getenv("GOPATH"), "sample", "test1.m"))

	// 扫描所有的 tokens
	toks := scanner.ScanTokens()

	// 调试 tokens 结果
	// printer.PrintTokens(toks)

	// 全局对象表
	globalObjs := ast.NewObjectList(nil)

	// 新建解析器
	p := ast.NewParser(toks, globalObjs)

	// 新建解释器
	e := rt.NewExec(p)

	// 运行
	e.Run()

}
