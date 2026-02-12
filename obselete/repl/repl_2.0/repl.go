package repl

import (
	"bufio"
	"fmt"
	"io"
	"monk/lexer_2"
	"monk/parser"
)

const BUDDHA_FACE = "\033[32m" + `

        _=_
      q(-_-)p
      '_) (_'
      /__/  \
    _(<_   / )_
   (__\_\_|_/__)

   Feel free to type in commands !

` + "\033[0m"

func Start(in io.Reader, out io.Writer) {
	const PROMPT = ">>"
	scanner := bufio.NewScanner(in)
	io.WriteString(out, BUDDHA_FACE)
	for {
		fmt.Fprint(out, PROMPT)

		if !scanner.Scan() {
			return
		}
		line := scanner.Text()
		l := lexer_2.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err+"\n")
	}
}
