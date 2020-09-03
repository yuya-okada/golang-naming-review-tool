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


func TestisSpecificPartOfSpeech(t *testing.T) {
	ans := golang_naming_review_tool.IsSpecificPartOfSpeech("understand", "n")
	if ans {
		t.Errorf("understand is noun = %t; want false", ans)
	}
	ans = golang_naming_review_tool.IsSpecificPartOfSpeech("pineapple", "n")
	if !ans {
		t.Errorf("pineapple is noun = %t; want true", ans)
	}
	ans = golang_naming_review_tool.IsSpecificPartOfSpeech("understanding", "n")
	if !ans {
		t.Errorf("understanding is noun = %t; want true", ans)
	}
}

