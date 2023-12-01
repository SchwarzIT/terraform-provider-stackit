package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID                 types.String   `tfsdk:"id"`
	Name               types.String   `tfsdk:"name"`
	ProjectID          types.String   `tfsdk:"project_id"`
	Plan               types.String   `tfsdk:"plan"`
	PlanID             types.String   `tfsdk:"plan_id"`
	Version            types.String   `tfsdk:"version"`
	ACL                types.List     `tfsdk:"acl"`
	DashboardURL       types.String   `tfsdk:"dashboard_url"`
	CFGUID             types.String   `tfsdk:"cf_guid"`
	CFSpaceGUID        types.String   `tfsdk:"cf_space_guid"`
	CFOrganizationGUID types.String   `tfsdk:"cf_organization_guid"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages %s instances\n%s",
			r.service.Display(),
			common.EnvironmentInfo(r.urls),
		),
		DeprecationMessage: func() string {
			if r.service.Display() == "ElasticSearch" {
				return "This resource is deprecated and will be removed in a future release. Please use the `stackit_opensearch_instance` resource instead."
			}
			return ""
		}(),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated. Changing this value requires the resource to be recreated.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s Plan. Default is `%s`", r.service.Display(), r.getDefaultPlan()),
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(r.getDefaultPlan()),
			},
			"plan_id": schema.StringAttribute{
				Description: "The selected plan ID",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("%s version. Default is %s", r.service.Display(), r.getDefaultVersion()),
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(r.getDefaultVersion()),
			},
			"acl": schema.ListAttribute{
				Description: "Access Control rules to whitelist IP addresses",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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
			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}
