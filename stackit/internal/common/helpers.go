package common

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

func IsNonRetryable(err error) bool {
	if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
		if !strings.Contains(err.Error(), ERR_NO_SUCH_HOST) {
			return true
		}
	}
	if strings.Contains(err.Error(), http.StatusText(http.StatusUnauthorized)) {
		return true
	}
	return false
}
