package database

// Config ...
type Config struct {
	DatabaseAddress string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		DatabaseAddress: "forum.db",
	}
}
