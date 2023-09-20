package entities

type Response struct {
	Query    Query    `json:"query"`
	Response string   `json:"response"`
	Sources  []string `json:"sources"`
}
