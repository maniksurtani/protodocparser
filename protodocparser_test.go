package protodocparser

import (
	"runtime/debug"
	"testing"
	//"os"
	"encoding/json"
	"fmt"
	"github.com/maniksurtani/protodocparser/impl"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func TestExtractAnnotation(t *testing.T) {
	assertEqualStrings(extractAnnotationContent("@API(   content here  )"), "content here", t)
	assertPanic(t, func() { extractAnnotationContent("no content here. it will panic") })
}

func TestExtractLanguageFromExample(t *testing.T) {
	assertEqualStrings(extractLanguageFromExample(`@Example(language="java")`), "java", t)
	assertEqualStrings(extractLanguageFromExample(`@Example(language="go")`), "go", t)
	assertEqualStrings(extractLanguageFromExample("@Example(language   =   java)"), "java", t)
	assertPanic(t, func() { extractLanguageFromExample("@Example(badparam=java)") })
}

func TestExtractDesignDoc(t *testing.T) {
	assertEqualStrings(extractDesignDoc(`design="http://example.com/design.html"`), "http://example.com/design.html", t)
	assertEqualStrings(extractDesignDoc(`design   =   "path"`), "path", t)
	assertEqualStrings(extractDesignDoc("no design doc here"), "", t)
}

func TestExtractOrg(t *testing.T) {
	assertEqualStrings(extractOrg(`org="theorg"`), "theorg", t)
	assertEqualStrings(extractOrg(`org   =   "theorg"`), "theorg", t)
	assertEqualStrings(extractOrg("no org here"), "", t)
}

func TestStartComment(t *testing.T) {
	assertTrue(startCommentRE.MatchString("/**"), t)
	assertTrue(startCommentRE.MatchString("          /**"), t)
	assertTrue(startCommentRE.MatchString("/*"), t)
	assertFalse(startCommentRE.MatchString("//"), t)
	assertFalse(startCommentRE.MatchString("this is not a comment /**"), t)
}

func TestEndComment(t *testing.T) {
	assertTrue(endCommentRE.MatchString("*/"), t)
	assertTrue(endCommentRE.MatchString("*/    "), t)
	assertFalse(endCommentRE.MatchString("*/    s"), t)
}

func TestIsRpc(t *testing.T) {
	assertTrue(rpcRE.MatchString("   rpc MyRPC"), t)
	assertFalse(rpcRE.MatchString("   rpc"), t)
}

func TestIsService(t *testing.T) {
	assertTrue(serviceRE.MatchString("   service MyService"), t)
	assertFalse(serviceRE.MatchString("   service"), t)
	assertFalse(serviceRE.MatchString("   serviceWithTypo MyService"), t)
}

func TestApiAnnotation(t *testing.T) {
	assertTrue(apiAnnotationRE.MatchString("* @API()"), t)
	assertTrue(apiAnnotationRE.MatchString(` * @API(design="http://link.to.design.com/design.html", org="payments")`), t)
	assertTrue(apiAnnotationRE.MatchString("* @API"), t)
	assertFalse(apiAnnotationRE.MatchString("* @Api()"), t)
	assertFalse(apiAnnotationRE.MatchString("* @api()"), t)
	assertFalse(apiAnnotationRE.MatchString("* api"), t)
}

func TestExampleAnnoations(t *testing.T) {
	assertTrue(exampleAnnoationRE.MatchString("* @Example()"), t)
	assertTrue(exampleAnnoationRE.MatchString(` * @Example(language="java")`), t)
	assertTrue(exampleAnnoationRE.MatchString("* @Example"), t)
	assertFalse(exampleAnnoationRE.MatchString("* @example()"), t)
}

func TestParseKeyValues(t *testing.T) {
	assertTrue(exampleAnnoationRE.MatchString("* @Example()"), t)
}

func TestServiceNames(t *testing.T) {
	if serviceName("service S{") != "S" {
		t.Errorf("Expected 'S' but was %v", serviceName("service S{"))
	}

	if serviceName("service S {") != "S" {
		t.Errorf("Expected 'S' but was %v", serviceName("service S {"))
	}

	if serviceName("    service    S_s   {") != "S_s" {
		t.Errorf("Expected 'S_s' but was %v", serviceName("    service    S_s   {"))
	}
}

func TestParseSimpleProto(t *testing.T) {
	protoFile, _ := os.Open("./sample.proto")

	pf := &ProtoFile{
		ProtoFileSource: protoFile,
		ProtoFilePath:   "./sample.proto",
		Url:             "http://some.repo/sample.proto",
		Sha:             "ABCD1234"}

	output := parse([]*ProtoFile{pf})

	expectedServices := make([]*impl.Service, 0)
	s := impl.NewService()
	s.Package = "squareup.test.stuff"
	s.Name = "MyService"
	s.Api = true
	s.Design = "http://example.com/design.html"
	s.Org = "organization"
	s.Doc = "The doc for this service\nThe second line of the doc"
	rpc := impl.NewRpc()
	rpc.Name = "MyEndpoint"
	rpc.Request = "Request"
	rpc.Response = "Response"
	rpc.Doc = "The doc for MyEndpoint\n"
	gocode := "conn := createRpcConnection()\nresponse, err := conn.MyEndpoint(&Request{})"
	s.Examples = append(s.Examples,
		&impl.Example{Language: "java", Code: `String s = new String("Blah");`},
		&impl.Example{Language: "go", Code: gocode},
	)
	rpc.Examples = append(rpc.Examples, &impl.Example{Language: "java", Code: `Future<Response> rsp = makeRequest();`})
	s.Rpcs = append(s.Rpcs, rpc)
	expectedServices = append(expectedServices, s)
	if !reflect.DeepEqual(output, expectedServices) {
		t.Errorf("Not the same: \n%+s\n%+s\n", asJson(output), asJson(expectedServices))
	}
}

func asJson(v interface{}) string {
	//out, _ := json.MarshalIndent(v, "", "  ")
	out, _ := json.Marshal(v)
	return string(out)
}

func TestRpcNames(t *testing.T) {
	// TODO
}

func assertTrue(expr bool, t *testing.T) {
	if !expr {
		printRelevantStacktrace()
		t.Fail()
	}
}

func assertFalse(expr bool, t *testing.T) {
	if expr {
		printRelevantStacktrace()
		t.Fail()
	}
}

func assertEqualStrings(input string, expected string, t *testing.T) {
	if input != expected {
		printRelevantStacktrace()
		t.Errorf("\nExpected: `%s`\n but got: `%s`", expected, input)
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()
	f()
}

/*
	Finds the line of the test that failed, and prints it.
*/
func printRelevantStacktrace() {
	stackTraces := strings.Split(string(debug.Stack()[:]), "\n")
	testinRunnerStackIndex := -1
	for index, stackTraceLine := range stackTraces {
		matches, _ := regexp.MatchString("^testing.tRunner", stackTraceLine)
		if matches {
			testinRunnerStackIndex = index - 1
		}
	}
	if testinRunnerStackIndex >= 0 {
		fmt.Printf("The test failed at:\n  %s\n", stackTraces[testinRunnerStackIndex])
	}
}
