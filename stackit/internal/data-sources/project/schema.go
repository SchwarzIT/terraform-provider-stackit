package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Project is the schema model
type Project struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	BillingRef types.String `tfsdk:"billing_ref"`
}

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Project data source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "the project ID",
				Type:        types.StringType,
				Required:    true,
			},

			"name": {
				Description: "the project name",
				Type:        types.StringType,
				Computed:    true,
			},

			"billing_ref": {
				Description: "billing reference for cost transparency",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}
