package main

import (
	"fmt"
	"github.com/ismdeep/rand"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type config struct {
	BaseURL    string `yaml:"base_url"`    // OJ网址
	SecureCode string `yaml:"secure_code"` // 服务器用于判题机认证的Secure Code，存放在 justoj 项目的 .env 文件中
}

var ClientName string
var WorkDir string
var Config *config

func init() {
	WorkDir, _ = os.Getwd()
	if os.Getenv("JUSTOJ_CORE_ROOT") != "" {
		WorkDir = os.Getenv("JUSTOJ_CORE_ROOT")
	}

	ClientName, _ = os.Hostname()
	if ClientName == "" {
		ClientName = rand.HexStr(8)
	}

	c := &config{}
	data, err := ioutil.ReadFile(fmt.Sprintf("%v/config.yaml", WorkDir))
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		panic(err)
	}

	Config = c
}
