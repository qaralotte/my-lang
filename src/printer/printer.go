package printer

import (
	"fmt"
	"my-compiler/ast"
	"my-compiler/scanner"
	"my-compiler/token"
	"my-compiler/variable"
)

func PrintTokens(s *scanner.Scanner) {
	for tok, lit := s.ScanNext(); tok != token.EOF; tok, lit = s.ScanNext() {
		fmt.Printf("%s: %s\n", token.String(tok), lit)
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
		PrintVariable(indentation+1, expr.Var)
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
		PrintVariable(indentation+1, stmt.Var)
		PrintExpr(indentation+1, stmt.Right)
	}
}

func PrintVariable(indentation int, va *variable.Variable) {
	for i := 0; i < indentation; i++ {
		fmt.Print("  ")
	}

	fmt.Printf("[Var] type: %s, name: %s\n", variable.TypeString(va.Type), va.Name)
}
