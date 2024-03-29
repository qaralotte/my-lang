package main

import (
	"my-lang/ast"
	"my-lang/rt"
	"my-lang/token"
	"os"
)

type MyArgs struct {
	mainFile string
}

func initArgs() (args MyArgs) {

	if len(os.Args) > 1 {
		args.mainFile = os.Args[1]
	}

	return
}

func main() {

	args := initArgs()

	// 新建扫描器
	scanner := token.NewScanner(args.mainFile)

	// 扫描所有的 tokens
	toks := scanner.ScanTokens()

	// 调试 tokens 结果
	// token.Debug(toks)

	// 全局对象表
	globalObjs := ast.NewObjectList(nil)

	// 新建解析器
	p := ast.NewParser(toks, globalObjs)

	// 新建解释器
	e := rt.NewExec(p)

	// 运行
	e.Run()

}
