package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

type instrumenter struct {
	traceImport string
	tracePkg    string
	traceFunc   string
}

func New(traceImport, tracePkg, traceFunc string) *instrumenter {
	return &instrumenter{
		traceImport: traceImport,
		tracePkg:    tracePkg,
		traceFunc:   traceFunc,
	}
}

// 判断是否含有函数定义语句
func hasFuncDecl(f *ast.File) bool {
	if len(f.Decls) == 0 {
		return false
	}

	for _, decl := range f.Decls {
		if _, ok := decl.(*ast.FuncDecl); ok {
			return true
		}
	}
	return false
}

//实现Instrumenter接口

func (a instrumenter) Instrument(filename string) ([]byte, error) {
	fset := token.NewFileSet()
	curAST, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w\n", filename, err)
	}

	if !hasFuncDecl(curAST) {
		return nil, nil
	}
	astutil.AddImport(fset, curAST, a.traceImport)

	a.addDeferTraceIntoFuncDecls(curAST)
	//在这个前面添加build注释
	buf := &bytes.Buffer{}

	err = format.Node(buf, fset, curAST)

	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w\n", err)
	}
	return buf.Bytes(), nil
}

func (a instrumenter) addDeferTraceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			a.addDeferStmt(fd)
		}
	}
}

func (a instrumenter) addDeferStmt(fd *ast.FuncDecl) bool {
	stmts := fd.Body.List

	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt) //是否是defer语句
		if !ok {
			continue
		}
		ce, ok := ds.Call.Fun.(*ast.CallExpr) //判断defer语句是否是函数调用
		if !ok {
			continue
		}

		se, ok := ce.Fun.(*ast.SelectorExpr) //判断这个函数调用是否是a.b()
		if !ok {
			continue
		}
		x, ok := se.X.(*ast.Ident) //取出a
		if !ok {
			continue
		}
		if (x.Name == a.tracePkg) && (se.Sel.Name == a.traceFunc) {
			return false
		}
	}

	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: a.tracePkg,
					},
					Sel: &ast.Ident{
						Name: a.traceFunc,
					},
				},
			},
		},
	}
	//comment := &ast.Comment{
	//	Text: "// +build ignore\n",
	//}
	//label := &ast.LabeledStmt{
	//	Label: ast.NewIdent("_"),
	//	Stmt:  &ast.EmptyStmt{},
	//}
	//comment.Slash = label.Pos()

	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	//newList[0] = label
	newList[0] = ds
	fd.Body.List = newList

	return true
}
