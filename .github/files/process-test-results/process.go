package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
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
	dirname := ".github/files/process-test-results/artifact"
	agg := TestsSummary{
		Packages: map[string]Summary{},
	}
	dir, _ := os.ReadDir(dirname)
	for _, file := range dir {
		var f TestsSummary
		b, err := os.ReadFile(path.Join(dirname, file.Name()))
		if err != nil {
			continue
		}
		fmt.Println("processing: " + file.Name())
		if err := json.Unmarshal(b, &f); err != nil {
			continue
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

	readme := "README.md"
	b, err := os.ReadFile(readme)
	if err != nil {
		panic(err)
	}
	rslice := strings.Split(string(b), "<!--workflow-badge-->")
	if len(rslice) == 3 {
		rslice[1] = getBadge(agg.Overall)
	}
	os.WriteFile(readme, []byte(strings.Join(rslice, "<!--workflow-badge-->")), 644)
	fmt.Println(getBadge(agg.Overall))
}

func getBadge(s Summary) string {
	badgeText := "All failed"
	switch true {
	case s.Fail == 0 && s.Pass > 0:
		badgeText = "All passed"
	case s.Fail > 0 && s.Pass > 0:
		badgeText = fmt.Sprintf("%d passed, %d failed", s.Pass, s.Fail)
	}

	color := "red"
	var colorIndicator float32
	sum := float32(s.Pass + s.Fail)
	if sum == 0 {
		sum = 1
	}
	colorIndicator = float32(s.Pass) / sum
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

	return fmt.Sprintf(`[![GitHub Workflow Status](https://img.shields.io/badge/Acceptance%%20Tests-%s-%s)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml)`, url.PathEscape(badgeText), color)
}
