package repl

import (
	"os"
)

/*
Interfaces :

	type Printer interface{
		Print()
	}

	now anything with Print method will satisfy this interface...

	type Person struct{
		Name string
	}

	Now to implement the interface we did :

	func (p Person) Print(){
		fmt.Println(p.Name)
	}

	func PrintPerson(p Printer){
		p.Print()
	} -> now this can be used with any type which implements the Print interface

	p := Person{"John Doe"}
	PrintPerson(p)
*/
type Reader interface {
	Read(p []byte) (n int, err error) // here p is a bucket(slice struct or more technically slice header) where we will fill our bytes subsequently we will send output : (n,err) where n is how many bytes written at p and err if any error occured.
}

type Writer interface {
	Write(p []byte) (n int, err error) // p is where to read from for writing
}

type ConsoleReader struct{} // this will implement th struct !
type ConsoleWriter struct{}

func (ConsoleReader) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
} // we implemented the interface method ! now we can use it via Reader reference and its a win as we will be writing a single Echo method which will work with this Reader interface ref and thus thats the abstraction win 1 method for all who wants to Read !

func (ConsoleWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// Now Read input and write it back (Echo)

func Echo(r Reader, w Writer) {
	buf := make([]byte, 1024) // buffer where we write the readed data...
	for {
		n, err := r.Read(buf)
		if err != nil {
			break // this is reading error
		}
		w.Write(buf[:n])
	}
}

//These interfaces are powerful as they abstract the source (for reading) and destination (for writing) of your data. This means that you can read from or write to a variety of sources and destinations (like files, network connections, buffers, etc.) using the same interface.
//
