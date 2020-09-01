package golang_naming_review_tool

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

const dictionaryFileName = "dictionary.json"
const doc = "go_naming_review is ..."

var wordDict map[string]map[string]bool


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
	wordDict = loadDictionary(dictionaryFileName)

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
	words := GetWordList(name)

}

func GetWordList(name string) []string{
	var words []string
	wordStartIndex := 0
	for i, c := range name {
		if unicode.IsUpper(c) && i!= wordStartIndex{
			words = append(words, strings.ToLower(name[wordStartIndex:i]))
			wordStartIndex=i
		}
	}
	words = append(words, strings.ToLower(name[wordStartIndex:len(name)]))

	return words
}

func IsPlural(word string) bool{
	val, ok := wordDict[word]
	if ok {
		isPl, ok := val["pl"]
		if !ok {
			isPl = false
		}

		return isPl
	} else {
		return word[len(word)-1] ==  's'
	}
}