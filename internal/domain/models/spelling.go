package models

type SpellError struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}
