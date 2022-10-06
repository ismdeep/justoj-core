package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ismdeep/log"
	"github.com/ismdeep/rand"
)

var solutionQueue chan string
var solutionTimeoutMap map[string]int64
var cpuCompensation string

func init() {
	solutionQueue = make(chan string, 300)
	solutionTimeoutMap = make(map[string]int64)
}

// SolutionPullWorker worker
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

// SolutionCleanWorker clean worker
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

// SolutionJudge judge
func SolutionJudge(solutionID string) (*SolutionResult, error) {
	runHexID := rand.HexStr(32)
	runDir := fmt.Sprintf("%v/run/%v-%v", WorkDir, solutionID, runHexID)
	if err := os.MkdirAll(runDir, 0777); err != nil {
		return nil, err
	}
	defer func() {
		// 清理目录
		go func() {
			_ = exec.Command("rm", "-rf", runDir).Start()
		}()
	}()
	// 1. 准备数据 ${WorkDir}/run/${solution_id}-${rand-hex}
	solutionInfo, err := GetSolutionInfo(solutionID)
	if err != nil {
		return nil, err
	}
	problemInfo, err := GetProblemInfo(solutionInfo.ProblemID)
	if err != nil {
		return nil, err
	}
	spjV := 0
	if problemInfo.Spj {
		spjV = 1
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%v/config", runDir),
		[]byte(fmt.Sprintf("LANGUAGE=%v\nSPECIAL_JUDGE=%v\nCPU_COMPENSATION=%v\nTIME_LIMIT=%v\nMEMORY_LIMIT=%v",
			solutionInfo.Language,    // LANGUAGE
			spjV,                     // SPECIAL_JUDGE
			cpuCompensation,          // CPU_COMPENSATION
			problemInfo.TimeLimit,    // TIME_LIMIT
			problemInfo.MemoryLimit), // MEMORY_LIMIT
		), 0777); err != nil {
		return nil, err
	}

	cmdCopy := exec.Command(
		"cp",
		"-r",
		"-v",
		fmt.Sprintf("%v/justoj-data/data/%v", WorkDir, problemInfo.ID),
		fmt.Sprintf("%v/data", runDir))
	if err := cmdCopy.Start(); err != nil {
		return nil, err
	}
	if err := cmdCopy.Wait(); err != nil {
		return nil, err
	}

	src, err := GetSolutionSourceCode(solutionID)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%v/code", runDir), []byte(src), 0777); err != nil {
		return nil, err
	}

	// 2. 执行判题 justoj-core-client -d ${WorkDir}/run/${solution_id}-${rand-hex}

	cmd := exec.Command("justoj-core-client", "-d", runDir)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	// 3. 解析结果 ${WorkDir}/run/${solution_id}-${rand-hex}/run/results.txt
	compileErr, err := ioutil.ReadFile(fmt.Sprintf("%v/run/ce.txt", runDir))
	if err != nil {
		return nil, err
	}
	if string(compileErr) != "" {
		return &SolutionResult{
			SolutionID:     solutionID,
			RunDir:         runDir,
			Result:         ResultCompileError,
			CompileError:   string(compileErr),
			ResultInfoList: nil,
			MaxTime:        0,
			MaxMem:         0,
		}, nil
	}

	content, err := ioutil.ReadFile(fmt.Sprintf("%v/run/results.txt", runDir))
	if err != nil {
		return nil, err
	}
	results := make([]*ResultInfo, 0)
	if err := json.Unmarshal(content, &results); err != nil {
		return nil, err
	}

	if len(results) <= 0 {
		return &SolutionResult{
			SolutionID:     solutionID,
			RunDir:         runDir,
			Result:         ResultAccepted,
			CompileError:   "",
			ResultInfoList: nil,
			MaxTime:        0,
			MaxMem:         0,
		}, nil
	}

	result := ResultAccepted
	maxTime := int64(0)
	maxMem := int64(0)
	for _, v := range results {
		if v.Time > maxTime {
			maxTime = v.Time
		}
		if v.Mem > maxMem {
			maxMem = v.Mem
		}
		if v.Result != ResultAccepted {
			result = v.Result
			break
		}
	}

	return &SolutionResult{
		SolutionID:     solutionID,
		RunDir:         runDir,
		Result:         result,
		CompileError:   "",
		ResultInfoList: results,
		MaxTime:        maxTime,
		MaxMem:         maxMem,
	}, nil
}

// SolutionJudgeWorker worker
func SolutionJudgeWorker() {
	for {
		solutionID := <-solutionQueue
		result, err := SolutionJudge(solutionID)
		if err != nil {
			log.Error("SolutionJudgeWorker", log.FieldErr(err))
			continue
		}

		fmt.Printf("Solution: %v    Result: %18v    Time: %6v    Mem: %10v    @    %v\n", result.SolutionID, ResultText[result.Result], result.MaxTime, result.MaxMem, result.RunDir)
		if err := UpdateSolutionResult(result); err != nil {
			log.Error("SolutionJudgeWorker", log.FieldErr(err))
		}
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
	var err error
	cpuCompensation, err = GetCPUBench()
	if err != nil {
		panic(err)
	}

	go SolutionCleanWorker()
	go SolutionJudgeWorker()
	SolutionPullWorker()
}
