package main

import (
	"fmt"
	"github.com/maniksurtani/protodocparser/impl"
	"encoding/json"
	"regexp"
	"bufio"
	"os"
	"strings"
)

// Regexps
var startCommentRE = regexp.MustCompile("^\\s*/\\*")
var endCommentRE = regexp.MustCompile("\\s*\\*/\\s*$")
var rpcRE = regexp.MustCompile("\\s*rpc\\s+")
var serviceRE = regexp.MustCompile("\\s*service\\s+")
var serviceNameRE = regexp.MustCompile("\\s*service\\s+(\\w+)\\s*\\{")
var rpcNameRE = regexp.MustCompile("\\s*rpc\\s+(\\w+)\\s*\\(\\s*([\\w]+)\\s*\\)\\s+returns\\s+\\(\\s*(\\w+)\\s*\\)")
var pkgNameRE = regexp.MustCompile("\\s*package\\s+([\\w.]+)\\s*;")

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

func addRpcToLastService(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string, currentLine int) {
	rpc := impl.NewRpc()
	rpc.Name, rpc.Request, rpc.Response = rpcName(lines[currentLine])

	// TODO add other sections of the rpc.
	lastService := services[len(services) - 1]
	lastService.Rpcs = append(lastService.Rpcs, rpc)
}

func addServiceToServices(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string, currentLine int) []*impl.Service {
	s := impl.NewService()
	s.Name = serviceName(lines[currentLine])

	// TODO add other sections of the service.
	return append(services, s)
}

func serviceName(line string) string {
	r := serviceNameRE.FindStringSubmatch(line)
	if len(r) > 0 {
		return r[1]
	} else {
		return ""
	}
}

func rpcName(line string) (name, req, rsp string) {
	r := rpcNameRE.FindStringSubmatch(line)
	lenR := len(r)
	if lenR > 0 {
		name = r[1]
	}

	if lenR > 1 {
		req = r[2]
	}

	if lenR > 2 {
		rsp = r[3]
	}

	return
}

func matchPkgName(line string) (string, bool) {
	matches := pkgNameRE.FindStringSubmatch(line)
	if len(matches) > 0 {
		return matches[1], true
	}

	return "", false
}

func parseString(fileAsString string) string {
	return parse(strings.Split(fileAsString, "\n"))
}

func parse(lines []string) string {
	//fmt.Printf("Got %v lines of txt\n", len(lines))
	// Create an array of services.
	services := make([]*impl.Service, 0)

	var currentBlock *impl.CommentBlock
	pkgName := ""
	matched := false

	for ln, line := range lines {
		if pkgName == "" {
			pkgName, matched = matchPkgName(line)
			if matched {
				continue
			}
		}
		//fmt.Printf("Line %v is [%v]\n", ln, line)
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
			addRpcToLastService(services, currentBlock, lines, ln)
			currentBlock = nil
		} else if serviceRE.MatchString(line) && currentBlock != nil && currentBlock.End > 0 {
			// Mark block as a Service type.
			currentBlock.Type = impl.ServiceComment
			services = addServiceToServices(services, currentBlock, lines, ln)
			currentBlock = nil
		}
	}

	// Add package names to services.
	for _, svc := range services {
		svc.Package = pkgName
	}

	bytes, err := json.Marshal(services)
	if err != nil {
		panic(fmt.Sprintf("Caught error %v trying to serialize %v into JSON.", err, services))
	}
	return string(bytes)
}
