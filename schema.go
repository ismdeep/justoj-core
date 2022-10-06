package main

// SolutionInfo struct
type SolutionInfo struct {
	ID        string
	ProblemID string
	UserID    string
	Language  string
	ContestID string
}

// ProblemInfo struct
type ProblemInfo struct {
	ID          string
	TimeLimit   string
	MemoryLimit string
	Spj         bool
}

// ResultInfo struct
type ResultInfo struct {
	Name   string `json:"name"`
	Result int    `json:"result"`
	Time   int64  `json:"time"`
	Mem    int64  `json:"mem"`
}

// SolutionResult struct
type SolutionResult struct {
	SolutionID     string
	RunDir         string
	Result         int
	CompileError   string
	ResultInfoList []*ResultInfo
	MaxTime        int64
	MaxMem         int64
}
