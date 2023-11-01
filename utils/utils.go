package utils

// TODO: лучше url.Values (или просто мапой) ??
type SearchParams struct {
	Start int `json:"start"`
	Step  int `json:"step"`
	Limit int `json:"limit"`
}
