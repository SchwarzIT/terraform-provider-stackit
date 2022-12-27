package credential

import (
	"context"

	credentialsgroup "github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/credentials-group"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gorm.io/gorm/schema"
)

// Credential is the schema model
type Credential struct {
	ID                     types.String `tfsdk:"id"`
	ObjectStorageProjectID types.String `tfsdk:"object_storage_project_id"`
	CredentialsGroupID     types.String `tfsdk:"credentials_group_id"`
	Expiry                 types.String `tfsdk:"expiry"`
	DisplayName            types.String `tfsdk:"display_name"`
	AccessKey              types.String `tfsdk:"access_key"`
	SecretAccessKey        types.String `tfsdk:"secret_access_key"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Object Storage credentials",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "the credential ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
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

			"credentials_group_id": {
				Description: "credential group ID. changing this field will recreate the credential.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validate.StringWith(credentialsgroup.ValidateCredentialsGroupID, "credentials group ID"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"expiry": {
				Description: "specifies if the credential should expire. changing this field will recreate the credential.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validate.StringWith(clientValidate.ISO8601, "validate expiry is ISO-8601 compatible"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
	}
}
