package main

import (
	"fmt"
	"github.com/maniksurtani/protodocparser/impl"
	"encoding/json"
	"regexp"
	"bufio"
	"os"
)

// Regexps
var startCommentRE = regexp.MustCompile("\\s*/\\*\\*")
var endCommentRE = regexp.MustCompile("\\s*\\*\\\\\\s*$")
var rpcRE = regexp.MustCompile("\\s*rpc\\s+")
var serviceRE = regexp.MustCompile("\\s*service\\s+")

func main() {
	// Read a proto file from StdIn
	protoContents := readFromStdIn()
	j := parse(protoContents)
	fmt.Println(j)
}

func readFromStdIn() []string {
	s := bufio.NewScanner(os.Stdin)
	txt := make([]string, 0)
	for s.Scan() {
		txt = append(txt, s.Text())
	}
	return txt
}

func addRpcToLastService(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string) {
	// TODO
}

func addServiceToServices(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string) []*impl.Service {
	// TODO
	//return &impl.Service{}
	return nil;
}

func parse(lines []string) string {
	fmt.Printf("Got %v lines of txt\n", len(lines))
	// Create an array of services.
	services := make([]*impl.Service, 0)

	var currentBlock *impl.CommentBlock

	for ln, line := range lines {
		if startCommentRE.MatchString(line) && currentBlock == nil {
			// Create a new comment block.
			currentBlock = &impl.CommentBlock{}
			currentBlock.Start = ln
			currentBlock.Type = impl.OtherComment
		} else if endCommentRE.MatchString(line) && currentBlock != nil {
			currentBlock.End = ln
		} else if rpcRE.MatchString(line) && currentBlock != nil && currentBlock.End > 0 {
			// Mark block as an RPC type.
			currentBlock.Type = impl.RpcComment
			addRpcToLastService(services, currentBlock, lines)
			currentBlock = nil
		} else if serviceRE.MatchString(line) && currentBlock != nil && currentBlock.End > 0 {
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
