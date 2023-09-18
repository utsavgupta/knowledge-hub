package entities

type Query struct {
	Question string
	DomainId string
	Concepts []Concept
}
