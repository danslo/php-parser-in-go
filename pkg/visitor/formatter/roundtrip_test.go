package formatter_test

import (
	"bytes"
	"testing"

	"github.com/rectorphp/php-parser-in-go/pkg/ast"
	"github.com/rectorphp/php-parser-in-go/pkg/conf"
	"github.com/rectorphp/php-parser-in-go/pkg/parser"
	"github.com/rectorphp/php-parser-in-go/pkg/version"
	"github.com/rectorphp/php-parser-in-go/pkg/visitor/formatter"
	"github.com/rectorphp/php-parser-in-go/pkg/visitor/printer"
)

// format parses src, lays it out from scratch and prints it.
func format(test *testing.T, src string) string {
	test.Helper()
	ver, err := version.New("8.4")
	if err != nil {
		test.Fatal(err)
	}
	root, err := parser.Parse([]byte(src), conf.Config{Version: ver})
	if err != nil {
		test.Fatal(err)
	}
	root.Accept(formatter.NewFormatter())
	var buffer bytes.Buffer
	root.(*ast.Root).Accept(printer.NewPrinter(&buffer))
	return buffer.String()
}

func TestFormatter_KeepsNowdoc(test *testing.T) {
	got := format(test, "<?php $a = <<<'TXT'\nval $x \\t end\nTXT;\n")
	expected := "<?php\n\n$a = <<<'TXT'\nval $x \\t end\nTXT;\n"
	if got != expected {
		test.Errorf("\nexpected: %q\ngot: %q\n", expected, got)
	}
}

func TestFormatter_KeepsHeredocLabel(test *testing.T) {
	got := format(test, "<?php $a = <<<TXT\nEOT; x\nTXT;\n")
	expected := "<?php\n\n$a = <<<TXT\nEOT; x\nTXT;\n"
	if got != expected {
		test.Errorf("\nexpected: %q\ngot: %q\n", expected, got)
	}
}

func TestFormatter_NewlineAfterCloseTag(test *testing.T) {
	got := format(test, "<?php echo 1; ?>\n\nX<?php echo 2;")
	expected := "<?php\n\necho 1;?>\n\nX<?php echo 2;"
	if got != expected {
		test.Errorf("\nexpected: %q\ngot: %q\n", expected, got)
	}
}

// print lays out a tree built without tokens and prints it.
func print(root *ast.Root) string {
	root.Accept(formatter.NewFormatter())
	var buffer bytes.Buffer
	root.Accept(printer.NewPrinter(&buffer))
	return buffer.String()
}

func echo(value string) ast.Vertex {
	return &ast.StmtEcho{Exprs: []ast.Vertex{&ast.ScalarLnumber{Value: []byte(value)}}}
}

func variable(name string) ast.Vertex {
	return &ast.ExprVariable{Name: &ast.Identifier{Value: []byte("$" + name)}}
}

func TestFormatter_BracesDanglingElse(test *testing.T) {
	inner := &ast.StmtIf{Cond: variable("a"), Stmt: echo("1")}
	loop := &ast.StmtWhile{Cond: variable("b"), Stmt: &ast.StmtIf{Cond: variable("c"), Stmt: echo("2")}}
	root := &ast.Root{Stmts: []ast.Vertex{
		&ast.StmtIf{Cond: variable("p"), Stmt: inner, Else: &ast.StmtElse{Stmt: echo("3")}},
		&ast.StmtIf{
			Cond:   variable("q"),
			Stmt:   echo("4"),
			ElseIf: []ast.Vertex{&ast.StmtElseIf{Cond: variable("r"), Stmt: loop}},
			Else:   &ast.StmtElse{Stmt: echo("5")},
		},
	}}
	got := print(root)
	expected := "<?php\n\nif ($p) {\n    if ($a) echo 1;\n} else echo 3;\nif ($q) echo 4; elseif ($r) {\n    while($b) if ($c) echo 2;\n} else echo 5;"
	if got != expected {
		test.Errorf("\nexpected: %q\ngot: %q\n", expected, got)
	}
}
