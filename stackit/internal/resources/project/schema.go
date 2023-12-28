package project

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Project is the schema model
type Project struct {
	ID                types.String      `tfsdk:"id"`
	ContainerID       types.String      `tfsdk:"container_id"`
	ParentContainerID types.String      `tfsdk:"parent_container_id"`
	Name              types.String      `tfsdk:"name"`
	BillingRef        types.String      `tfsdk:"billing_ref"`
	OwnerEmail        types.String      `tfsdk:"owner_email"`
	Timeouts          timeouts.Value    `tfsdk:"timeouts"`
	Labels            map[string]string `tfsdk:"labels"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages STACKIT projects\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "the project ID",
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"container_id": schema.StringAttribute{
				Description: "the project container ID",
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"parent_container_id": schema.StringAttribute{
				Description: "the container ID in which the project will be created",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				Description: "the project name",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectName(),
				},
			},

			"billing_ref": schema.StringAttribute{
				Description: "billing reference for cost transparency",
				Required:    true,
				Validators: []validator.String{
					validate.BillingRef(),
				},
			},

			"owner_email": schema.StringAttribute{
				Description: "Email address of owner of the project. This value is only considered during creation. changing it afterwards will have no effect.",
				Required:    true,
			},

			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),

			"labels": schema.MapAttribute{
				Description: "Extend project information with custom label values.",
				Required:    true,
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
				Validators: []validator.Map{
					validate.ReserveProjectLabels(),
				},
			},
		},
	}
}
