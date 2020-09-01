package golang_naming_review_tool_test

import (
	"strings"
	"testing"

	"github.com/yuya-okada/golang_naming_review_tool"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, golang_naming_review_tool.Analyzer, "a")
}

func TestGetWordList(t *testing.T) {

	ans := golang_naming_review_tool.GetWordList("fooBarBar")
	if !(len(ans)==3 && ans[0]=="foo" &&  ans[1]=="bar" &&  ans[2]=="bar") {
		t.Errorf("fooBarBar  = %s; want ['foo','bar','bar']", strings.Join(ans, ","))
	}
	ans = golang_naming_review_tool.GetWordList("Foo")
	if !(len(ans)==1 && ans[0]=="foo") {
		t.Errorf("foo  = %s; want ['foo']", strings.Join(ans, ","))
	}
}

func TestIsVerb(t *testing.T) {
	ans := golang_naming_review_tool.IsVerb("play")
	if !ans {
		t.Errorf("play is verb = %t; want true", ans)
	}
	ans = golang_naming_review_tool.IsVerb("pineapple")
	if ans {
		t.Errorf("pineapple is verb = %t; want false", ans)
	}
}

func TestIsPlural(t *testing.T) {
	ans := golang_naming_review_tool.IsPlural("apple")
	if ans {
		t.Errorf("apples is plural = %t; want false", ans)
	}
	ans = golang_naming_review_tool.IsPlural("apples")
	if !ans {
		t.Errorf("apples is plural = %t; want true", ans)
	}
	ans = golang_naming_review_tool.IsPlural("children")
	if !ans {
		t.Errorf("children is plural = %t; want true", ans)
	}
	ans = golang_naming_review_tool.IsPlural("foobar")
	if ans {
		t.Errorf("foobar is plural = %t; want false", ans)
	}
	ans = golang_naming_review_tool.IsPlural("foobars")
	if !ans {
		t.Errorf("foobars is plural = %t; want true", ans)
	}
}
