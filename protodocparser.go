package main

import "fmt"

func main() {
	// Read a proto file from StdIn
	protoContents := readFromStdIn()
	j := parse(protoContents)
	printToStdOut(j)
}

func readFromStdIn() []string {
	// TODO actually read from StdIn, and split on newlines
	return []string{"a", "b", "c"}
}

func parse(lines []string) string {
	// TODO parse lines, only read comments, return serialized JSON
	return "[{}]"
}

func printToStdOut(j string) {
	fmt.Print("TODO: print JSON\n\n")
}
