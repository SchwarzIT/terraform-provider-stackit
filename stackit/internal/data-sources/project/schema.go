package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for projects",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "the project UUID",
				Computed:    true,
			},

			"container_id": schema.StringAttribute{
				Description: "the project container ID",
				Required:    true,
			},

			"parent_container_id": schema.StringAttribute{
				Description: "the project's parent container ID",
				Computed:    true,
			},

			"name": schema.StringAttribute{
				Description: "the project name",
				Computed:    true,
			},

			"billing_ref": schema.StringAttribute{
				Description: "billing reference for cost transparency",
				Computed:    true,
			},
		},
	}
}
