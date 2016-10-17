package main

import (
	"fmt"
	"github.com/maniksurtani/protodocparser/impl"
	"encoding/json"
)

func main() {
	// Read a proto file from StdIn
	protoContents := readFromStdIn()
	j := parse(protoContents)
	fmt.Println(j)
}

func readFromStdIn() []string {
	// TODO actually read from StdIn, and split on newlines
	return []string{"a", "b", "c"}
}

func parse(lines []string) []byte {
	// Create an array of services.
	services := make([]*impl.Service, 0)
	// TODO parse lines, only read comments, add to services array.
	bytes, err := json.Marshal(services)
	if err != nil {
		panic(fmt.Sprintf("Caught error %v trying to serialize %v into JSON.", err, services))
	}
	return bytes
}
