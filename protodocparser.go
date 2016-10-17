package main

import (
	"fmt"
	"github.com/maniksurtani/protodocparser/impl"
	"encoding/json"
	"io/ioutil"
	"strings"
)

func main() {
	// Read a proto file from StdIn
	protoContents := readFromStdIn()
	j := parse(protoContents)
	fmt.Println(j)
}

func readFromStdIn() []string {
	//todo read from stdin and convert this into a proper test
	dat, err := ioutil.ReadFile("./sample.proto")
	check(err)
	return strings.Split(string(dat), "\n")
}

func parse(lines []string) string {
	// Create an array of services.
	services := make([]*impl.Service, 0)
	// TODO parse lines, only read comments, add to services array.
	bytes, err := json.Marshal(services)
	if err != nil {
		panic(fmt.Sprintf("Caught error %v trying to serialize %v into JSON.", err, services))
	}
	return string(bytes)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}