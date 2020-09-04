package golang_naming_review_tool_test

import (
	"strings"
	"testing"

	"github.com/yuya-okada/golang_naming_review_tool"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testData := analysistest.TestData()
	analysistest.Run(t, testData, golang_naming_review_tool.Analyzer, "a")
}

func TestGetWordList(t *testing.T) {

	answers := golang_naming_review_tool.GetWordList("fooBarBar")
	if !(len(answers)==3 && answers[0]=="foo" &&  answers[1]=="bar" &&  answers[2]=="bar") {
		t.Errorf("fooBarBar  = %s; want ['foo','bar','bar']", strings.Join(answers, ","))
	}
	answers = golang_naming_review_tool.GetWordList("Foo")
	if !(len(answers)==1 && answers[0]=="foo") {
		t.Errorf("foo  = %s; want ['foo']", strings.Join(answers, ","))
	}
}


func TestIsSpecificPartOfSpeech(t *testing.T) {
	isSpecificPartOfSpeech := golang_naming_review_tool.IsSpecificPartOfSpeech("understand", "n")
	if isSpecificPartOfSpeech {
		t.Errorf("understand is noun = %t; want false", isSpecificPartOfSpeech)
	}
	isSpecificPartOfSpeech = golang_naming_review_tool.IsSpecificPartOfSpeech("pineapple", "n")
	if !isSpecificPartOfSpeech {
		t.Errorf("pineapple is noun = %t; want true", isSpecificPartOfSpeech)
	}
	isSpecificPartOfSpeech = golang_naming_review_tool.IsSpecificPartOfSpeech("understanding", "n")
	if !isSpecificPartOfSpeech {
		t.Errorf("understanding is noun = %t; want true", isSpecificPartOfSpeech)
	}
}

