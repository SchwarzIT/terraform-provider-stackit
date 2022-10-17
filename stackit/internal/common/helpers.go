package common

import (
	"context"
	"fmt"
	"os"

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
