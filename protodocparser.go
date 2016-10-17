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

func isStartComment(line string) bool {
	// TODO
	return true
}

func isEndComment(line string) bool {
	// TODO
	return true
}

func isRpc(line string) bool {
	// TODO
	return true
}

func isService(line string) bool {
	// TODO
	return true
}

func addRpcToLastService(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string) {
	// TODO
}

func addServiceToServices(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string) []*impl.Service {
	// TODO
	return &impl.Service{}
}

func parse(lines []string) string {
	fmt.Printf("Got %v lines of txt\n", len(lines))
	// Create an array of services.
	services := make([]*impl.Service, 0)

	var currentBlock *impl.CommentBlock

	for ln, line := range lines {
		if isStartComment(line) && currentBlock == nil {
			// Create a new comment block.
			currentBlock = &impl.CommentBlock{}
			currentBlock.Start = ln
			currentBlock.Type = impl.OtherComment
		} else if isEndComment(line) && currentBlock != nil {
			currentBlock.End = ln
		} else if isRpc(line) && currentBlock != nil && currentBlock.End > 0 {
			// Mark block as an RPC type.
			currentBlock.Type = impl.RpcComment
			addRpcToLastService(services, currentBlock, lines)
			currentBlock = nil
		} else if isService(line) && currentBlock != nil && currentBlock.End > 0 {
			// Mark block as a Service type.
			currentBlock.Type = impl.ServiceComment
			services = addServiceToServices(services, currentBlock, lines)
			currentBlock = nil
		}
	}

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