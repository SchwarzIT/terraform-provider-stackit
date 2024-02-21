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

func ReserveProjectLabels() *Validator {
	return &Validator{
		description: "reserve project labels",
		validateMap: func(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
			for k := range req.ConfigValue.Elements() {
				// do not allow internal / hidden ones which are directly resolved by attributes
				// this is to ensure backwards compatibility
				if k == "billingReference" || k == "scope" {
					resp.Diagnostics.AddError("Reserved Project Labels", "billingReference and scope are reserved project labels")
					return
				}
			}
		},
	}
}

func NetworkID() *Validator {
	return &Validator{
		description: "validate project ID",
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

func NetworkName() *Validator {
	return &Validator{
		description: "validate network name",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := clientValidate.NetworkName(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func NameServers() *Validator {
	return &Validator{
		description: "validate name servers",
		validateMap: func(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
			for _, v := range req.ConfigValue.Elements() {
				if err := clientValidate.NameServer(v); err != nil {
					resp.Diagnostics.AddError(err.Error(), err.Error())
				}
			}
		},
	}
}

func Prefixes() *Validator {
	return &Validator{
		description: "validate prefixes",
		validateMap: func(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
			for _, v := range req.ConfigValue.Elements() {
				if err := clientValidate.Prefixes(v); err != nil {
					resp.Diagnostics.AddError(err.Error(), err.Error())
				}
			}
		},
	}
}

func PrefixLenghtV4() *Validator {
	return &Validator{
		description: "validate prefix length",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := clientValidate.PrefixLengthV4(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func PublicIP() *Validator {
	return &Validator{
		description: "validate public IP",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			v, diag := req.ConfigValue.ToStringValue(ctx)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			if err := clientValidate.PublicIp(v.ValueString()); err != nil {
				resp.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}