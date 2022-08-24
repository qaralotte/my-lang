package rt

import (
	"fmt"
	"my-compiler/ast"
	"my-compiler/token"
	"reflect"
	"strconv"
	"unsafe"
)

type Exec struct {
	Parser *ast.Parser
}

func NewExec(parser *ast.Parser) *Exec {
	return &Exec{
		Parser: parser,
	}
}

func (e *Exec) Run() {
	for {
		stmt, isFinish := e.Parser.ParseStmt()
		if isFinish {
			break
		}
		if stmt != nil {
			e.stmt(stmt)
		}
	}
}

func (e *Exec) RunWithReturn() interface{} {
	for {
		stmt, isFinish := e.Parser.ParseStmt()
		if isFinish {
			break
		}
		if stmt != nil {
			if reflect.TypeOf(stmt).String() == "*ast.ReturnStmt" {
				// 返回语句，立即返回并返回值
				stmt := stmt.(*ast.ReturnStmt)
				return e.expr(stmt.Expr)
			}
			e.stmt(stmt)
		}
	}
	return nil
}

func (e *Exec) stmt(stmt ast.Stmt) {
	switch stmt.(type) {
	case *ast.ExprStmt:
		// 仅表达式的语句
		stmt := stmt.(*ast.ExprStmt)
		e.expr(stmt.Expr)
	case *ast.AssignStmt:
		// 赋值语句
		stmt := stmt.(*ast.AssignStmt)
		e.Parser.Objects.Add(&ast.Variable{
			Name:  stmt.Name,
			Value: e.expr(stmt.Value),
		})
	case *ast.PrintStmt:
		// 打印语句
		stmt := stmt.(*ast.PrintStmt)
		fmt.Println(e.expr(stmt.Expr))
	}
}

func (e *Exec) expr(expr ast.Expr) interface{} {
	switch expr.(type) {
	case *ast.BinaryExpr:
		expr := expr.(*ast.BinaryExpr)
		leftType := reflect.TypeOf(e.expr(expr.Left)).String()
		rightType := reflect.TypeOf(e.expr(expr.Right)).String()

		switch expr.Op {
		// 加法表达式
		case ast.ADD:
			// 字符串相加: 'abc' + 'def' = 'abcdef'
			if leftType == rightType && leftType == "string" {
				return e.expr(expr.Left).(string) + e.expr(expr.Right).(string)
			}

			// 整数相加: 1 + 2 = 3
			if leftType == rightType && leftType == "int64" {
				return e.expr(expr.Left).(int64) + e.expr(expr.Right).(int64)
			}

			// 浮点数相加: 1.1 + 2.2 = 3.3
			if leftType == rightType && leftType == "float64" {
				return e.expr(expr.Left).(float64) + e.expr(expr.Right).(float64)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s + %s", leftType, rightType))
		case ast.SUB:
			// 整数相减: 1 - 2 = -1
			if leftType == rightType && leftType == "int64" {
				return e.expr(expr.Left).(int64) - e.expr(expr.Right).(int64)
			}

			// 浮点数相减: 1.1 - 2.2 = -1.1
			if leftType == rightType && leftType == "float64" {
				return e.expr(expr.Left).(float64) - e.expr(expr.Right).(float64)
			}
		case ast.MUL:

			// 字符串乘整数: 'str' * 3 = 'strstrstr'
			if leftType == "string" && rightType == "int64" {
				result := ""
				right := e.expr(expr.Right).(int64)
				length := *(*int)(unsafe.Pointer(&right))
				for i := 0; i < length; i++ {
					result += e.expr(expr.Left).(string)
				}
				return result
			}

			// 整数相乘: 1 * 2 = 2
			if leftType == rightType && leftType == "int64" {
				return e.expr(expr.Left).(int64) * e.expr(expr.Right).(int64)
			}

			// 浮点数相乘: 1.1 * 2.2 = 2.42
			if leftType == rightType && leftType == "float64" {
				return e.expr(expr.Left).(float64) * e.expr(expr.Right).(float64)
			}
		case ast.DIV:
			// 整数相减: 1 / 2 = 0.5
			if leftType == rightType && leftType == "int64" {
				return e.expr(expr.Left).(int64) / e.expr(expr.Right).(int64)
			}

			// 浮点数相减: 1.1 / 2.2 = 0.5
			if leftType == rightType && leftType == "float64" {
				return e.expr(expr.Left).(float64) / e.expr(expr.Right).(float64)
			}
		}
	case *ast.LitExpr:
		expr := expr.(*ast.LitExpr)
		switch expr.Type {
		case ast.INT:
			// 整数字面量
			val, _ := strconv.ParseInt(expr.Lit, 10, 64)
			return val
		case ast.FLOAT:
			// 浮点数字面量
			val, _ := strconv.ParseFloat(expr.Lit, 64)
			return val
		case ast.STRING, ast.BOOL:
			// 字符串字面量 布尔值字面量
			return expr.Lit
		}
	case *ast.IdentityExpr:
		// 变量
		expr := expr.(*ast.IdentityExpr)
		return expr.Object.(*ast.Variable).Value
	case *ast.CallFnExpr:
		// 方法调用
		expr := expr.(*ast.CallFnExpr)
		return e.callFn(expr.Fn, expr.Params)
	}
	return nil
}

func (e *Exec) callFn(fn *ast.Function, params []ast.Expr) (value interface{}) {
	oldParser := e.Parser.Copy()

	// parser 定位到方法语句块内
	e.Parser.Load(fn.Parser)

	// 加载局部变量表
	e.Parser.Objects = fn.Objects

	// 将具体的表达式传入具体的参数上
	for i := 0; i < len(params); i++ {
		e.Parser.Objects.Add(&ast.Variable{
			Name:  fn.Args[i],
			Value: e.expr(params[i]),
		})
	}

	// 开始语法分析
	e.Parser.EndTokens.Push(token.RBRACE)
	value = e.RunWithReturn()

	// 清空局部变量表并返回到全局变量表
	e.Parser.Objects.Clear()
	e.Parser.Objects = e.Parser.Objects.Get(0).(*ast.Channel).Next

	e.Parser.Load(oldParser)

	return
}
