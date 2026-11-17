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
