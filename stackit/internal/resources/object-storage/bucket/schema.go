package bucket

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

// Bucket is the schema model
type Bucket struct {
	ID                     types.String   `tfsdk:"id"`
	Name                   types.String   `tfsdk:"name"`
	ProjectID              types.String   `tfsdk:"project_id"`
	ObjectStorageProjectID types.String   `tfsdk:"object_storage_project_id"`
	Region                 types.String   `tfsdk:"region"`
	HostStyleURL           types.String   `tfsdk:"host_style_url"`
	PathStyleURL           types.String   `tfsdk:"path_style_url"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages Object Storage buckets\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},

			"name": schema.StringAttribute{
				Description: "Bucket name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"object_storage_project_id": schema.StringAttribute{
				DeprecationMessage: "This attribute is deprecated and will be removed in a future version of the provider. Please use `project_id` instead.",
				Description:        "The ID returned from `stackit_object_storage_project`",
				Optional:           true,
				Computed:           true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"project_id": schema.StringAttribute{
				Description: "The project UUID.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"region": schema.StringAttribute{
				Description: "the region where the bucket was created",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"host_style_url": schema.StringAttribute{
				Description: "url with dedicated host name",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"path_style_url": schema.StringAttribute{
				Description: "url with path to the bucket",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}
