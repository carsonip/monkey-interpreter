package main

import (
	"github.com/carsonip/monkey-interpreter/repl"
	"os"
)

func main() {
	r := repl.Repl{}
	r.Start(os.Stdin, os.Stdout)
}
