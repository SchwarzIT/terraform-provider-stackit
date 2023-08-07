package common

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/baseurl"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
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

func EnvironmentInfo(u baseurl.BaseURL) string {
	return fmt.Sprintf(`
<br />

-> __Environment support__<small>To set a custom API base URL, set <code>%s</code> environment variable </small>
	`,
		u.OverrideWith,
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

func GetDefaultACL() defaults.List {
	av := []attr.Value{}
	for _, r := range KnownRanges {
		av = append(av, types.StringValue(r))
	}
	return listdefault.StaticValue(types.ListValueMust(types.StringType, av))
}

func Validate(d *diag.Diagnostics, res interface{}, err error, checkNullFields ...string) error {
	agg := validate.Response(res, err, checkNullFields...)
	if agg != nil && res != nil {
		resValue := reflect.ValueOf(res)
		if resValue.Kind() == reflect.Ptr && resValue.Elem().Kind() == reflect.Struct {
			resValue = resValue.Elem()
		}
		if resValue.Kind() != reflect.Struct {
			return agg
		}
		body := resValue.FieldByName("Body")
		if body.IsValid() && !body.IsNil() {
			if b, ok := body.Interface().([]byte); ok {
				Dump(d, b)
			}
		}
	}
	return agg
}
