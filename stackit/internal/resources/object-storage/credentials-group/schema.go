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
	Name      types.String `tfsdk:"name"`
	ProjectID types.String `tfsdk:"project_id"`
	URN       types.String `tfsdk:"urn"`
}

// GetSchema returns the terraform schema structure
func (r Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages Object Storage credential groups",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "the credential group ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
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
				Description: "the credential group display name",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
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
	}, nil
}
