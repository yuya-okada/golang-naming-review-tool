package namingreview

import (
	"encoding/json"
	"fmt"
	pluralizePkg "github.com/gertd/go-pluralize"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

const dictionaryFileName = "dictionary.json"
const codingWordDictionaryFileName = "coding_word_dictionary.json"

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
	initWordDict()

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilterArray := []ast.Node{
		(*ast.GenDecl)(nil),
		(*ast.AssignStmt)(nil),
	}
	inspect.Preorder(nodeFilterArray, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.GenDecl:
			reviewGenDecl(pass, n)
		case *ast.AssignStmt:
			reviewAssignStmt(pass, n)
		}

	})

	nodeFilterArray = []ast.Node{
		(*ast.FuncDecl)(nil),
	//	(*ast.FuncLit)(nil),
	}
	inspect.Preorder(nodeFilterArray, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			reviewFuncDecl(pass, n)
		}
	})

	return nil, nil
}

func printIfError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func initWordDict() {
	wordDict = loadDictionary(dictionaryFileName)
	codingWordDict := loadDictionary(codingWordDictionaryFileName)
	for codingWord, partOfSpeechDict := range codingWordDict {
		if _, ok := wordDict[codingWord]; ok {
			for partOfSpeech, value := range partOfSpeechDict {
				wordDict[codingWord][partOfSpeech] = value
			}
		} else {
			wordDict[codingWord] = partOfSpeechDict
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("reviewCustomDict.json does not found")
	} else {
		customWordDict := loadDictionary(path.Join(dir, "reviewCustomDict.json"))

		for customWord, partOfSpeechDict := range customWordDict {
			if _, ok := wordDict[customWord]; ok {
				for partOfSpeech, value := range partOfSpeechDict {
					wordDict[customWord][partOfSpeech] = value
				}
			} else {
				wordDict[customWord] = partOfSpeechDict
			}
		}
	}
}

func loadDictionary(fileName string) map[string]map[string]bool {
	jsonFile, err := os.Open(fileName)
	printIfError(err)
	defer jsonFile.Close()

	byteValues, err := ioutil.ReadAll(jsonFile)
	printIfError(err)

	var words map[string]map[string]bool
	json.Unmarshal(byteValues, &words)
	return words
}


func reviewAssignStmt(pass *analysis.Pass, statement *ast.AssignStmt) {
	if statement.Tok == token.DEFINE {
		for _, expr := range statement.Lhs {
			id, ok := expr.(*ast.Ident)
			if ok {
				error := reviewVariableName(pass, id)
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
	// ValueSpec
	valSpec, isValueSpec := spec.(*ast.ValueSpec)
	if isValueSpec {
		for _, id := range valSpec.Names {
			if id.Name != "_" {
				error := reviewVariableName(pass, id)
				if error != nil {
					pass.Reportf(id.Pos(), error.Error())
				}
			}
		}
	}

	// TypeSpec
	typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
	if isTypeSpec {
		error := reviewVariableName(pass, typeSpec.Name)
		if error != nil {
			pass.Reportf(typeSpec.Pos(), error.Error())
		}
	}
}

func reviewVariableName(pass *analysis.Pass,id *ast.Ident) error{
	if id == nil {
		return nil
	}
	name := id.Name
	obj := pass.TypesInfo.ObjectOf(id)

	// Single-character words are allowed (Using in large scope is deprecated)
	if len(name) < 1 {
		return nil
	}
	words := GetWordList(name)

	// The boolean variable name should contain a verb.
	if obj != nil && types.Identical(obj.Type(), types.Typ[types.Bool]) {
		containsVerb:= false
		for _, word := range words {
			if IsVerb(word) {
				containsVerb = true
			}
		}
		if !containsVerb {
			return NewNamingError("The boolean variable name should start with a verb.  ex. selected->isSelected, updatable->canUpdate")
		}

		return nil
	}

	// The first word in the variable name should be a noun or an adjective
	if !(IsNoun(words[0]) || IsAdjective(words[0])) {
		return NewNamingError("The variable name should start with a noun or an adjective")
	}

	// Check variable whether name contains at least one noun
	containsNoun := false
	for _, word := range words {
		if IsNoun(word) {
			containsNoun = true
		}
	}

	if !containsNoun {
		return NewNamingError("The variable name should contain at least one noun")
	}

	// The final noun in the name of Array or Slice must be plural
	// On the contrary, the final noun NOT in the name of Array or Slice must be singular
	// For the name of Map, they don't matter
	if obj != nil {
		_, isSlice := obj.Type().(*types.Slice)
		_, isArray := obj.Type().(*types.Array)
		_, isMap := obj.Type().(*types.Map)
		if !isMap && len(name) > 1 {
			if isSlice || isArray {
				if  !canBeArrayName(name) {
					return NewNamingError("The final noun in the name of Array or Slice should be 'list', 'array', 'slice' or plural")
				}
			} else if !isSingularName(name) {
				return NewNamingError("The final noun not in the name of Array or Slice shouldn't be 'list', 'array', 'slice' or plural")
			}
		}
	}


	return nil
}


func reviewFuncDecl(pass *analysis.Pass, decl *ast.FuncDecl) {
	if decl.Name == nil || decl.Name.Name == "" {
		return
	}
	name := decl.Name.Name
	if name == "main" || name == "init" {
		return
	}
	words := GetWordList(name)
	if !IsVerb(words[0]) {
		error := NewNamingError("The function name should start with a verb")
		pass.Reportf(decl.Pos(), error.Error())
	}

	for _, field := range decl.Type.Params.List {
		for _, name := range field.Names {
			error := reviewVariableName(pass, name)
			if error != nil {
				pass.Reportf(decl.Pos(), error.Error())
			}
		}
	}
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


func canBeArrayName(name string) bool {
	finalNoun := ""
	for _, word := range GetWordList(name) {
		if IsNoun(word) {
			finalNoun = word
		}
	}


	if isPlural(finalNoun) {
		return true
	} else {
		return finalNoun == "list" || finalNoun == "slice" || finalNoun == "array"
	}
}

func isSingularName(name string) bool {
	finalNoun := ""
	for _, word := range GetWordList(name) {
		if IsNoun(word) {
			finalNoun = word
		}
	}

	if finalNoun == "list" || finalNoun == "slice" || finalNoun == "array"{
		return false
	} else{
		return isSingular(finalNoun)
	}
}


func IsSpecificPartOfSpeech(word string, partOfSpeech string) bool {
	if isPlural(word) {
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

func isPlural(word string) bool {
	if partOfSpeechDict, ok := wordDict[word]; ok {
		if _, ok := partOfSpeechDict["pl"]; ok {
			return true
		}
	}
	return pluralize.IsPlural(word)
}
func isSingular(word string) bool {
	if partOfSpeechDict, ok := wordDict[word]; ok {
		if _, ok := partOfSpeechDict["s"]; ok {
			return true
		}
	}
	return pluralize.IsSingular(word)
}