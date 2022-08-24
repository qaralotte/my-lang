package printer

import (
	"fmt"
	"my-compiler/ast"
	"my-compiler/token"
)

func PrintTokens(s *token.Scanner) {
	for tok, lit := s.ScanNext(); tok != token.EOF; tok, lit = s.ScanNext() {
		fmt.Print(token.String(tok))
		if lit != "" {
			fmt.Printf(": %s", lit)
		}
		fmt.Println()
	}
}

func PrintExpr(indentation int, expr ast.Expr) {
	if expr == nil {
		return
	}

	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	switch expr.(type) {
	case *ast.LitExpr:
		// 字面量
		expr := expr.(*ast.LitExpr)
		fmt.Printf("[LitExpr] Lit: %s\n", expr.Lit)
	case *ast.IdentityExpr:
		// 变量
		expr := expr.(*ast.IdentityExpr)
		fmt.Println("[IdentityExpr]")
		PrintObject(indentation+1, expr.Object)
	case *ast.BinaryExpr:
		// Left [Op] Right
		expr := expr.(*ast.BinaryExpr)
		fmt.Printf("[BinaryExpr] Op: %s\n", ast.OperatorString(expr.Op))
		PrintExpr(indentation+1, expr.Left)
		PrintExpr(indentation+1, expr.Right)
	case *ast.CallFnExpr:
		// 调用方法
		expr := expr.(*ast.CallFnExpr)
		fmt.Println("[CallFnExpr]")
		PrintObject(indentation+1, expr.Fn)
		for i := 0; i < len(expr.Params); i++ {
			PrintExpr(indentation+1, expr.Params[i])
		}

	}
}

func PrintStmt(indentation int, stmt ast.Stmt) {

	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	switch stmt.(type) {
	// 纯表达式
	case *ast.ExprStmt:
		fmt.Println("[ExprStmt]")
		PrintExpr(indentation+1, stmt.(*ast.ExprStmt).Expr)
	case *ast.AssignStmt:
		stmt := stmt.(*ast.AssignStmt)
		fmt.Printf("[AssignStmt] name: %s\n", stmt.Name)
		PrintExpr(indentation+1, stmt.Value)
	case *ast.PrintStmt:
		stmt := stmt.(*ast.PrintStmt)
		fmt.Println("[PrintStmt]")
		PrintExpr(indentation+1, stmt.Expr)
	case *ast.ReturnStmt:
		stmt := stmt.(*ast.ReturnStmt)
		fmt.Println("[ReturnStmt]")
		PrintExpr(indentation+1, stmt.Expr)
	}
}

func PrintStmts(indentation int, stmts []ast.Stmt) {
	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	fmt.Println("[Stmts]")
	for i := 0; i < len(stmts); i++ {
		PrintStmt(indentation+1, stmts[i])
	}
}

func PrintObject(indentation int, obj ast.Object) {
	if obj == nil {
		return
	}

	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	switch obj.(type) {
	case *ast.Variable:
		o := obj.(*ast.Variable)
		fmt.Printf("[Variable] name: %s, value: %s\n", o.Name, o.Value)
	case *ast.Function:
		o := obj.(*ast.Function)
		fmt.Printf("[Function] name: %s, args: %v\n", o.Name, o.Args)
	}
}

func PrintObjectList(indentation int, objs *ast.ObjectList) {
	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	fmt.Println("[Objects]")
	for i := 1; i < len(*objs.Objects); i++ {
		PrintObject(indentation+1, (*objs.Objects)[i])
	}
}
