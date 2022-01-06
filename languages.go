package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type languageInfo struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Enable bool   `yaml:"enable"`
}

var Languages []languageInfo

func init() {
	bytes, err := ioutil.ReadFile(WorkDir + "/languages.yaml")
	if err != nil {
		panic(err)
	}
	Languages = make([]languageInfo, 0)
	if err := yaml.Unmarshal(bytes, &Languages); err != nil {
		panic(err)
	}
}

func GetAvailableLanguageIDs() []string {
	ids := make([]string, 0)
	for _, info := range Languages {
		if info.Enable {
			ids = append(ids, fmt.Sprintf("%v", info.ID))
		}
	}
	return ids
}
