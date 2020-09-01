package main

import (
	"github.com/yuya-okada/golang_naming_review_tool"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(golang_naming_review_tool.Analyzer) }

