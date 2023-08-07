package credential

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Credential is the schema model
type Credential struct {
	ID                     types.String `tfsdk:"id"`
	ObjectStorageProjectID types.String `tfsdk:"object_storage_project_id"`
	CredentialsGroupID     types.String `tfsdk:"credentials_group_id"`
	Expiry                 types.String `tfsdk:"expiry"`
	DisplayName            types.String `tfsdk:"display_name"`
}

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for Object Storage credentials\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "the credential ID",
				Optional:    true,
				Computed:    true,
			},

			"object_storage_project_id": schema.StringAttribute{
				Description: "The ID returned from `stackit_object_storage_project`",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			"credentials_group_id": schema.StringAttribute{
				Description: "the credentials group ID",
				Required:    true,
			},

			"expiry": schema.StringAttribute{
				Computed: true,
			},

			"display_name": schema.StringAttribute{
				Description: "the credential's display name in the portal",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}
