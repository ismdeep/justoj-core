package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetPendingSolutions() ([]string, error) {
	ids := GetAvailableLanguageIDs()
	if len(ids) <= 0 {
		return nil, errors.New("no language supported")
	}

	h := &http.Client{}
	resp, err := h.Get(fmt.Sprintf("%v/api/judge_api/get_pending?query_size=10&secure_code=%v&oj_lang_set=%v", Config.BaseURL, Config.SecureCode, strings.Join(ids, ",")))
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	if len(lines) <= 0 {
		return nil, errors.New("system error")
	}

	if lines[0] != "solution_ids" {
		return nil, errors.New("unauthorized")
	}

	results := make([]string, 0)
	for _, s := range lines[1:] {
		if strings.TrimSpace(s) != "" {
			results = append(results, s)
		}
	}

	return results, nil
}

func GetSolutionInfo(solutionID string) (*SolutionInfo, error) {
	h := &http.Client{}
	resp, err := h.Get(fmt.Sprintf("%v/api/judge_api/get_solution_info?sid=%v&secure_code=%v", Config.BaseURL, solutionID, Config.SecureCode))
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) < 4 {
		return nil, errors.New("unauthorized")
	}
	return &SolutionInfo{
		ID:        solutionID,
		ProblemID: lines[0],
		UserID:    lines[1],
		Language:  lines[2],
		ContestID: lines[3],
	}, nil
}

func GetProblemInfo(problemID string) (*ProblemInfo, error) {
	h := &http.Client{}
	resp, err := h.Get(fmt.Sprintf("%v/api/judge_api/get_problem_info?secure_code=%v&pid=%v", Config.BaseURL, Config.SecureCode, problemID))
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) < 3 {
		return nil, errors.New("unauthorized")
	}

	return &ProblemInfo{
		ID:          problemID,
		TimeLimit:   lines[0],
		MemoryLimit: lines[1],
		Spj:         lines[2] == "1",
	}, nil
}

func GetSolutionSourceCode(solutionID string) (string, error) {
	h := &http.Client{}
	resp, err := h.Get(fmt.Sprintf("%v/api/judge_api/get_solution?secure_code=%v&sid=%v", Config.BaseURL, Config.SecureCode, solutionID))
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
