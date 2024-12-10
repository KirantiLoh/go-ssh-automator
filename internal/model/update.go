package model

type Update struct {
	Host       string
	Message    string
	IsComplete bool
	IsError    bool
}
