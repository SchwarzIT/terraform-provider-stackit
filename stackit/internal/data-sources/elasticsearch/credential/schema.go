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
	URI             types.String `tfsdk:"uri"`
	Host            types.String `tfsdk:"host"`
	Port            types.Int64  `tfsdk:"port"`
	Hosts           types.List   `tfsdk:"hosts"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	CACert          types.String `tfsdk:"password"`
	Schema          types.String `tfsdk:"password"`
	SyslogDrainURL  types.String `tfsdk:"syslog_drain_url"`
	RouteServiceURL types.String `tfsdk:"route_service_url"`
}

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages Elasticsearch credentials",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
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
				Description: "Elasticsearch instance ID the credential belongs to",
				Type:        types.StringType,
				Required:    true,
			},

			"host": {
				Description: "Credential host",
				Type:        types.StringType,
				Computed:    true,
			},
			"port": {
				Description: "Credential port",
				Type:        types.Int64Type,
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
			"cacrt": {
				Description: "Credential CA Certificate",
				Type:        types.StringType,
				Computed:    true,
			},
			"schema": {
				Description: "Credential schema",
				Type:        types.StringType,
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
		},
	}, nil
}
