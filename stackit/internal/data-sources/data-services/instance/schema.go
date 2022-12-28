package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProjectID          types.String `tfsdk:"project_id"`
	Plan               types.String `tfsdk:"plan"`
	PlanID             types.String `tfsdk:"plan_id"`
	Version            types.String `tfsdk:"version"`
	ACL                types.List   `tfsdk:"acl"`
	DashboardURL       types.String `tfsdk:"dashboard_url"`
	CFGUID             types.String `tfsdk:"cf_guid"`
	CFSpaceGUID        types.String `tfsdk:"cf_space_guid"`
	CFOrganizationGUID types.String `tfsdk:"cf_organization_guid"`
}

// Schema returns the terraform schema structure

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: fmt.Sprintf(`Data source for %s instances

~> **Note:** %s API (Part of DSA APIs) currently has issues reflecting updates & configuration correctly. Therefore, this data source is not ready for production usage.		
		`, d.service.Display(), d.service.Display()),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"plan": schema.StringAttribute{
				Description: "The RabbitMQ Plan",
				Required:    true,
			},
			"plan_id": schema.StringAttribute{
				Description: "The selected plan ID",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "RabbitMQ version",
				Computed:    true,
			},
			"acl": schema.ListAttribute{
				Description: "Access Control rules to whitelist IP addresses",
				ElementType: types.StringType,
				Computed:    true,
			},
			"dashboard_url": schema.StringAttribute{
				Description: "Dashboard URL",
				Computed:    true,
			},
			"cf_guid": schema.StringAttribute{
				Description: "Cloud Foundry GUID",
				Computed:    true,
			},
			"cf_space_guid": schema.StringAttribute{
				Description: "Cloud Foundry Space GUID",
				Computed:    true,
			},
			"cf_organization_guid": schema.StringAttribute{
				Description: "Cloud Foundry Organization GUID",
				Computed:    true,
			},
		},
	}
}
