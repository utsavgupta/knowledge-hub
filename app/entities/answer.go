package entities

type Answer struct {
	Query    Query  `json:"query"`
	Response string `json:"response"`
}
