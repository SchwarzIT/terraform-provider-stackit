package project

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// KubernetesProject is the schema model
type KubernetesProject struct {
	ID        types.String   `tfsdk:"id"`
	ProjectID types.String   `tfsdk:"project_id"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("This resource enables STACKIT Kubernetes Engine (SKE) in a project\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "kubernetes project ID",
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"project_id": schema.StringAttribute{
				Description: "the project ID that SKE will be enabled in",
				Required:    true,
			},

			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}
