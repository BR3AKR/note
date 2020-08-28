package cmd

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Config contains all configurations required for note to run
type Config struct {
	FullName string `yaml:"fullname"`
	Paths    struct {
		Blog    string `yaml:"blog"`
		Book    string `yaml:"book"`
		Morning string `yaml:"morning"`
	} `yaml:paths`
}

func parseConfigs(r io.Reader) (Config, error) {
	var cfg Config
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&cfg)
	return cfg, err
}

func getConfigFile() (*os.File, error) {
	if f, err := os.Open("./.note.yml"); err == nil {
		return f, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if f, err := os.Open(home + "/.config/note.yml"); err == nil {
		return f, nil
	}
	if f, err := os.Open(home + "/.note.yml"); err == nil {
		return f, nil
	}
	return nil, errors.New("unable to find a note configuration")
}
