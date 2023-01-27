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
	return fmt.Sprintf(`
<br />

-> __Environment support__<br /><table style='border-collapse: separate; border-spacing: 5px; margin-top:-20px; margin-left: 24px; font-size: smaller;'>
<tr><td style='width: 100px'>Production</td><td>%s<td></tr>
<tr><td>QA</td><td>%s<td></tr>
<tr><td>Dev</td><td>%s<td></tr>
</table><br />
<small style='margin-left: 24px; margin-top: -5px; display: inline-block;'><a href="https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs#environment">By default</a>, production is used.<br />To set a custom URL, set an environment variable <code>%s</code></small>
	`,
		u.Prod, u.QA, u.Dev, u.OverrideWith,
	)
}
