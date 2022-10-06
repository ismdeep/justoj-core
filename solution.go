package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GenerateURL generate url
func GenerateURL(uri string, params map[string]interface{}) string {
	l := make([]string, 0)
	l = append(l, fmt.Sprintf("client_name=%v", ClientName))
	l = append(l, fmt.Sprintf("secure_code=%v", Config.SecureCode))
	for key, value := range params {
		l = append(l, fmt.Sprintf("%v=%v", key, value))
	}
	return fmt.Sprintf("%v%v?%v", Config.BaseURL, uri, strings.Join(l, "&"))
}

// GetPendingSolutions get pending solutions
func GetPendingSolutions() ([]string, error) {
	ids := GetAvailableLanguageIDs()
	if len(ids) <= 0 {
		return nil, errors.New("no language supported")
	}

	h := &http.Client{}
	resp, err := h.Get(GenerateURL("/api/judge_api/get_pending", map[string]interface{}{
		"query_size":  10,
		"oj_lang_set": strings.Join(ids, ","),
	}))
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

// GetSolutionInfo get solution info
func GetSolutionInfo(solutionID string) (*SolutionInfo, error) {
	h := &http.Client{}
	resp, err := h.Get(GenerateURL("/api/judge_api/get_solution_info", map[string]interface{}{
		"sid": solutionID,
	}))
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

// GetProblemInfo get problem info
func GetProblemInfo(problemID string) (*ProblemInfo, error) {
	h := &http.Client{}
	resp, err := h.Get(GenerateURL("/api/judge_api/get_problem_info", map[string]interface{}{
		"pid": problemID,
	}))
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

// GetSolutionSourceCode get source code
func GetSolutionSourceCode(solutionID string) (string, error) {
	h := &http.Client{}
	resp, err := h.Get(GenerateURL("/api/judge_api/get_solution", map[string]interface{}{
		"sid": solutionID,
	}))
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// UpdateSolutionResult update solution result
func UpdateSolutionResult(res *SolutionResult) error {
	if _, err := (&http.Client{}).Get(GenerateURL("/api/judge_api/update_solution", map[string]interface{}{
		"sid":    res.SolutionID,
		"result": res.Result,
		"time":   res.MaxTime,
		"memory": res.MaxMem / 1024,
	})); err != nil {
		return err
	}

	if res.Result == ResultCompileError {
		data := struct {
			CeInfo string `json:"ceinfo"`
		}{
			CeInfo: res.CompileError,
		}
		content, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if _, err := (&http.Client{}).Post(GenerateURL("/api/judge_api/add_ce_info", map[string]interface{}{
			"sid": res.SolutionID,
		}),
			"application/json",
			bytes.NewReader(content)); err != nil {
			return err
		}
	}

	return nil
}
