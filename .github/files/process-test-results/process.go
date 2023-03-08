package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
)

type TestsSummary struct {
	Overall  Summary            `json:"overall"`
	Packages map[string]Summary `json:"packages"`
}

type Summary struct {
	Pass int `json:"pass"`
	Fail int `json:"fail"`
}

func main() {
	dirname := "test"
	agg := TestsSummary{
		Packages: map[string]Summary{},
	}
	dir, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	for _, file := range dir {
		var f TestsSummary
		b, err := os.ReadFile(path.Join(dirname, file.Name()))
		if err != nil {
			continue
		}
		fmt.Println("processing: " + file.Name())
		if err := json.Unmarshal(b, &f); err != nil {
			panic(err)
		}
		agg.Overall.Fail += f.Overall.Fail
		agg.Overall.Pass += f.Overall.Pass

		for k := range f.Packages {
			v1 := agg.Packages[k]
			v2 := f.Packages[k]

			v1.Fail += v2.Fail
			v1.Pass += v2.Pass
			agg.Packages[k] = v1
		}
	}

}

func getBadge(s Summary) string {
	badgeText := ""
	switch true {
	case s.Fail > 0 && s.Pass > 0:
		badgeText = fmt.Sprintf("%d passed, %d failed", s.Pass, s.Fail)
	case s.Pass > 0:
		badgeText = fmt.Sprintf("%d passed", s.Pass)
	case s.Fail > 0:
		badgeText = fmt.Sprintf("%d failed", s.Fail)
	}

	color := "red"
	colorIndicator := s.Pass / (s.Pass + s.Fail)
	if colorIndicator > 0.5 {
		color = "orange"
	}
	if colorIndicator > 0.7 {
		color = "yellow"
	}
	if colorIndicator > 0.8 {
		color = "green"
	}
	if colorIndicator > 0.9 {
		color = "success"
	}
	return fmt.Sprintf("https://img.shields.io/badge/Acceptance%%20Tests-%s-%s", url.QueryEscape(badgeText), color)
}
