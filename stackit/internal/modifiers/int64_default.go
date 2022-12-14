package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// int64DefaultModifier is a plan modifier that sets a default value for a
// types.Int64Type attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type int64DefaultModifier struct {
	Default int64
}

var _ planmodifier.Int64 = (*int64DefaultModifier)(nil)

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m int64DefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %d", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m int64DefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%d`", m.Default)
}

// PlanModifyBool runs the logic of the plan modifier.
func (m int64DefaultModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	resp.PlanValue = types.Int64Value(m.Default)
}

func Int64Default(defaultValue int64) int64DefaultModifier {
	return int64DefaultModifier{
		Default: defaultValue,
	}
}
