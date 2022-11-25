package instance

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `Data source for MariaDB instances

~> **Note:** MariaDB API (Part of DSA APIs) currently has issues reflecting updates & configuration correctly. Therefore, this data source is not ready for production usage.		
		`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "The instance ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				Description: "Specifies the instance name.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"project_id": {
				Description: "The project ID.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
			},
			"plan": {
				Description: "The MariaDB plan name",
				Type:        types.StringType,
				Computed:    true,
			},
			"plan_id": {
				Description: "MariaDB plan ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"version": {
				Description: "MariaDB version",
				Type:        types.StringType,
				Computed:    true,
			},
			"acl": {
				Description: "Access control rules",
				Type:        types.ListType{ElemType: types.StringType},
				Computed:    true,
			},
			"dashboard_url": {
				Description: "Dashboard URL",
				Type:        types.StringType,
				Computed:    true,
			},
			"cf_guid": {
				Description: "Cloud Foundry GUID",
				Type:        types.StringType,
				Computed:    true,
			},
			"cf_space_guid": {
				Description: "Cloud Foundry Space GUID",
				Type:        types.StringType,
				Computed:    true,
			},
			"cf_organization_guid": {
				Description: "Cloud Foundry Organization GUID",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}
