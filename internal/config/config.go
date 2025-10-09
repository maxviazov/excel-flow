package config

// AppConfig holds application-specific configuration.
type AppConfig struct {
	// Name is the name of the application.
	Name string `mapstructure:"name"`
	// Version is the version of the application.
	Version string `mapstructure:"version"`
	// Port is the port the application will listen to on.
	Port int `mapstructure:"port"`
}

// LoggerConfig holds logging configuration.
type LoggerConfig struct {
	// Level is the logging level (e.g., "debug", "info", "warn", "error").
	Level string `mapstructure:"level"`
	// File is the path to the log file.
	File string `mapstructure:"file"`
	// Console determines if logs should be output to the console.
	Console bool `mapstructure:"console"`
	// Color enables or disables colored output in console logs.
	Color bool `mapstructure:"color"`
	// Timestamp enables or disables timestamps in logs.
	Timestamp bool `mapstructure:"timestamp"`
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	// Host is the database host.
	Host string `mapstructure:"host"`
	// Port is the database port.
	Port int `mapstructure:"port"`
	// User is the database user.
	User string `mapstructure:"user"`
	// Password is the database user's password.
	Password string `mapstructure:"password"`
	// Name is the name of the database.
	Name string `mapstructure:"name"`
}

// Config is the root configuration struct that holds all other configuration structs.
type Config struct {
	// App holds the application configuration.
	App AppConfig `mapstructure:"app"`
	// Logger holds the logging configuration.
	Logger LoggerConfig `mapstructure:"logging"`
	// Database holds the database configuration.
	Database DatabaseConfig `mapstructure:"database"`
}
