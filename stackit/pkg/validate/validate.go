package validate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Validator struct {
	description         string
	markdownDescription string
	validate            ValidationFn
	validateMap         ValidationMapFn
	validateInt         ValidationIntFn
	validateList        ValidationListFn
}

type ValidationFn func(context.Context, validator.StringRequest, *validator.StringResponse)

type ValidationMapFn func(context.Context, validator.MapRequest, *validator.MapResponse)

type ValidationIntFn func(context.Context, validator.Int64Request, *validator.Int64Response)

type ValidationListFn func(context.Context, validator.ListRequest, *validator.ListResponse)

var _ = validator.String(&Validator{})

func (v *Validator) Description(ctx context.Context) string {
	return v.description
}

func (v *Validator) MarkdownDescription(ctx context.Context) string {
	return v.markdownDescription
}

func (v *Validator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	v.validate(ctx, req, resp)
}

func (v *Validator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	v.validateMap(ctx, req, resp)
}

func (v *Validator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	v.validateInt(ctx, req, resp)
}

func (v *Validator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	v.validateList(ctx, req, resp)
}
