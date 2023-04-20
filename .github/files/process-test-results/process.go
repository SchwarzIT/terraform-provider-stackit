package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

	// modify aggregation because every error is printed twice
	agg.Overall.Fail = agg.Overall.Fail / 2
	agg.Overall.Pass = agg.Overall.Pass / 2
	packagesTmp := map[string]Summary{}
	for k, v := range agg.Packages {
		packagesTmp[k] = Summary{
			Fail: v.Fail / 2,
			Pass: v.Pass / 2,
		}
	}
	agg.Packages = packagesTmp

	readme := "README.md"
	errored := false
	urlBadge, badge := getBadge(agg.Overall)
	if err := injectToMarkdownFile(readme, "<!--workflow-badge-->", badge); err != nil {
		errored = true
		fmt.Println(err)
	}

	urlOverview, imageTag := generateImage(agg)
	if err := injectToMarkdownFile(readme, "<!--summary-image-->", imageTag); err != nil {
		errored = true
		fmt.Println(err)
	}
	if err := sendTeamsNotification(agg, urlBadge, urlOverview); err != nil {
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

func getBadge(s Summary) (string, string) {
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

	return fmt.Sprintf(`https://img.shields.io/badge/Acceptance%%20Tests-%s-%s`, url.PathEscape(badgeText), color),
		fmt.Sprintf(`[![GitHub Workflow Status](https://img.shields.io/badge/Acceptance%%20Tests-%s-%s)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml)`, url.PathEscape(badgeText), color)
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

func sendTeamsNotification(v TestsSummary, badgeImgURL, overviewImgURL string) error {
	if v.Overall.Fail == 0 {
		return nil
	}
	webhookURL := os.Getenv("TEAMS_WEBHOOK_URL")
	if webhookURL == "" {
		return errors.New("webhookURL is empty. please set TEAMS_WEBHOOK_URL")
	}
	githubRunID := os.Getenv("GITHUB_RUN_ID")

	text := ""
	for k, v := range v.Packages {
		if v.Fail == 0 {
			continue
		}
		text += fmt.Sprintf("• **%s**: %d failed, %d succeeded\\n\\n", k, v.Fail, v.Pass)
	}

	card := generateAdaptiveCard(githubRunID, badgeImgURL, overviewImgURL, text)
	fmt.Println(card)
	payload := strings.NewReader(card)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, webhookURL, payload)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error posting message: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", res.StatusCode)
	}

	return nil
}

func generateAdaptiveCard(runID, badgeImageURL, overviewImageURL, resourceTestsText string) string {
	return fmt.Sprintf(`{
		"type": "message",
		"attachments": [
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"contentUrl": null,
				"content": {
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"type": "AdaptiveCard",
					"msTeams": {
						"width": "full"
					},
					"version": "1.0",
					"body": [
						{
							"type": "Container",
							"id": "1c2fc7ce-39c2-670b-e12a-6fff4159f487",
							"padding": "None",
							"style": "default",
							"spacing": "Medium",
							"verticalContentAlignment": "Center",
							"width": "stretch",
							"items": [
								{
									"type": "Container",
									"id": "9c1b7408-2c73-5328-86b0-73fd57eb91a9",
									"padding": "Medium",
									"width": "stretch",
									"items": [
										{
											"type": "ColumnSet",
											"id": "21c3caa3-9b86-79e1-2fa3-3cc61351149d",
											"columns": [
												{
													"type": "Column",
													"id": "b05d4957-5e66-6185-b0ca-40a8f5da7b70",
													"padding": "Small",
													"width": "stretch",
													"items": [
														{
															"type": "TextBlock",
															"id": "4cc2f99c-131e-12e6-41a7-3934b95fc0ec",
															"text": "Acceptance Tests Results",
															"wrap": true,
															"size": "Medium",
															"weight": "Bolder"
														}
													]
												},
												{
													"type": "Column",
													"id": "c9b7d464-0c6f-45eb-e8bb-97cb0896640c",
													"padding": "Small",
													"height": "15px",
													"items": [
														{
															"type": "Image",
															"id": "b72abea5-26db-0795-e045-e6cc5009db97",
															"url": "%s",
															"selectAction": {
																"type": "Action.OpenUrl",
																"url": "https://github.com/SchwarzIT/terraform-provider-stackit/actions/runs/%s"
															},
															"height": "20px",
															"spacing": "None",
															"horizontalAlignment": "Right"
														}
													]
												}
											],
											"padding": "None"
										}
									],
									"style": "emphasis"
								},
								{
									"type": "ColumnSet",
									"id": "8486ba45-e628-45ef-df1f-7049dbe45046",
									"columns": [
										{
											"type": "Column",
											"id": "43c124d5-fe13-aed2-93d7-d647e8ba8252",
											"padding": "None",
											"width": "stretch",
											"items": [
												{
													"type": "TextBlock",
													"id": "32600619-df1a-0d14-5b61-79862b15f149",
													"text": "%s",
													"wrap": true,
													"spacing": "None"
												}
											]
										},
										{
											"type": "Column",
											"items": [
												{
													"type": "Image",
													"id": "5dee8083-3822-5a1d-e782-c21e6d9d7ba1",
													"url": "%s",
													"selectAction": {
														"type": "Action.OpenUrl",
														"url": "https://github.com/SchwarzIT/terraform-provider-stackit/actions/runs/%s"
													},
													"size": "Large",
													"width": "250px",
													"spacing": "None"
												}
											],
											"padding": "None",
											"width": "auto"
										}
									],
									"padding": "Medium",
									"spacing": "None"
								},
								{
									"type": "Container",
									"id": "4c0b33c7-8fb0-f7ff-3d5b-9a5181fb7edf",
									"padding": "Default",
									"items": [
										{
											"type": "ActionSet",
											"actions": [
												{
													"type": "Action.OpenUrl",
													"id": "4a0a75e4-b785-ad20-c61a-c18de7dfc6a9",
													"title": "View run",
													"url": "https://github.com/SchwarzIT/terraform-provider-stackit/actions/runs/%s",
													"style": "positive",
													"isPrimary": true
												}
											],
											"spacing": "Small"
										}
									],
									"spacing": "Medium",
									"separator": true
								}
							]
						}
					],
					"padding": "None"
				}
			}
		]
	}`, badgeImageURL, runID, resourceTestsText, overviewImageURL, runID, runID)
}
