package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token    string `json:"token"`
	Prefix   string `json:"prefix"`
	EmojiDir string `json:"emoji_dir"`
}

func LoadConfigFile(file string) (config *Config, err error) {
	var (
		rawConfig []byte
	)

	rawConfig, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	config = new(Config)

	err = json.Unmarshal(rawConfig, &config)

	return
}
