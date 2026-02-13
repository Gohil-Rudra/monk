package repl

import (
	"bufio"
	"fmt"
	"io"
	"monk/evaluator"
	"monk/lexer_2"
	"monk/object"
	"monk/parser"
	"strings"
)

const BUDDHA_FACE = `

        _=_
      q(-_-)p
      '_) (_'
      /__/  \
    _(<_   / )_
   (__\_\_|_/__)

   Version 1 : 
   
   Birla Vishvakarma Mahavidyalaya
   IT batch 2023-24( Sem 6 Mini Project )

   Team : Rudra,Jagdish,Yug | Advisor : Dr Nilesh
   
   Feel free to type in commands !

`

func Start(in io.Reader, out io.Writer) {
	const PROMPT = ">>"
	const CONT = "..  "
	scanner := bufio.NewScanner(in)
	io.WriteString(out, BUDDHA_FACE)
	env := object.NewEnvironment()

	line := ""
	braceCount := 0

	for {

		if braceCount == 0 {
			fmt.Fprint(out, PROMPT)

		} else {
			fmt.Fprint(out, CONT)
		}

		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		line += input

		braceCount += strings.Count(input, "{")
		braceCount -= strings.Count(input, "}")

		if braceCount > 0 {
			continue
		}
		l := lexer_2.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			braceCount = 0
			line = ""
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		braceCount = 0
		line = ""

	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err+"\n")
	}
}
