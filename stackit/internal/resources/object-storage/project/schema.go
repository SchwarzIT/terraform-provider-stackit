package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gorm.io/gorm/schema"
)

// ObjectStorageProject is the schema model
type ObjectStorageProject struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource enables STACKIT Object Storage in a project",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "object storage project ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"project_id": {
				Description: "the project ID that Object Storage will be enabled in",
				Type:        types.StringType,
				Required:    true,
			},
		},
	}
}
