package common

import (
	"context"
	"fmt"
	"os"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/urls"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ToString converts attr.Value to string
func ToString(ctx context.Context, v attr.Value) (string, error) {
	if t := v.Type(ctx); t != types.StringType {
		return "", fmt.Errorf("type mismatch. expected 'types.StringType' but got '%s'", t.String())
	}
	if v.IsNull() || v.IsUnknown() {
		return "", fmt.Errorf("value is unknown or null")
	}
	tv, err := v.ToTerraformValue(ctx)
	if err != nil {
		return "", err
	}
	var s string
	if err := tv.Copy().As(&s); err != nil {
		return "", err
	}
	return s, nil
}

// ShouldAccTestRun returns true of the provided flag is true or if
// an env variable ACC_TEST_CI has any value
func ShouldAccTestRun(runFlag bool) bool {
	if v, ok := os.LookupEnv("ACC_TEST_CI"); (ok && v != "") || runFlag {
		return true
	}
	return false
}

// GetAcceptanceTestsProjectID returns the project ID for acceptance test
// can be overridden by setting ACC_TEST_PROJECT_ID
func GetAcceptanceTestsProjectID() string {
	if v, ok := os.LookupEnv("ACC_TEST_PROJECT_ID"); ok && v != "" {
		return v
	}
	return ACC_TEST_PROJECT_ID
}

func EnvironmentInfo(u urls.ByEnvs) string {
	return fmt.Sprintf(`<div class="warning" style='color: #69337A; border: solid #E9D8FD 1px; border-radius: 4px; padding-left:0.7em;margin-top:5px;'>
<span>
<p style='margin-top:1em;'>
<b>Environment support</b>
<table style='border-collapse: separate; margin:0;'>
<tr><td style='width: 100px'>Production</td><td>%s<td></tr>
<tr><td>QA</td><td>%s<td></tr>
<tr><td>Dev</td><td>%s<td></tr>
</table>
<br />
<small>By default, <a href="https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs#environment">production</a> is used.<br />To set a custom URL, set an environment variable %s</small>
</p>
</span>
</div>`,
		u.Prod, u.QA, u.Dev, u.OverrideWith,
	)
}
