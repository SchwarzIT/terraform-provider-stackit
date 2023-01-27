package instance

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

// Instance is the schema model
type Instance struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ProjectID      types.String `tfsdk:"project_id"`
	MachineType    types.String `tfsdk:"machine_type"`
	Version        types.String `tfsdk:"version"`
	Replicas       types.Int64  `tfsdk:"replicas"`
	BackupSchedule types.String `tfsdk:"backup_schedule"`
	ACL            types.List   `tfsdk:"acl"`
	Storage        types.Object `tfsdk:"storage"`
}

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for Postgres Flex instance\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"machine_type": schema.StringAttribute{
				Description: "The Machine Type",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "Postgres version",
				Computed:    true,
			},
			"replicas": schema.Int64Attribute{
				Description: "Number of replicas",
				Computed:    true,
			},
			"backup_schedule": schema.StringAttribute{
				Description: "Specifies the backup schedule (cron style)",
				Computed:    true,
			},
			"storage": schema.SingleNestedAttribute{
				Description: "A signle `storage` block as defined below",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"class": schema.StringAttribute{
						Description: "Specifies the storage class. Available option: `premium-perf6-stackit`",
						Computed:    true,
					},
					"size": schema.Int64Attribute{
						Description: "The storage size in GB",
						Computed:    true,
					},
				},
			},
			"acl": schema.ListAttribute{
				Description: "Access Control rules to whitelist IP addresses",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}
