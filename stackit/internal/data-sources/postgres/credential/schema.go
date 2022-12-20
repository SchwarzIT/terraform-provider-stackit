package credential

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Credential is the schema model
type Credential struct {
	ID              types.String `tfsdk:"id"`
	ProjectID       types.String `tfsdk:"project_id"`
	InstanceID      types.String `tfsdk:"instance_id"`
	Host            types.String `tfsdk:"host"`
	Hosts           types.List   `tfsdk:"hosts"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	Port            types.Int64  `tfsdk:"port"`
	SyslogDrainURL  types.String `tfsdk:"syslog_drain_url"`
	RouteServiceURL types.String `tfsdk:"route_service_url"`
	URI             types.String `tfsdk:"uri"`
}

// GetSchema returns the terraform schema structure
func (d *DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages Postgres credentials",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"project_id": {
				Description: "Project ID the credential belongs to",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"instance_id": {
				Description: "Postgres instance ID the credential belongs to",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			// Computed:
			"host": {
				Description: "Credential host",
				Type:        types.StringType,
				Computed:    true,
			},
			"hosts": {
				Description: "Credential hosts",
				Type:        types.ListType{ElemType: types.StringType},
				Computed:    true,
			},
			"username": {
				Description: "Credential username",
				Type:        types.StringType,
				Computed:    true,
			},
			"password": {
				Description: "Credential password",
				Type:        types.StringType,
				Computed:    true,
			},
			"port": {
				Description: "Credential port",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"syslog_drain_url": {
				Description: "Credential syslog_drain_url",
				Type:        types.StringType,
				Computed:    true,
			},
			"route_service_url": {
				Description: "Credential route_service_url",
				Type:        types.StringType,
				Computed:    true,
			},
			"uri": {
				Description: "The instance URI",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}
