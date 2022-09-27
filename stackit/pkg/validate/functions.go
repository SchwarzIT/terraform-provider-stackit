package validate

import (
	"context"

	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func StringWith(fn func(string) error, description string) *Validator {
	return &Validator{
		description: description,
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
			v, err := common.ToString(ctx, req.AttributeConfig)
			if err != nil {
				res.Diagnostics.AddError("failed to get string attribute", err.Error())
				return
			}
			if err := fn(v); err != nil {
				res.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func ProjectName() *Validator {
	return &Validator{
		description: "validate project name",
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
			v, err := common.ToString(ctx, req.AttributeConfig)
			if err != nil {
				res.Diagnostics.AddError("failed to get string attribute", err.Error())
				return
			}
			if err := clientValidate.ProjectName(v); err != nil {
				res.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func ProjectID() *Validator {
	return &Validator{
		description: "validate project ID",
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
			v, err := common.ToString(ctx, req.AttributeConfig)
			if err != nil {
				res.Diagnostics.AddError("failed to get string attribute", err.Error())
				return
			}
			if err := clientValidate.ProjectID(v); err != nil {
				res.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func BillingRef() *Validator {
	return &Validator{
		description: "validate billing reference",
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
			v, err := common.ToString(ctx, req.AttributeConfig)
			if err != nil {
				res.Diagnostics.AddError("failed to get string attribute", err.Error())
				return
			}
			if err := clientValidate.BillingRef(v); err != nil {
				res.Diagnostics.AddError(err.Error(), err.Error())
			}
		},
	}
}

func UUID() *Validator {
	return &Validator{
		description: "validate string is UUID",
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
			v, err := common.ToString(ctx, req.AttributeConfig)
			if err != nil {
				res.Diagnostics.AddError("failed to get string attribute", err.Error())
				return
			}
			if err := clientValidate.UUID(v); err != nil {
				res.Diagnostics.AddError(err.Error(), err.Error())
				return
			}
		},
	}
}
