package model

type Config struct {
	Servers       []Server      `json:"servers"`
	DefaultConfig DefaultConfig `json:"defaults"`
}

type DefaultConfig struct {
	Username     string `json:"username"`
	Port         string `json:"port"`
	IdentityFile string `json:"identity_file"`
}
