package repl

import (
	"bufio"
	"fmt"
	"github.com/carsonip/monkey-interpreter/parser"
	"github.com/carsonip/monkey-interpreter/token"
	"io"
)

type Repl struct {}

const PROMPT = ">> "

func (r *Repl) Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		if ok := scanner.Scan(); !ok {
			break
		}
		if line := scanner.Text(); line != "" {
			lex := token.NewLexer(line)
			p := parser.NewParser(&lex)
			for node := p.NextNode(); node != nil; node = p.NextNode(){
				fmt.Fprintf(out, "%v\n", node)
			}
		}
	}
}
