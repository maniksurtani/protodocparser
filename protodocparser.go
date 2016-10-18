package main

import (
	"fmt"
	"github.com/maniksurtani/protodocparser/impl"
	"encoding/json"
	"regexp"
	"bufio"
	"os"
	"strings"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Regexps
var (
	startCommentRE = regexp.MustCompile("^\\s*/\\*")
	endCommentRE = regexp.MustCompile("\\s*\\*/\\s*$")
	rpcRE = regexp.MustCompile("\\s*rpc\\s+")
	serviceRE = regexp.MustCompile("\\s*service\\s+")
	serviceNameRE = regexp.MustCompile("\\s*service\\s+(\\w+)\\s*\\{")
	rpcNameRE = regexp.MustCompile("\\s*rpc\\s+(\\w+)\\s*\\(\\s*([\\w]+)\\s*\\)\\s+returns\\s+\\(\\s*(\\w+)\\s*\\)")
	pkgNameRE = regexp.MustCompile("\\s*package\\s+([\\w.]+)\\s*;")
)

// Command-line args
var (
	debug = kingpin.Flag("verbose", "Enable verbose mode").Bool()
	// TODO should we allow multiple proto files? Separated by ':' or something?
	protoFile = kingpin.Arg("protofile", "Proto file").Required().String()
	outFile = kingpin.Arg("out", "Output file").Required().String()
)

func main() {
	// TODO: parse command line parameters. Do do so, use the kingpin package. See https://github.com/alecthomas/kingpin#examples
	kingpin.Version("0.1") // Version of protodocparser
	kingpin.Parse() // Parse command-line options


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

	// TODO: set Options, Doc and Examples
	// TODO: Doc is the comment block before the first @Example annotation
	// TODO: @Examples exists in the comment block
	// TODO: Options are the protobuf options. This might be harder to figure out since they may be split across multiple lines. :/

	lastService := services[len(services) - 1]
	lastService.Rpcs = append(lastService.Rpcs, rpc)
}

func addServiceToServices(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string, currentLine int) []*impl.Service {
	s := impl.NewService()
	s.Name = serviceName(lines[currentLine])

	// TODO: set Api, Design, Doc, Examples, File, Org and Url
	// TODO: to set Api, just check whether the @API annotation exists in the comment block
	// TODO: to set Design and Org, look at the params passed in to @API
	// TODO: if Org isn't set, attempt to "guess" what it might be by looking at the path/package of the proto, and look up in Registry
	// TODO: @Examples exists in the comment block
	// TODO: Get File and Url - TODO, have these passed in as command-line params
	// TODO: Doc is the comment block after @API and before the first @Example annotation

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
