package printer

import (
	"fmt"
	"my-compiler/ast"
	"my-compiler/object"
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

func PrintObject(indentation int, obj object.Object) {
	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	switch obj.(type) {
	case *object.Variable:
		o := obj.(*object.Variable)
		fmt.Printf("[Variable] type: %s, name: %s\n", object.TypeString(o.Type), o.Name)
	case *object.Function:
		o := obj.(*object.Function)
		fmt.Printf("[Function] type: %s, name: %s\n", object.TypeString(o.Type), o.Name)
		for i := 0; i < len(o.Args); i++ {
			PrintObject(indentation+1, o.Args[i])
		}
	}
}
