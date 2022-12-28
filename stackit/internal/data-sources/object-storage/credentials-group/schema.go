package credentialsgroup

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Data source for Object Storage credential groups",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "the credential group ID",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Description: "the credential group's display name",
				Type:        types.StringType,
				Computed:    true,
			},

			"urn": {
				Description: "credential group URN",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}
}
