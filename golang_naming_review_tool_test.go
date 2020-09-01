package golang_naming_review_tool_test

import (
	"testing"

	"github.com/yuya-okada/golang_naming_review_tool"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, golang_naming_review_tool.Analyzer, "a")
}

func TestIsPlural(t *testing.T) {

}
