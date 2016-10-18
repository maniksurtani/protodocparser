package main

import (
	"testing"
	d "runtime/debug"
	//"os"
	"strings"
	"fmt"
	"regexp"
	"io/ioutil"
	"github.com/maniksurtani/protodocparser/impl"
	"reflect"
)

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
	protoString, _ := ioutil.ReadFile("./sample.proto")
	output := parseString(string(protoString))

	expectedServices := make([]*impl.Service, 0)
	s := impl.NewService()
	s.Package = "squareup.test.stuff"
	s.Name = "MyService"
	s.Api = true
	rpc := impl.NewRpc()
	rpc.Name = "MyEndpoint"
	rpc.Request = "Request"
	rpc.Response = "Response"
	s.Rpcs = append(s.Rpcs, rpc)
	expectedServices = append(expectedServices, s)

	if !reflect.DeepEqual(output, expectedServices) {
		t.Errorf("Not the same: \n%+V\n%+V\n", output, expectedServices)
	}
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

/*
	Finds the line of the test that failed, and prints it.
 */
func printRelevantStacktrace() {
	stackTraces := strings.Split(string(d.Stack()[:]), "\n")
	testinRunnerStackIndex := -1
	for index, stackTraceLine := range stackTraces {
		matches, _ := regexp.MatchString("^testing.tRunner", stackTraceLine)
		if matches {
			testinRunnerStackIndex = index - 1
		}
	}
	if testinRunnerStackIndex >= 0 {
		fmt.Printf("The test failed at:\n  %s\n", stackTraces[testinRunnerStackIndex ])
	}
}
