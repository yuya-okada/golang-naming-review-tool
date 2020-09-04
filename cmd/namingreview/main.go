package main

import (
	"github.com/yuya-okada/namingreview"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(namingreview.Analyzer) }

