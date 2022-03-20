package config

import (
	"os"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type Channel string

type Exchange struct {
	Env     string `yaml:"env"`
	Name    string `yaml:"name"`
	Scheme  string `yaml:"scheme"`
	Host    string `yaml:"host"`
	Path    string `yaml:"path"`
	Maxsize int    `yaml:"maxsize"`
	Message struct {
		Type       string    `yaml:"type"`
		ProductIds []string  `yaml:"product_ids"`
		Channels   []Channel `yaml:"channels"`
	} `yaml:"message"`
}

func NewExchange(configFilename string) (*Exchange, error) {
	data, err := os.ReadFile(configFilename)
	if err != nil {
		return nil, xerrors.Errorf("unable to read file: %s: %w", configFilename, err)
	}

	exchange := Exchange{}
	if err = yaml.Unmarshal(data, &exchange); err != nil {
		return nil, xerrors.Errorf("unable to parse yaml: %s: %w", configFilename, err)
	}
	return &exchange, nil
}
