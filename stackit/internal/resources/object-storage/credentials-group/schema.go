package credentialsgroup

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CredentialsGroup is the schema model
type CredentialsGroup struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ObjectStorageProjectID types.String `tfsdk:"object_storage_project_id"`
	URN                    types.String `tfsdk:"urn"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Object Storage credential groups",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "the credential group ID",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},

			"object_storage_project_id": schema.StringAttribute{
				Description: "The ID returned from `stackit_object_storage_project`",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"name": schema.StringAttribute{
				Description: "the credential group display name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"urn": schema.StringAttribute{
				Description: "credential group URN",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}
}
