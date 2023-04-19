package common

import (
	"context"
	"fmt"
	"os"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/env"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
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

func EnvironmentInfo(u env.EnvironmentURLs) string {
	return fmt.Sprintf(`
<br />

-> __Environment support__<br /><table style='border-collapse: separate; border-spacing: 0px; margin-top:-20px; margin-left: 24px; font-size: smaller;'>
<tr><td style='width: 100px; background: #fbfcff; border: none;'>Production</td><td style='background: #fbfcff; border: none;'>%s</td></tr>
<tr><td style='background: #fbfcff; border: none;'>QA</td><td style='background: #fbfcff; border: none;'>%s</td></tr>
<tr><td style='background: #fbfcff; border: none;'>Dev</td><td style='background: #fbfcff; border: none;'>%s</td></tr>
</table><br />
<small style='margin-left: 24px; margin-top: -5px; display: inline-block;'><a href="https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs#environment">By default</a>, production is used.<br />To set a custom URL, set an environment variable <code>%s</code></small>
	`,
		u.Prod, u.QA, u.Dev, u.OverrideWith,
	)
}

func Dump(d *diag.Diagnostics, body []byte) {
	d.AddWarning("request body", string(body))
}

func Timeouts(ctx context.Context, opts timeouts.Opts) schema.SingleNestedAttribute {
	timeout := timeouts.Attributes(ctx, opts).(schema.SingleNestedAttribute)
	attr := map[string]attr.Type{}
	if opts.Create {
		attr["create"] = types.StringType
	}
	if opts.Read {
		attr["read"] = types.StringType
	}
	if opts.Update {
		attr["update"] = types.StringType
	}
	if opts.Delete {
		attr["delete"] = types.StringType
	}
	timeout.Computed = true
	timeout.Default = objectdefault.StaticValue(
		types.ObjectNull(attr),
	)
	return timeout
}
