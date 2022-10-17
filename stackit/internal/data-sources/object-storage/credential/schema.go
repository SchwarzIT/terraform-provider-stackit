package credential

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Credential is the schema model
type Credential struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Expiry      types.String `tfsdk:"expiry"`
	DisplayName types.String `tfsdk:"display_name"`
}

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Data source for Object Storage credentials",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "the credential ID",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},

			"project_id": {
				Description: "project ID the credential belongs to",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"expiry": {
				Type:     types.StringType,
				Computed: true,
			},

			"display_name": {
				Description: "the credential's display name in the portal",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
		},
	}, nil
}
