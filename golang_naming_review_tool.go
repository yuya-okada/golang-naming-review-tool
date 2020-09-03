package golang_naming_review_tool

import (
	"encoding/json"
	"fmt"
	pluralizePkg "github.com/gertd/go-pluralize"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

const dictionaryFileName = "dictionary.json"
const doc = "go_naming_review is ..."
var pluralize = pluralizePkg.NewClient()

var wordDict map[string]map[string]bool

func NewNamingError(text string) error {
	return &NamingError{text}
}
type NamingError struct {
	s string
}
func (e *NamingError) Error() string {
	return e.s
}

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "go_naming_review",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	wordDict = loadDictionary(dictionaryFileName)

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
		(*ast.AssignStmt)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.GenDecl:
			reviewGenDecl(pass, n)
		case *ast.AssignStmt:
			reviewAssignStmt(pass, n)
		}

	})

	//
	//for _, f := range pass.Files {
	//	for _, decl := range f.Decls {
	//		decl, ok := decl.(*ast.GenDecl)
	//		if ok {
	//			reviewGenDecl(pass, decl)
	//		}
	//	}
	//}
	return nil, nil
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


func reviewAssignStmt(pass *analysis.Pass, statement *ast.AssignStmt) {
	if statement.Tok == token.DEFINE {
		for _, expr := range statement.Lhs {
			id, ok := expr.(*ast.Ident)
			if ok {
				error := reviewVariableName(id.Name)
				if error != nil {
					pass.Reportf(id.Pos(), error.Error())
				}
			}
		}
	}
}


func reviewGenDecl(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		reviewSpec(pass, spec)
	}
}

func reviewSpec(pass *analysis.Pass, spec ast.Spec) {
	valSpec, ok := spec.(*ast.ValueSpec)
	if ok {
		for _, id := range valSpec.Names {
			if id.Name != "_" {
				error := reviewVariableName(id.Name)
				if error != nil {
					pass.Reportf(id.Pos(), error.Error())
				}
			}
		}
	}
}

func reviewVariableName(name string) error{
	// A single-character word is allowed (Using in large scope is deprecated)
	if len(name) < 1 {
		return nil
	}
	words := GetWordList(name)

	// The first word in the variable name should be a noun or an adjective
	if !(IsNoun(words[0]) || IsAdjective(words[0])) {
		return NewNamingError("Variable name should start with a noun or an adjective")
	}

	// Check variable whether name contains at least one noun
	nounContainsNoun := false
	for _, word := range words {
		if IsNoun(word) {
			nounContainsNoun = true
		}
	}
	if !nounContainsNoun {
		return NewNamingError("Variable name should contain at least one noun")
	}
	return nil
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


func IsSpecificPartOfSpeech(word string, partOfSpeech string) bool {
	if pluralize.IsPlural(word) {
		word = pluralize.Singular(word)
	}

	types, ok := wordDict[word]
	if !ok {
		return true
	}
	_, ok = types[partOfSpeech]
	return ok
}

func IsNoun(word string) bool {
	return IsSpecificPartOfSpeech(word, "n")
}
func IsVerb(word string) bool {
	return IsSpecificPartOfSpeech(word, "v")
}
func IsVerbBareForm(word string) bool {
	return IsSpecificPartOfSpeech(word, "vb")
}
func IsAdjective(word string) bool {
	return IsSpecificPartOfSpeech(word, "a") || IsVerbBareForm(word)
}


