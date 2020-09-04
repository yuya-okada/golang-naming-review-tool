package a


var beautifulBeautiful int	// want "The variable name should contain at least one noun"
var eatApple =0		// want "The variable name should start with a noun or an adjective"

var apple []string // want "The final noun in the name of Array or Slice must be plural"
var pineapples int // want "The final noun not in the name of Array or Slice must be singular"

var selected = true // want "The boolean variable name should start with a verb.  ex. selected->isSelected, updatable->canUpdate"
var isSelected = true

func main() {
	beautifulBeautiful := 0	// want "The variable name should contain at least one noun"
	eatApple:=0		// want "The variable name should start with a noun or an adjective"
	appleContains:=0
	abandoningMan:=0

	print(beautifulBeautiful, eatApple, appleContains, abandoningMan)

	var gopher int
	print(gopher)
}

func playBaseball(beautifulBeautiful int) { // want "The variable name should contain at least one noun"
	return
}

func max () {	// want "The function name should start with a verb"
	return
}

