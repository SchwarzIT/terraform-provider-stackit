package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Project is the schema model
type Project struct {
	ID                types.String `tfsdk:"id"`
	ContainerID       types.String `tfsdk:"container_id"`
	ParentContainerID types.String `tfsdk:"parent_container_id"`
	Name              types.String `tfsdk:"name"`
	BillingRef        types.String `tfsdk:"billing_ref"`
}

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Data source for projects",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "the project UUID",
				Type:        types.StringType,
				Computed:    true,
			},

			"container_id": {
				Description: "the project container ID",
				Type:        types.StringType,
				Required:    true,
			},

			"parent_container_id": {
				Description: "the project's parent container ID",
				Type:        types.StringType,
				Computed:    true,
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
	}
}
