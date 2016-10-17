package main

import (
	"testing"
	"runtime/debug"
	//"os"
	"strings"
	"fmt"
	"regexp"
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
		fmt.Printf("The test failed at:\n  %s\n", stackTraces[testinRunnerStackIndex ])
	}
}
