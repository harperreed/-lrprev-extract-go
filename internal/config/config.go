package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputDir        string `yaml:"input_dir"`
	InputFile       string `yaml:"input_file"`
	OutputDirectory string `yaml:"output_directory"`
	LightroomDB     string `yaml:"lightroom_db"`
	IncludeSize     bool   `yaml:"include_size"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}
