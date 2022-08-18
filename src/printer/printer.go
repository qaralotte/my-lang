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
		fmt.Println("[AssignStmt]")
		PrintObject(indentation+1, stmt.Object)
		PrintExpr(indentation+1, stmt.Right)
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
	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	switch obj.(type) {
	case *ast.Variable:
		o := obj.(*ast.Variable)
		fmt.Printf("[Variable] type: %s, name: %s\n", ast.TypeString(o.Type), o.Name)
	case *ast.Function:
		o := obj.(*ast.Function)
		fmt.Printf("[Function] type: %s, name: %s\n", ast.TypeString(o.Type), o.Name)
		for i := 0; i < len(o.Args); i++ {
			PrintObject(indentation+1, o.Args[i])
		}
		PrintStmts(indentation+1, o.Stmts)
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
