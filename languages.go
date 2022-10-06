package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type languageInfo struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Enable bool   `yaml:"enable"`
}

// Languages slice
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

// GetAvailableLanguageIDs get language id list
func GetAvailableLanguageIDs() []string {
	ids := make([]string, 0)
	for _, info := range Languages {
		if info.Enable {
			ids = append(ids, fmt.Sprintf("%v", info.ID))
		}
	}
	return ids
}
