package a


var beautifulBeautiful int	// want "Variable name should contain at least one noun"
var eatApple =0		// want "Variable name should start with a noun or an adjective"

func f() {
	beautifulBeautiful := 0	// want "Variable name should contain at least one noun"
	eatApple:=0		// want "Variable name should start with a noun or an adjective"
	appleContains:=0
	abandoningMan:=0

	print(beautifulBeautiful, eatApple, appleContains, abandoningMan)

	var gopher int
	print(gopher)
}

