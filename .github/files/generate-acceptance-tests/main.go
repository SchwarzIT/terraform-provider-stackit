package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	rk, dsk := []string{}, []string{}
	sr, sds := map[string]string{}, map[string]string{}
	globalKeysRes := map[string]interface{}{}
	globalKeysDS := map[string]interface{}{}
	err := filepath.Walk("stackit/internal/",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(path, "stackit/internal/data-sources") && strings.HasSuffix(path, "_test.go") {
				sl := strings.Split(path, "/")
				key := strings.Join(sl[3:len(sl)-1], " ")
				if _, ok := sds[key]; ok {
					return nil
				}
				globalKeysRes[sl[3]] = nil
				dsk = append(dsk, key)
				sds[key] = strings.Join(sl[:len(sl)-1], "/")
			}
			if strings.HasPrefix(path, "stackit/internal/resources") && strings.HasSuffix(path, "_test.go") {
				sl := strings.Split(path, "/")
				key := strings.Join(sl[3:len(sl)-1], " ")
				if _, ok := sr[key]; ok {
					return nil
				}
				globalKeysDS[sl[3]] = nil
				rk = append(rk, key)
				sr[key] = strings.Join(sl[:len(sl)-1], "/")
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	sortedGlobalKeysRes := []string{}
	for g := range globalKeysRes {
		sortedGlobalKeysRes = append(sortedGlobalKeysRes, g)
	}

	sortedGlobalKeysDS := []string{}
	for g := range globalKeysDS {
		sortedGlobalKeysDS = append(sortedGlobalKeysDS, g)
	}

	sort.Strings(sortedGlobalKeysRes)
	sort.Strings(sortedGlobalKeysDS)
	sort.Strings(rk)
	sort.Strings(dsk)

	// fmt.Println("found resources:")
	// printOutcome(sortedGlobalKeysRes, rk, sr)

	// fmt.Println("\nfound data sources:")
	// printOutcome(sortedGlobalKeysDS, dsk, sds)

	s := "# this is a generated file, DO NOT EDIT\n# to generate this file run make pre-commit\n"
	data, err := ioutil.ReadFile(".github/files/generate-acceptance-tests/template.yaml")
	if err != nil {
		fmt.Println(err)
	}
	sData := string(data)

	sData = strings.Replace(sData, "__data_sources_names__", "          - "+strings.Join(dsk, "\n          - "), 1)
	sData = strings.Replace(sData, "__data_sources_include__", printOutcome(sortedGlobalKeysDS, dsk, sds), 1)

	err = ioutil.WriteFile(".github/workflows/acceptance_test.yml", []byte(s+sData), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func printOutcome(sortedglobalKeys []string, sortedKeys []string, keyAndPathMap map[string]string) string {
	s := ""
	for _, key := range sortedglobalKeys {
		for _, v := range sortedKeys {
			if !strings.HasPrefix(v, key) {
				continue
			}
			s = s + fmt.Sprintf("          - name: %s\n            path: %s\n            concurrency: %s\n", v, keyAndPathMap[v], key)
		}
	}
	return s
}
