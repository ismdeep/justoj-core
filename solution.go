package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetPendingSolutions() ([]string, error) {
	h := &http.Client{}
	resp, err := h.Get(fmt.Sprintf("%v/api/judge_api/get_pending?query_size=10&secure_code=%v&oj_lang_set=1,2,3", Config.BaseUrl, Config.SecureCode))
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
