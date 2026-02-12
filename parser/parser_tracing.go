package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel)
}

func tracePrint(ms string) {
	fmt.Printf("%s%s\n", identLevel(), ms)
}

func incIdent() {
	traceLevel += 1
}

func decIdent() {
	traceLevel -= 1
}

func trace(msg string) string {
	tracePrint("Begin " + msg)
	incIdent()
	return msg
}

func untrace(msg string) {
	decIdent()
	tracePrint("End " + msg)

}
