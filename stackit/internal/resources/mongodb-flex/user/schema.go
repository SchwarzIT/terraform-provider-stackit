package user

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// User is the schema model
type User struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.String `tfsdk:"instance_id"`
	ProjectID  types.String `tfsdk:"project_id"`
	Password   types.String `tfsdk:"password"`
	Username   types.String `tfsdk:"username"`
	Database   types.String `tfsdk:"database"`
	Host       types.String `tfsdk:"host"`
	Port       types.Int64  `tfsdk:"port"`
	URI        types.String `tfsdk:"uri"`
	Roles      types.List   `tfsdk:"roles"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages MongoDB Flex instance users\n%s",
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
				Description: "the mongo db flex instance id.",
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
				Description: "Specifies the user's username",
				Optional:    true,
				Computed:    true,
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
			"database": schema.StringAttribute{
				Description: "Specifies the database the user can access",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(DefaultDatabase),
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
			"roles": schema.ListAttribute{
				Description: "Specifies the role assigned to the user, valid options are: `readWrite` or `read`",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("readWrite", "read"),
					),
					listvalidator.SizeAtMost(1),
				},
				Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(DefaultRole),
				})),
			},
		},
	}
}
