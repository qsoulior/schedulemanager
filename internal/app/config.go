package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Port           int    `json:"port"`
		AllowedOrigins string `json:"allowed_origins"`
	} `json:"server"`

	Mongo struct {
		URI string `json:"uri"`
	} `json:"mongo"`

	Moodle struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		CourseId int    `json:"course_id"`
	} `json:"moodle"`
}

func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := new(Config)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
