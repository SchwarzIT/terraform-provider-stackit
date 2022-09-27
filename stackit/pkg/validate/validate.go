package validate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type Validator struct {
	description         string
	markdownDescription string
	validate            ValidationFn
}

type ValidationFn func(context.Context, tfsdk.ValidateAttributeRequest, *tfsdk.ValidateAttributeResponse)

var _ = tfsdk.AttributeValidator(&Validator{})

func (v *Validator) Description(ctx context.Context) string {
	return v.description
}

func (v *Validator) MarkdownDescription(ctx context.Context) string {
	return v.markdownDescription
}

func (v *Validator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig.IsUnknown() || req.AttributeConfig.IsNull() {
		return
	}
	v.validate(ctx, req, res)
}
