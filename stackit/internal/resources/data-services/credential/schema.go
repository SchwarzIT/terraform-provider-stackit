package credential

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Credential is the schema model
type Credential struct {
	ID              types.String `tfsdk:"id"`
	ProjectID       types.String `tfsdk:"project_id"`
	InstanceID      types.String `tfsdk:"instance_id"`
	Host            types.String `tfsdk:"host"`
	Hosts           types.List   `tfsdk:"hosts"`
	DatabaseName    types.String `tfsdk:"database_name"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	Port            types.Int64  `tfsdk:"port"`
	SyslogDrainURL  types.String `tfsdk:"syslog_drain_url"`
	RouteServiceURL types.String `tfsdk:"route_service_url"`
	URI             types.String `tfsdk:"uri"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: fmt.Sprintf("Manages %s credentials\n\n"+
			"~> **Note:** The following URLs are set for each environment:<br />"+
			"	Prod: 	`%s`<br />"+
			"	QA:		`%s`<br />"+
			"	Dev:	`%s`<br />"+
			"	Override the URLs by setting `%s` environment variable",
			r.service.Display(),
			r.urls.Prod,
			r.urls.QA,
			r.urls.Dev,
			r.urls.OverrideWith,
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"project_id": schema.StringAttribute{
				Description: "Project ID the credential belongs to",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Description: "Instance ID the credential belongs to",
				Required:    true,
			},

			"host": schema.StringAttribute{
				Description: "Credential host",
				Computed:    true,
			},

			"hosts": schema.ListAttribute{
				Description: "Credential hosts",
				ElementType: types.StringType,
				Computed:    true,
			},

			"username": schema.StringAttribute{
				Description: "Credential username",
				Computed:    true,
			},

			"database_name": schema.StringAttribute{
				Description: "Database name",
				Computed:    true,
			},

			"password": schema.StringAttribute{
				Description: "Credential password",
				Computed:    true,
			},

			"port": schema.Int64Attribute{
				Description: "Credential port",
				Computed:    true,
			},

			"syslog_drain_url": schema.StringAttribute{
				Description: "Credential syslog_drain_url",
				Computed:    true,
			},

			"route_service_url": schema.StringAttribute{
				Description: "Credential route_service_url",
				Computed:    true,
			},

			"uri": schema.StringAttribute{
				Description: "The instance URI",
				Computed:    true,
			},
		},
	}
}
