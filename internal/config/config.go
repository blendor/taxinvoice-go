package config

type Config struct {
    DatabaseURL string
    ServerPort  string
}

func LoadConfig() *Config {
    // In a real app, you'd load this from env vars or a config file
    return &Config{
        DatabaseURL: "postgres://username:password@localhost/dbname?sslmode=disable",
        ServerPort:  "8080",
    }
}