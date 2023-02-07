package db

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	DBName     string `yaml:"dbname" json:"dbname"`
	DBHost     string `yaml:"dbhost" json:"dbhost"`
	DBUser     string `yaml:"dbuser" json:"dbuser"`
	DBPassword string `yaml:"dbpassword" json:"dbpassword"`
	Charset    string `yaml:"charset" json:"charset"`
	ParseTime  bool   `yaml:"parse_time" json:"parseTime"`
	Loc        string `yaml:"loc" json:"loc"`
}

// Validate validates the config struct
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DBName, validation.Required),
		validation.Field(&c.DBHost, validation.Required),
		validation.Field(&c.DBUser, validation.Required),
		validation.Field(&c.Loc, validation.Required),
	)
}
