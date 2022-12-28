package instance

import (
	"context"

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
		MarkdownDescription: `Data source for MongoDB Flex instance
		
~> **Note:** MongoDB Flex is in 'beta' stage in STACKIT
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID the instance runs in.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"machine_type": schema.StringAttribute{
				Description: "The Machine Type.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "MongoDB version.",
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
				Description: "A signle `storage` block as defined below.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"class": schema.StringAttribute{
						Description: "Specifies the storage class. Available option: `premium-perf2-mongodb`",
						Computed:    true,
					},
					"size": schema.Int64Attribute{
						Description: "The storage size in GB.",
						Computed:    true,
					},
				},
			},
			"user": schema.SingleNestedAttribute{
				Description: "The databse admin user",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "Specifies the user id",
						Computed:    true,
					},
					"username": schema.StringAttribute{
						Description: "Specifies the user's username",
						Computed:    true,
					},
					"password": schema.StringAttribute{
						Description: "Specifies the user's password",
						Computed:    true,
						Sensitive:   true,
					},
					"database": schema.StringAttribute{
						Description: "Specifies the database the user can access",
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "Specifies the allowed user hostname",
						Computed:    true,
					},
					"port": schema.Int64Attribute{
						Description: "Specifies the port",
						Computed:    true,
					},
					"uri": schema.StringAttribute{
						Description: "Specifies connection URI",
						Computed:    true,
						Sensitive:   true,
					},
					"roles": schema.ListAttribute{
						Description: "Specifies the roles assigned to the user",
						ElementType: types.StringType,
						Computed:    true,
					},
				},
			},
			"options": schema.MapAttribute{
				Description: "Specifies mongodb instance options",
				ElementType: types.StringType,
				Computed:    true,
			},
			"labels": schema.MapAttribute{
				Description: "Instance Labels",
				ElementType: types.StringType,
				Computed:    true,
			},
			"acl": schema.ListAttribute{
				Description: "Access Control rules to whitelist IP addresses",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}
