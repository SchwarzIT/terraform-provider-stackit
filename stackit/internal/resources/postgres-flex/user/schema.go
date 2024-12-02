package user

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// User is the schema model
type User struct {
	ID         types.String   `tfsdk:"id"`
	InstanceID types.String   `tfsdk:"instance_id"`
	ProjectID  types.String   `tfsdk:"project_id"`
	Password   types.String   `tfsdk:"password"`
	Username   types.String   `tfsdk:"username"`
	Host       types.String   `tfsdk:"host"`
	Port       types.Int64    `tfsdk:"port"`
	URI        types.String   `tfsdk:"uri"`
	RoleSet    types.Set      `tfsdk:"role_set"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages Postgres Flex instance users\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Description: "the postgres db flex instance id.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID the instance runs in. Changing this value requires the resource to be recreated.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Specifies the username. Defaults to `psqluser`",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.NoneOfCaseInsensitive("stackit"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(DefaultUsername),
			},
			"password": schema.StringAttribute{
				Description: "Specifies the user's password",
				Computed:    true,
				Sensitive:   true,
			},
			"host": schema.StringAttribute{
				Description: "Specifies the allowed user hostname",
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Specifies the port",
				Computed:    true,
			},
			"uri": schema.StringAttribute{
				Description: "Specifies connection URI",
				Computed:    true,
				Sensitive:   true,
			},
			"role_set": schema.SetAttribute{
				Description: "Specifies the roles assigned to the user, valid options are: `login`, `createdb`",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("login", "createdb"),
					),
				},
				Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue(DefaultRole),
				})),
			},
			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Delete: true,
			}),
		},
	}
}
