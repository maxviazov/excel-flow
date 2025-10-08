package config

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

type LoggerConfig struct {
	Level     string `mapstructure:"level"`
	File      string `mapstructure:"file"`
	Console   bool   `mapstructure:"console"`
	Color     bool   `mapstructure:"color"`
	Timestamp bool   `mapstructure:"timestamp"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
}
