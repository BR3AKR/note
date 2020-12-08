package cmd

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Config contains all configurations required for note to run
type Config struct {
	FullName string      `yaml:"fullname"`
	Paths    ConfigPaths `yaml:"paths"`
}

type ConfigPaths struct {
	Base    string `yaml:"base"`
	Blog    string `yaml:"blog"`
	Book    string `yaml:"book"`
	Morning string `yaml:"morning"`
	Meeting string `yaml:"meeting"`
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
	filename := home + "/.config/note.yml"
	fmt.Println("unable to find a note configuration, creating one at: " + filename)
	config := &Config{
		FullName: prompt("full name: "),
		Paths: ConfigPaths{
			Base:    prompt("base path (~/notes): ", "~/notes"),
			Blog:    "~/blog",
			Book:    "books",
			Morning: "morning-pages",
			Meeting: "inbox",
		},
	}
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(out)
	if err != nil {
		return nil, err
	}

	return file, nil
}
