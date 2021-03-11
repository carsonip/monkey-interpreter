package repl

import (
	"bufio"
	"fmt"
	"github.com/carsonip/monkey-interpreter/eval"
	"github.com/carsonip/monkey-interpreter/object"
	"github.com/carsonip/monkey-interpreter/parser"
	"github.com/carsonip/monkey-interpreter/token"
	"io"
)

type Repl struct {}

const PROMPT = ">> "

func (r *Repl) Start(in io.Reader, out io.Writer) {
	env := object.NewEnv()
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		if ok := scanner.Scan(); !ok {
			break
		}
		if line := scanner.Text(); line != "" {
			lex := token.NewLexer(line)
			p := parser.NewParser(&lex)
			ev := eval.NewEvaluator(&p, env)
			for obj := ev.EvalNext(env); obj != nil; obj = ev.EvalNext(env) {
				fmt.Fprintf(out, "%s\n", obj)
			}
		}
	}
}
