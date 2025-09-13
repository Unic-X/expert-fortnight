package model

type DBConfig struct {
	URL           string `json:"url"`
	Host          string `json:"host"`
	Port          string `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	DBName        string `json:"dbname"`
	SSLMode       string `json:"sslmode"`
	MigrationsDir string `json:"migrations_dir"`
}

type Config struct {
	DB DBConfig
}
