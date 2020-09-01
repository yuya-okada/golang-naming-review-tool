package golang_naming_review_tool

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"io/ioutil"
	"os"
)

const dictionaryFileName = "dictionary.json"
const doc = "go_naming_review is ..."

var words map[string]map[string]bool


// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "go_naming_review",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func printIfError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func loadDictionary(fileName string) map[string]map[string]bool {
	jsonFile, err := os.Open(fileName)
	printIfError(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	printIfError(err)

	var words map[string]map[string]bool
	json.Unmarshal(byteValue, &words)
	return words
}

func run(pass *analysis.Pass) (interface{}, error) {
	words = loadDictionary(dictionaryFileName)

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			decl, ok := decl.(*ast.GenDecl)
			if ok {
				reviewGenDecl(decl)
			}


		}
	}

	return nil, nil
}


func reviewGenDecl(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		reviewSpec(spec)
	}
}

func reviewSpec(spec ast.Spec) {
	valSpec, ok := spec.(*ast.ValueSpec)
	if ok {
		for _, id := range valSpec.Names {
			if id.Name != "_" {
				reviewValueName(id.Name)
			}
		}
	}
}

func reviewValueName(name string) {
	if len(name) < 1 {
		return
	}

}