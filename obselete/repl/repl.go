package repl

import (
	"bufio"
	"fmt"
	"io"
	"monk/lexer_2"
	"monk/token"
)

func Start(in io.Reader, out io.Writer) {
	const PROMPT = ">>"
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)

		if !scanner.Scan() {
			return
		}
		line := scanner.Text()
		l := lexer_2.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() { // this is a for loop where init ; condition ; Post is present
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
