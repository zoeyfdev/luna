package neoparser

var ScopeTicker int = 1
func CreateScope(Parent int) int {
	ID := ScopeTicker
	Scopes = append(Scopes, Scope {
		Parent: Parent,
		ID: ID,
	})
	ScopeTicker++
	return ID
} 
