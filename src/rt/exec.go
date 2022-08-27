package rt

import (
	"fmt"
	"my-compiler/ast"
	"my-compiler/data"
	"strconv"
)

type Exec struct {
	Parser        *ast.Parser
	ReturnObjects *data.Stack
}

func NewExec(parser *ast.Parser) *Exec {
	return &Exec{
		Parser: parser,
	}
}

func (e *Exec) Run() interface{} {
	for {
		result := e.Parser.ParseStmt()

		if result.IsEndBlock {
			// 如果是endBlock，应该直接结束
			break
		}

		stmt := result.Stmt
		if stmt != nil {
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

		lval, rval := e.expr(expr.Left), e.expr(expr.Right)
		ltype, rtype := ast.ProcessType(ast.GetType(lval), ast.GetType(rval))

		switch expr.Op {
		// 加法表达式
		case ast.ADD:
			// 字符串相加: 'abc' + 'def' = 'abcdef'
			if ast.SameType(ltype, rtype, ast.STRING) {
				return lval.(string) + rval.(string)
			}

			// 整数相加: 1 + 2 = 3
			if ast.SameType(ltype, rtype, ast.INT) {
				return lval.(int64) + rval.(int64)
			}

			// 浮点数相加: 1.1 + 2.2 = 3.3
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) + rval.(float64)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s + %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.SUB:
			// 整数相减: 1 - 2 = -1
			if ast.SameType(ltype, rtype, ast.INT) {
				return lval.(int64) - rval.(int64)
			}

			// 浮点数相减: 1.1 - 2.2 = -1.1
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) - rval.(float64)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s - %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.MUL:

			// 字符串乘整数: 'str' * 3 = 'strstrstr'
			if ltype == ast.INT && rtype == ast.STRING {
				lval, rval = ast.SortByType(lval, rval)

				result := ""
				for i := 0; i < ast.Int64ToInt(lval.(int64)); i++ {
					result += rval.(string)
				}
				return result
			}

			// 整数相乘: 1 * 2 = 2
			if ast.SameType(ltype, rtype, ast.INT) {
				return lval.(int64) * rval.(int64)
			}

			// 浮点数相乘: 1.1 * 2.2 = 2.42
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) * rval.(float64)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s * %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.DIV:
			// 整数相除: 1 / 2 = 0.5
			if ast.SameType(ltype, rtype, ast.INT) {
				return float64(lval.(int64)) / float64(rval.(int64))
			}

			// 浮点数相除: 1.1 / 2.2 = 0.5
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) / rval.(float64)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s / %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.EQ:
			// 1 == 2
			if ast.SameType(ltype, rtype, ast.INT) ||
				ast.SameType(ltype, rtype, ast.FLOAT) ||
				ast.SameType(ltype, rtype, ast.STRING) ||
				ast.SameType(ltype, rtype, ast.BOOL) {
				return lval == rval
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s == %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.NQ:
			// 1 != 2
			if ast.SameType(ltype, rtype, ast.INT) ||
				ast.SameType(ltype, rtype, ast.FLOAT) ||
				ast.SameType(ltype, rtype, ast.STRING) ||
				ast.SameType(ltype, rtype, ast.BOOL) {
				return lval != rval
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s != %s", ast.TypeString(ltype), ast.TypeString(rtype)))

		case ast.GT:
			// 1 > 2
			if ast.SameType(ltype, rtype, ast.INT) {
				return lval.(int64) > rval.(int64)
			}

			// 1.2 > 2.1
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) > rval.(float64)
			}

			// 1.2 > 2.1
			if ast.SameType(ltype, rtype, ast.STRING) {
				return lval.(string) > rval.(string)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s > %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.GE:
			// 1 >= 2
			if ast.SameType(ltype, rtype, ast.INT) {
				return lval.(int64) >= rval.(int64)
			}

			// 1.2 >= 2.1
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) >= rval.(float64)
			}

			// 1.2 >= 2.1
			if ast.SameType(ltype, rtype, ast.STRING) {
				return lval.(string) >= rval.(string)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s >= %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.LT:
			// 1 < 2
			if ast.SameType(ltype, rtype, ast.INT) {
				return lval.(int64) < rval.(int64)
			}

			// 1.2 < 2.1
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) < rval.(float64)
			}

			// 1.2 < 2.1
			if ast.SameType(ltype, rtype, ast.STRING) {
				return lval.(string) < rval.(string)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s < %s", ast.TypeString(ltype), ast.TypeString(rtype)))
		case ast.LE:
			// 1 <= 2
			if ast.SameType(ltype, rtype, ast.INT) {
				return lval.(int64) <= rval.(int64)
			}

			// 1.2 < 2.1
			if ast.SameType(ltype, rtype, ast.FLOAT) {
				return lval.(float64) <= rval.(float64)
			}

			// 1.2 < 2.1
			if ast.SameType(ltype, rtype, ast.STRING) {
				return lval.(string) <= rval.(string)
			}

			panic(fmt.Sprintf("错误: 不合法的运算 %s <= %s", ast.TypeString(ltype), ast.TypeString(rtype)))
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
		case ast.STRING:
			// 字符串字面量
			return expr.Lit
		case ast.BOOL:
			// 布尔值字面量
			return expr.Lit

		}
	case *ast.IdentityExpr:
		// 变量
		expr := expr.(*ast.IdentityExpr)
		return expr.Object.(*ast.Variable).Value
	case *ast.CallFnExpr:
		// 方法调用
		// expr := expr.(*ast.CallFnExpr)
		// return e.callFn(expr.Fn, expr.Params)

	}
	return nil
}
