package utils

import (
	"encoding/json"
	"io/ioutil"
)

type Language struct {
	Id       int    `json:"id"`
	Language string `json:"lang"`
	Enabled  bool   `json:"enabled"`
	ShellCmd string `json:"shell"`
}

type Config struct {
	BaseUrl         string     `json:"base_url"`          // OJ网址
	QuerySize       int        `json:"query_size"`        // 查询时间
	RunningSize     int        `json:"running_size"`      // 同时运行并发量
	SleepTime       int        `json:"sleep_time"`        // 服务器无用户提交后程序休眠时间
	ClientName      string     `json:"client_name"`       // 客户端名称
	CpuCompensation float64    `json:"cpu_compensation"`  // CPU性能放大倍数
	JavaTimeBonus   float64    `json:"java_time_bonus"`   // Java程序运行时间红利倍数
	JavaMemoryBonus float64    `json:"java_memory_bonus"` // Java程序运行内容红利倍数
	JavaXms         string     `json:"java_xms"`          // Java Xms 参数
	JavaXmx         string     `json:"java_xmx"`          // Java Xmx 参数
	DataPath        string     `json:"data"`              // 数据存放地址
	ClientExec      string     `json:"client"`            // justoj-core-client 绝对路径
	SecureCode      string     `json:"secure_code"`       // 服务器用于判题机认证的Secure Code，存放在 justoj 项目的 .env 文件中
	Languages       []Language `json:"languages"`         // 语言列表
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
