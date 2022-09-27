package credentialsgroup

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CredentialsGroup is the schema model
type CredentialsGroup struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	URN       types.String `tfsdk:"urn"`
}

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Data source for Object Storage buckets",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "the credential group ID",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"project_id": {
				Description: "project ID the credential group belongs to",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"name": {
				Description: "the credential group's display name",
				Type:        types.StringType,
				Computed:    true,
				Optional:    false,
				Required:    false,
			},

			"urn": {
				Description: "credential group URN",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}, nil
}
