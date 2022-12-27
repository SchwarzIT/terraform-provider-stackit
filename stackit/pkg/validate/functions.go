package validate

import (
	"context"

	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func StringWith(fn func(string) error, description string) *Validator {
	return &Validator{
		description: description,
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := fn(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func ProjectName() *Validator {
	return &Validator{
		description: "validate project name",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := clientValidate.ProjectName(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func ProjectID() *Validator {
	return &Validator{
		description: "validate project ID",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := clientValidate.ProjectID(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func BillingRef() *Validator {
	return &Validator{
		description: "validate billing reference",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := clientValidate.BillingRef(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func UUID() *Validator {
	return &Validator{
		description: "validate string is UUID",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := clientValidate.UUID(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
				return
			}
		},
	}
}
