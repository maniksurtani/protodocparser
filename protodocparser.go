package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/maniksurtani/protodocparser/impl"
)

// Regexps
var (
	startCommentRE      = regexp.MustCompile("^\\s*/\\*")
	endCommentRE        = regexp.MustCompile("\\s*\\*/\\s*$")
	rpcRE               = regexp.MustCompile("\\s*rpc\\s+")
	serviceRE           = regexp.MustCompile("\\s*service\\s+")
	serviceNameRE       = regexp.MustCompile("\\s*service\\s+(\\w+)\\s*\\{")
	rpcNameRE           = regexp.MustCompile("\\s*rpc\\s+(\\w+)\\s*\\(\\s*([\\w]+)\\s*\\)\\s+returns\\s+\\(\\s*(\\w+)\\s*\\)")
	pkgNameRE           = regexp.MustCompile("\\s*package\\s+([\\w.]+)\\s*;")
	apiAnnotationRE     = regexp.MustCompile("^\\s*\\*\\s*@API")
	exampleAnnoationRE  = regexp.MustCompile("^\\s*\\*\\s*@Example")
	annotationContentRE = regexp.MustCompile("\\(([^\\)]+)\\)")
	designDocRE         = regexp.MustCompile("design\\s*=\\s*\"([^\"]+)\"")
	orgRE               = regexp.MustCompile("org\\s*=\\s*\"([^\"]+)\"")
)

type ProtoFile struct {
	// Can be used to access the contents of the proto
	ProtoFileSource io.Reader

	// Metadata, purely for display in the output JSON
	ProtoFilePath string
	Url           string
	Sha           string
}

type ParsingContext struct {
	currentBlock           *impl.CommentBlock
	pkgName                string
	matched                bool
	apiAnnotation          string
	designDoc              string
	org                    string
	examples               []*impl.Example
	currentExample         []string
	currentExampleLanguage string
}

func NewParsingContext() *ParsingContext {
	return &ParsingContext{examples: make([]*impl.Example, 0)}
}

func (p *ParsingContext) closeCurrentExample() []string {
	if p.currentExample != nil {
		ret := p.currentExample
		p.examples = append(p.examples, &impl.Example{Language: p.currentExampleLanguage, Code: strings.Join(p.currentExample, "\n")})
		p.currentExample = nil
		p.currentExampleLanguage = ""
		return ret
	}
	return nil
}

func (p *ParsingContext) initializeNewExample(line string) {
	p.closeCurrentExample()
	p.currentExampleLanguage = extractLanguageFromExample(line)
	p.currentExample = make([]string, 0)
}

func (p *ParsingContext) addLineToCurrentExample(line string) {
	p.currentExample = append(p.currentExample, strings.Trim(line, "* "))
}

// ParseAsString parses proto manifests and returns a JSON string. Used externally also.
func ParseAsString(protoFiles []*ProtoFile) string {
	services := parse(protoFiles)
	bytes, err := json.Marshal(services)
	if err != nil {
		panic(fmt.Sprintf("Caught error %v trying to serialize %v into JSON.", err, services))
	}
	return string(bytes)
}

func (p *ParsingContext) createNewCommentBlock(ln int) {
	p.currentBlock = &impl.CommentBlock{}
	p.currentBlock.Start = ln
	p.currentBlock.Type = impl.OtherComment
}

func (p *ParsingContext) reset() {
	p.currentBlock = nil
	p.apiAnnotation = ""
	p.designDoc = ""
	p.org = ""
}

func addRpcToLastService(services []*impl.Service, commentBlock *impl.CommentBlock, lines []string, currentLine int) {
	rpc := impl.NewRpc()
	rpc.Name, rpc.Request, rpc.Response = rpcName(lines[currentLine])

	// TODO: set Options, Doc and Examples
	// TODO: Doc is the comment block before the first @Example annotation
	// TODO: @Examples exists in the comment block
	// TODO: Options are the protobuf options. This might be harder to figure out since they may be split across multiple lines. :/

	lastService := services[len(services)-1]
	lastService.Rpcs = append(lastService.Rpcs, rpc)
}

func addServiceToServices(services []*impl.Service, p *ParsingContext, lines []string, currentLine int) []*impl.Service {
	s := impl.NewService()
	s.Name = serviceName(lines[currentLine])
	if len(p.apiAnnotation) > 0 {
		s.Api = true
		s.Design = p.designDoc
		s.Org = p.org
	}

	// TODO: set Doc, File, and Url
	// TODO: if Org isn't set, attempt to "guess" what it might be by looking at the path/package of the proto, and look up in Registry
	// TODO: Get File and Url - TODO, have these passed in as params
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

func extractAnnotationContent(line string) string {
	return strip(string(annotationContentRE.FindStringSubmatch(line)[1]))
}

func extractLanguageFromExample(line string) string {
	kv := strings.Split(extractAnnotationContent(line), "=")
	key, value := strip(kv[0]), strings.Trim(strip(kv[1]), "\"")
	if key != "language" {
		panic(fmt.Sprintf("Not a valid @Example parameter: `%s`", key))
	}
	return value
}

func exractSingleRegex(regex *regexp.Regexp, line string) string {
	match := regex.FindStringSubmatch(line)
	if len(match) < 1 {
		return ""
	}
	return match[1]
}

func extractDesignDoc(line string) string {
	return exractSingleRegex(designDocRE, line)
}

func extractOrg(line string) string {
	return exractSingleRegex(orgRE, line)
}

func strip(s string) string {
	return strings.Trim(s, " ")
}

// Can test from here rather than ParseAsString, since it makes testing easier
func parse(protoFiles []*ProtoFile) []*impl.Service {
	services := make([]*impl.Service, 0)

	for _, p := range protoFiles {
		contents, err := ioutil.ReadAll(p.ProtoFileSource)
		if err != nil {
			// TODO do something
			panic(err)
		}

		services = parseLines(strings.Split(string(contents), "\n"), p, services)
	}

	return services
}

func parseLines(lines []string, profoFile *ProtoFile, services []*impl.Service) []*impl.Service {
	p := NewParsingContext()
	for ln, line := range lines {
		if p.pkgName == "" {
			p.pkgName, p.matched = matchPkgName(line)
			if p.matched {
				continue
			}
		}
		fmt.Printf("Line %v is [%v]\n", ln, line)
		if startCommentRE.MatchString(line) && p.currentBlock == nil {
			p.createNewCommentBlock(ln)
		} else if endCommentRE.MatchString(line) && p.currentBlock != nil {
			p.currentBlock.End = ln
		} else if rpcRE.MatchString(line) && p.currentBlock != nil && p.currentBlock.End > 0 {
			// Mark block as an RPC type.
			p.currentBlock.Type = impl.RpcComment
			addRpcToLastService(services, p.currentBlock, lines, ln)
			p.currentBlock = nil
		} else if serviceRE.MatchString(line) && p.currentBlock != nil && p.currentBlock.End > 0 {
			// Mark block as a Service type.
			p.currentBlock.Type = impl.ServiceComment
			services = addServiceToServices(services, p, lines, ln)
			p.reset()
			currentExample := p.closeCurrentExample()
			if currentExample != nil {
				lastService := services[len(services)-1]
				lastService.Examples = p.examples
				p.examples = make([]*impl.Example, 0)
			}

		} else if apiAnnotationRE.MatchString(line) && p.currentBlock != nil {
			p.apiAnnotation = line
			if annotationContentRE.MatchString(line) {
				annotationContent := extractAnnotationContent(line)
				if designDocRE.MatchString(line) {
					p.designDoc = extractDesignDoc(annotationContent)
				}
				if orgRE.MatchString(line) {
					p.org = extractOrg(annotationContent)
				}
			}

		} else if exampleAnnoationRE.MatchString(line) && p.currentBlock != nil {
			p.initializeNewExample(line)
		} else if p.currentBlock != nil && p.currentExample != nil {
			p.addLineToCurrentExample(line)
		} else {
			// Todo : remove this entire block
			fmt.Printf("What?: %v\n", line)
			fmt.Printf(">>>> p.currentBlock: %v\n\n", p.currentBlock)
			if p.currentBlock != nil {
				fmt.Printf(">>>> p.currentBlock.End: %v\n\n", p.currentBlock.End)
			}
			fmt.Printf(">>>> len(services): %d\n\n", len(services))
		}
	}

	// Add package names to services.
	for _, svc := range services {
		svc.Package = p.pkgName
	}

	return services
}
