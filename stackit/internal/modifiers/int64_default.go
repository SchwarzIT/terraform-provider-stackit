package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m int64DefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %d", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m int64DefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%d`", m.Default)
}

// Modify runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m int64DefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	// types.Int64 must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic plan modifiers, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var i types.Int64
	diags := tfsdk.ValueAs(ctx, req.AttributePlan, &i)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if !i.Null {
		return
	}
	resp.AttributePlan = types.Int64Value(m.Default)
}

func Int64Default(defaultValue int64) int64DefaultModifier {
	return int64DefaultModifier{
		Default: defaultValue,
	}
}
