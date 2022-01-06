package main

import (
	"fmt"
	"github.com/ismdeep/log"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

var solutionQueue chan string
var solutionTimeoutMap map[string]int64

func init() {
	solutionQueue = make(chan string, 300)
	solutionTimeoutMap = make(map[string]int64)
}

func SolutionPullWorker() {
	for {
		solutionIDs, err := GetPendingSolutions()
		if err != nil {
			log.Error("main", log.FieldErr(err))
			time.Sleep(1 * time.Second)
			continue
		}
		for _, id := range solutionIDs {
			if _, found := solutionTimeoutMap[id]; found {
				continue
			}
			solutionTimeoutMap[id] = time.Now().Unix() + 6
			solutionQueue <- id
		}

		time.Sleep(1 * time.Second)
	}
}

func SolutionCleanWorker() {
	for {
		for key, v := range solutionTimeoutMap {
			if v < time.Now().Unix() {
				delete(solutionTimeoutMap, key)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func SolutionJudgeWorker() {
	for {
		solutionID := <-solutionQueue
		fmt.Println(solutionID)
		// 等待判题结果

		// 1. 准备数据， ${WorkDir}/run/${solution_id}-${rand-hex}

		// 2. 执行判题 justoj-core-client -d ${WorkDir}/run/${solution_id}-${rand-hex}

		// 3. 解析结果 ${WorkDir}/run/${solution_id}-${rand-hex}/run/results.txt

		// 4. 清理目录

		fmt.Printf("solution judge done. [%v]\n", solutionID)
		delete(solutionTimeoutMap, solutionID)
	}
}

// GetCPUBench 获取 CPU_BENCHMARK 结果
func GetCPUBench() (string, error) {
	cmd := exec.Command("justoj-cpu-benchmark")
	p, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(p)
	if err != nil {
		return "", err
	}

	cpuBench := strings.TrimSpace(string(bytes))

	return cpuBench, nil
}

func main() {
	cpuBench, err := GetCPUBench()
	if err != nil {
		panic(err)
	}
	fmt.Println(cpuBench)

	go SolutionCleanWorker()
	go SolutionJudgeWorker()
	SolutionPullWorker()
}
