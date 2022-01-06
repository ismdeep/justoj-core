package main

import (
	"fmt"
	"github.com/ismdeep/log"
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
			return
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
		fmt.Printf("solution judge done. [%v]\n", solutionID)
		delete(solutionTimeoutMap, solutionID)
	}
}

func main() {
	fmt.Println(Config)
	go SolutionCleanWorker()
	go SolutionJudgeWorker()
	SolutionPullWorker()
}
