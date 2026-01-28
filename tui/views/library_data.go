package views

var ShelfOrder = []string{
	"to_read",
	"currently_reading",
	"read",
}

var MockShelves = map[string][]string{
	"to_read": {
		"Dune",
		"1984",
		"Emma",
	},
	"currently_reading": {
		"Foundation",
	},
	"read": {
		"Frankenstein",
		"Dracula",
	},
}
