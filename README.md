
# NamingReview fo Golang

This is a static analysis tool for Golang to review the variable names, function names, and  type names in your project and find out the bad practices for naming.

This tool check whether you are properly naming your identifiers based on [The Online Plain Text English Dictionary](http://www.mso.anu.edu.au/~ralph/OPTED/).

## Install
```
go get -u github.com/yuya-okada/namingreview/cmd/namingreview
```

## How to use
### Review your identifiers

```
go vet -vettool=${which namingreview}
```

### Custom Dictionary File
To define or rewrite some words, you can put `reviewCustomDict.json` in your project directory. See the default dictionary [1](https://github.com/yuya-okada/namingreview/blob/master/dictionary.json) [2](https://github.com/yuya-okada/namingreview/blob/master/coding_word_dictionary.json) to learn how to write your dictionary.


## Currently Implemented Rules
### Variable Name
- Boolean variable names should contain verbs （ex. selected->isSelected, updatable->canUpdate）. It's desirable for them to start with verbs but it is not necessary (ex. you can use such as "userCanAccess").
-Variable names that are not boolean should contain at least one noun.
- Variable names that are not boolean should start with nouns, adjectives, or participles.

- The last noun in a name of Array or Slice should be "list", "array", "slice", or plural. On the contrary, the last noun NOT in a name of Array or Slice should be singular. And you can use both for Map.

- You can use a variable name with just 1 character if you need for readability.

### Type name
- Should contain at least one noun
- Should start with either a noun, adjective, or participle
- You can use a variable name with just 1 character if you need for readability.

 ### Fucntion name
 - Should start with a verb (main function is an exception)
 
