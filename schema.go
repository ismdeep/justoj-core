package main

type SolutionInfo struct {
	ID        string
	ProblemID string
	UserID    string
	Language  string
	ContestID string
}

type ProblemInfo struct {
	ID          string
	TimeLimit   string
	MemoryLimit string
	Spj         bool
}

type ResultInfo struct {
	Name   string `json:"name"`
	Result int    `json:"result"`
	Time   int64  `json:"time"`
	Mem    int64  `json:"mem"`
}

type SolutionResult struct {
	SolutionID     string
	RunDir         string
	Result         int
	CompileError   string
	ResultInfoList []*ResultInfo
	MaxTime        int64
	MaxMem         int64
}
