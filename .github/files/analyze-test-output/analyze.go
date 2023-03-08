package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TestsSummary struct {
	Overall  Summary            `json:"overall"`
	Packages map[string]Summary `json:"packages"`
}

type Summary struct {
	Pass int `json:"pass"`
	Fail int `json:"fail"`
}

type TestEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func main() {
	var ts TestsSummary
	ts.Packages = map[string]Summary{}

	file, err := os.Open("testoutput")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event TestEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			panic(err)
		}
		fmt.Print(event.Output)
		sp := strings.Split(event.Package, "/")
		if len(sp) < 7 {
			continue
		}
		name := sp[6]
		if name == "data-services" {
			name = getDataServiceName(event.Test)
			event.Action = getDataServiceEventAction(event.Output)
			if name == "" {
				continue
			}
		}
		v := ts.Packages[name]
		switch event.Action {
		case "fail":
			ts.Overall.Fail++
			v.Fail++
		case "pass":
			ts.Overall.Pass++
			v.Pass++
		}
		ts.Packages[name] = v
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(ts)
	if err := os.WriteFile("result/"+uuid.New().String()+".json", b, 0644); err != nil {
		panic(err)
	}
	if ts.Overall.Fail > 0 {
		panic("test failed.")
	}
}

func getDataServiceName(testName string) string {
	if strings.Contains(testName, "ElasticSearch") {
		return "elasticsearch-dsa"
	}
	if strings.Contains(testName, "LogMe") {
		return "logme-dsa"
	}
	if strings.Contains(testName, "MariaDB") {
		return "mariadb-dsa"
	}
	if strings.Contains(testName, "Postgres") {
		return "postgres-dsa"
	}
	if strings.Contains(testName, "RabbitMQ") {
		return "rabbitmq-dsa"
	}
	if strings.Contains(testName, "Redis") {
		return "redis-dsa"
	}
	return testName
}

func getDataServiceEventAction(output string) string {
	if strings.Contains(output, "PASS") {
		return "pass"
	}
	if strings.Contains(output, "SKIP") {
		return "pass"
	}
	if strings.Contains(output, "FAIL") {
		return "fail"
	}
	return ""
}
