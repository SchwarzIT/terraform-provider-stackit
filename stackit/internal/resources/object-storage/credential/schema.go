package credential

import (
	"context"

	credentialsgroup "github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/credentials-group"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Credential is the schema model
type Credential struct {
	ID                 types.String `tfsdk:"id"`
	ProjectID          types.String `tfsdk:"project_id"`
	CredentialsGroupID types.String `tfsdk:"credentials_group_id"`
	Expiry             types.String `tfsdk:"expiry"`
	DisplayName        types.String `tfsdk:"display_name"`
	AccessKey          types.String `tfsdk:"access_key"`
	SecretAccessKey    types.String `tfsdk:"secret_access_key"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages Object Storage credentials",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "the credential ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},

			"project_id": {
				Description: "project ID the credential belongs to. changing this field will recreate the credential.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"credentials_group_id": {
				Description: "credential group ID. changing this field will recreate the credential.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(credentialsgroup.ValidateCredentialsGroupID, "credentials group ID"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"expiry": {
				Description: "specifies if the credential should expire. changing this field will recreate the credential.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(clientValidate.ISO8601, "validate expiry is ISO-8601 compatible"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"display_name": {
				Description: "the credential's display name in the portal",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"access_key": {
				Description: "access key (sensitive)",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Sensitive:   true,
			},

			"secret_access_key": {
				Description: "secret access key (sensitive)",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Sensitive:   true,
			},
		},
	}, nil
}
