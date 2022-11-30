package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
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
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `Manages PostgresDB instances

~> **Note:** PostgresDB API (Part of DSA APIs) currently has issues reflecting updates & configuration correctly. Therefore, this resource is not ready for production usage.
		`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated. Changing this value requires the resource to be recreated.",
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
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"plan": {
				Description: fmt.Sprintf("The PostgresDB Plan. Default is `%s`.\nOptions are: `stackit-postgresql-cluster-big`, `stackit-postgresql-cluster-extra-large`, `stackit-postgresql-cluster-medium`, `stackit-postgresql-cluster-small`, `stackit-postgresql-single-medium`, `stackit-postgresql-single-small`, `stackit-postgresql-cluster-extra-large-non-ssl`, `stackit-postgresql-cluster-medium-non-ssl`, `stackit-postgresql-single-small-non-ss`, `stackit-postgresql-cluster-big-non-ssl`, `stackit-postgresql-cluster-small-non-ssl`, `stackit-postgresql-single-medium-non-ssl`", default_plan),
				Type:        types.StringType,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.StringDefault(default_plan),
				},
			},
			"plan_id": {
				Description: "The selected plan ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"version": {
				Description: "PostgresDB version. Options: `9.4`, `10`, `11`, `13`. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
					modifiers.StringDefault(default_version),
				},
			},
			"acl": {
				Description: "Access Control rules to whitelist IP addresses",
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
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
