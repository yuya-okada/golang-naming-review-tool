package a


var beautifulBeautiful int	// want "Variable name should contain at least one noun"
var eatApple =0		// want "Variable name should start with a noun or an adjective"

var apple []string // want "The final noun in the name of Array or Slice must be plural"
var pineapples int // want "The final noun not in the name of Array or Slice must be singular"

func main() {
	beautifulBeautiful := 0	// want "Variable name should contain at least one noun"
	eatApple:=0		// want "Variable name should start with a noun or an adjective"
	appleContains:=0
	abandoningMan:=0

	print(beautifulBeautiful, eatApple, appleContains, abandoningMan)

	var gopher int
	print(gopher)
}

func playBaseball() {
	return
}

func baseball () {	// want "Function name should start with a verb"
	return
}
