package project

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Project is the schema model
type Project struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	BillingRef          types.String `tfsdk:"billing_ref"`
	Owner               types.String `tfsdk:"owner"`
	EnableKubernetes    types.Bool   `tfsdk:"enable_kubernetes"`
	EnableObjectStorage types.Bool   `tfsdk:"enable_object_storage"`
}

// GetSchema returns the terraform schema structure
func (r Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages projects",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "the project ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
			},

			"name": {
				Description: "the project name",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectName(),
				},
			},

			"billing_ref": {
				Description: "billing reference for cost transparency",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.BillingRef(),
				},
			},

			"owner": {
				Description: "user ID of the owner of the project. This value is only considered during creation. changing it afterwards will have no effect.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.UUID(),
				},
			},

			"enable_kubernetes": {
				Description: "should kubernetes be enabled for this project? `false` by default.",
				Type:        types.BoolType,
				Optional:    true,
			},

			"enable_object_storage": {
				Description: "should object storage be enabled for this project? `false` by default.",
				Type:        types.BoolType,
				Optional:    true,
			},
		},
	}, nil
}
