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
	ProjectID  types.String `tfsdk:"project_id"`
	InstanceID types.String `tfsdk:"instance_id"`

	ID  types.String `tfsdk:"id"`
	URI types.String `tfsdk:"uri"`

	SyslogDrainURL  types.String `tfsdk:"syslog_drain_url"`
	RouteServiceURL types.String `tfsdk:"route_service_url"`
	VolumeMounts    types.List   `tfsdk:"volume_mounts"`

	Host      types.String `tfsdk:"host"`
	Port      types.Int64  `tfsdk:"port"`
	Hosts     types.List   `tfsdk:"hosts"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	Protocols types.Map    `tfsdk:"protocols"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"volume_mounts": {
				Description: "Credential volume_mounts",
				Type:        types.ListType{ElemType: types.MapType{ElemType: types.StringType}},
				Computed:    true,
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
			"protocols": {
				Description: "Credential protocols",
				Type:        types.MapType{ElemType: types.StringType},
				Computed:    true,
			},
		},
	}, nil
}
