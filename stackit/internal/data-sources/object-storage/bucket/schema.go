package bucket

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for Object Storage buckets\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},

			"name": schema.StringAttribute{
				Description: "the bucket name",
				Required:    true,
			},

			"object_storage_project_id": schema.StringAttribute{
				Description: "The ID returned from `stackit_object_storage_project`",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			"region": schema.StringAttribute{
				Description: "the region where the bucket was created",
				Computed:    true,
			},

			"host_style_url": schema.StringAttribute{
				Description: "url with dedicated host name",
				Computed:    true,
			},

			"path_style_url": schema.StringAttribute{
				Description: "url with path to the bucket",
				Computed:    true,
			},
		},
	}
}
