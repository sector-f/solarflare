package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type config struct {
	ApiKey         string   `toml:"api_key"`
	Interface      string   `toml:"interface"`
	Zone           string   `toml:"zone"`
	Subdomains     []string `toml:"subdomains"`
	UpdateInterval duration `toml:"update_interval"`
}

func getConfig(path string) (*config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf config
	metadata, err := toml.Decode(string(file), &conf)
	if err != nil {
		return nil, err
	}
	missing := []string{}

	if !metadata.IsDefined("api_key") {
		missing = append(missing, "api_key")
	}

	if !metadata.IsDefined("interface") {
		missing = append(missing, "interface")
	}

	if !metadata.IsDefined("zone") {
		missing = append(missing, "zone")
	}

	if !metadata.IsDefined("subdomains") {
		missing = append(missing, "subdomains")
	}

	if !metadata.IsDefined("update_interval") {
		missing = append(missing, "update_interval")
	}

	if len(missing) > 0 {
		return nil, errors.New(fmt.Sprintf("Missing config variables: %s", strings.Join(missing, ", ")))
	}

	return &conf, nil
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
