package ast

import (
	"fmt"
	"my-compiler/object"
	"my-compiler/token"
)

func defFn() {
	// [fn] name(...) {...}
	globalParser.Require(token.FN, true)

	// fn [name](...) {...}
	name := globalParser.Require(token.IDENTITY, true)
	obj := object.FindObject(globalParser.Objects, name)
	if obj != nil {
		panic(fmt.Sprintf("错误: 重复定义的名称 %s", name))
	}

	// fn name[(...)] {...}
	var args []object.Object = nil
	if globalParser.Token == token.LPAREN {
		// todo 函数参数定义
		globalParser.NextToken()
		globalParser.Require(token.RPAREN, true)
	}

	// 添加方法进对象表
	fn := object.NewFunction(name, args, globalParser.Objects)
	globalParser.Objects = append(globalParser.Objects, fn)

	// fn name(...) [{]...}
	globalParser.Require(token.LBRACE, true)
	// 穿梭进局部对象表
	globalParser.ShuttleChannel(fn.Next)
	// todo GET AST
	globalParser.ParseStmts(token.RBRACE)
	// 穿梭出全局对象表
	globalParser.ShuttleChannel(globalParser.Objects[0].(*object.Channel))
	// fn name(...) {...[}]
	globalParser.Require(token.RBRACE, true)
}
