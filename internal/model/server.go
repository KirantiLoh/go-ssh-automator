package model

type Server struct {
	Username string   `json:"username"`
	IP       string   `json:"ip"`
	Port     string   `json:"port"`
	Commands []string `json:"commands"`
}
