package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
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
	errored := false
	if err := injectToMarkdownFile(readme, "<!--workflow-badge-->", getBadge(agg.Overall)); err != nil {
		errored = true
		fmt.Println(err)
	}
	url, imageTag := generateImage(agg)
	if err := injectToMarkdownFile(readme, "<!--summary-image-->", imageTag); err != nil {
		errored = true
		fmt.Println(err)
	}
	if err := sendTeamsNotification(agg, url); err != nil {
		errored = true
		fmt.Println(err)
	}

	b, _ := json.MarshalIndent(agg, "", "  ")
	fmt.Println(string(b))

	if errored {
		panic("an error occured")
	}
}

func injectToMarkdownFile(name, separator, injected string) error {
	b, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	rslice := strings.Split(string(b), separator)
	if len(rslice) == 3 {
		rslice[1] = injected + fmt.Sprintf("<!--revision-%s-->", uuid.New().String())
	}
	if err := os.WriteFile(name, []byte(strings.Join(rslice, separator)), 644); err != nil {
		return err
	}
	return nil
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

func generateImage(v TestsSummary) (string, string) {
	img, err := callAPI(generateHTML(v))
	if err != nil {
		fmt.Println(err.Error())
	}
	if img != "" {
		return img, fmt.Sprintf(`
<img src="%s" width="250" align="right" />
`, img)
	}
	return "", ""
}

const css = `<style type="text/css">
.tg  {border-collapse:collapse;border: none;margin-bottom:20px;}
.tg td { padding: 2px 5px; border: none; font-size: 12px; font-family: 'courier' }
</style>`

const table = `<table class="tg">
<tbody><tr>
%s
</tr></tbody>
</table>
`

const td = `<td class="tg-0lax">%s</td><td class="tg-0lax">%s</td>`

func generateHTML(v TestsSummary) string {
	md := ""
	i := 0
	for pkg, sum := range v.Packages {
		md += generateRow(pkg, sum)
		if i++; i%2 == 0 {
			md += "<tr></tr>"
		}
	}
	return fmt.Sprintf(table, md)
}

func generateRow(pkg string, sum Summary) string {
	return fmt.Sprintf(td, getIcon(sum), pkg)
}

func getIcon(sum Summary) string {
	color := "🔥"
	var pc float32
	s := float32(sum.Pass + sum.Fail)
	if s == 0 {
		s = 1
	}
	pc = float32(sum.Pass) / s
	if pc > 0.7 {
		color = "⚠️"
	}
	if pc == 1 {
		color = "🟢"
	}
	return color
}

func callAPI(html string) (string, error) {
	data := map[string]string{
		"html": html,
		"css":  css,
	}
	reqBody, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://hcti.io/v1/image", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	userID := os.Getenv("HCTI_USER_ID")
	apiKey := os.Getenv("HCTI_API_KEY")
	req.SetBasicAuth(userID, apiKey)
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var v struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return "", err
	}
	return v.URL, nil
}

type TeamsMessage struct {
	Type            string         `json:"@type"`
	Context         string         `json:"@context"`
	ThemeColor      string         `json:"themeColor,omitempty"`
	Summary         string         `json:"summary,omitempty"`
	Title           string         `json:"title,omitempty"`
	Text            string         `json:"text,omitempty"`
	Sections        []TeamsSection `json:"sections,omitempty"`
	Markdown        bool           `json:"markdown"`
	PotentialAction []TeamsAction  `json:"potentialAction,omitempty"`
}

type TeamsSection struct {
	ActivityTitle    string      `json:"activityTitle,omitempty"`
	ActivitySubtitle string      `json:"activitySubtitle,omitempty"`
	ActivityImage    string      `json:"activityImage,omitempty"`
	Facts            []TeamsFact `json:"facts,omitempty"`
	Markdown         bool        `json:"markdown"`
}

type TeamsAction struct {
	Type   string `json:"@type"`
	Name   string `json:"name"`
	Target string `json:"target,omitempty"`
}

func sendTeamsNotification(v TestsSummary, imgURL string) error {
	if v.Overall.Fail == 0 {
		return nil
	}
	webhookURL := os.Getenv("ODJ_STACKIT_GENERAL_CHANNEL_WEBHOOK_URL")
	section := TeamsSection{
		ActivityTitle: "Overview",
		ActivityImage: imgURL,
		Markdown:      false,
	}

	githubServerURL := os.Getenv("GITHUB_SERVER_URL")
	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	githubRunID := os.Getenv("GITHUB_RUN_ID")
	workflowRunURL := fmt.Sprintf("%s/%s/actions/runs/%s", githubServerURL, githubRepository, githubRunID)

	action := TeamsAction{
		Type:   "OpenUri",
		Name:   "View Run",
		Target: workflowRunURL,
	}

	text := ""
	for k, v := range v.Packages {
		if v.Fail == 0 {
			continue
		}
		text += fmt.Sprintf("**%s**\n%d failed, %d succeeded\n\n", k, v.Fail, v.Pass)
	}

	msg := TeamsMessage{
		Type:            "MessageCard",
		Context:         "http://schema.org/extensions",
		ThemeColor:      "0078D7",
		Summary:         "Acceptance Tests failur summary",
		Title:           "Acceptance Tests failures",
		Text:            text,
		Markdown:        true,
		Sections:        []TeamsSection{section},
		PotentialAction: []TeamsAction{action},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error posting message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
}
