package credentialsgroup

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gorm.io/gorm/schema"
)

// CredentialsGroup is the schema model
type CredentialsGroup struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ObjectStorageProjectID types.String `tfsdk:"object_storage_project_id"`
	URN                    types.String `tfsdk:"urn"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Object Storage credential groups",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "the credential group ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
			},

			"object_storage_project_id": {
				Description: "The ID returned from `stackit_object_storage_project`",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"name": {
				Description: "the credential group display name",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"urn": {
				Description: "credential group URN",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}
}
