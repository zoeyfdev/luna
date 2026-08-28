package neoparser

var Scopes = []Scope {
	Scope {
		ID: 0,
		Parent: -1,
	},
}

var CurrentFunction *Variable
