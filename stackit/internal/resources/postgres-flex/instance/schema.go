package postgresinstance

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID             types.String      `tfsdk:"id"`
	Name           types.String      `tfsdk:"name"`
	ProjectID      types.String      `tfsdk:"project_id"`
	MachineType    types.String      `tfsdk:"machine_type"`
	Version        types.String      `tfsdk:"version"`
	Replicas       types.Int64       `tfsdk:"replicas"`
	BackupSchedule types.String      `tfsdk:"backup_schedule"`
	Options        map[string]string `tfsdk:"options"`
	Labels         map[string]string `tfsdk:"labels"`
	ACL            types.List        `tfsdk:"acl"`
	Storage        types.Object      `tfsdk:"storage"`
	User           types.Object      `tfsdk:"user"`
}

// Storage represent instance storage
type Storage struct {
	Class types.String `tfsdk:"class"`
	Size  types.Int64  `tfsdk:"size"`
}

// User represent database user
type User struct {
	ID       types.String `tfsdk:"id"`
	Password types.String `tfsdk:"password"`
	Username types.String `tfsdk:"username"`
	Database types.String `tfsdk:"database"`
	Hostname types.String `tfsdk:"hostname"`
	Port     types.Int64  `tfsdk:"port"`
	URI      types.String `tfsdk:"uri"`
	Roles    types.List   `tfsdk:"roles"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages Postgres Flex instances
		
~> **Note:** Postgres Flex is in 'alpha' stage in STACKIT
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID the instance runs in. Changing this value requires the resource to be recreated.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"machine_type": schema.StringAttribute{
				Description: "The Machine Type. Available options: `c1.2` `m1.2`, `c1.3`, `m1.3`, `c1.4`, `c1.5`, `m1.5`",
				Required:    true,
			},
			"version": schema.StringAttribute{
				Description: "Postgres version. Options: `13`, `14`. Changing this value requires the resource to be recreated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					modifiers.StringDefault(default_version),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replicas": schema.Int64Attribute{
				Description: "Number of replicas (Default is `1`). Changing this value requires the resource to be recreated.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					modifiers.Int64Default(default_replicas),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"backup_schedule": schema.StringAttribute{
				Description: "Specifies the backup schedule (cron style)",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					modifiers.StringDefault(default_backup_schedule),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage": schema.SingleNestedAttribute{
				Description: "A signle `storage` block as defined below.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"class": schema.StringAttribute{
						Description: "Specifies the storage class. Available option: `premium-perf6-stackit`",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefault(default_storage_class),
						},
					},
					"size": schema.Int64Attribute{
						Description: "The storage size in GB",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							modifiers.Int64Default(default_storage_size),
						},
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
					"hostname": schema.StringAttribute{
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
				Description: "Specifies postgres instance options",
				ElementType: types.StringType,
				Optional:    true,
			},
			"labels": schema.MapAttribute{
				Description: "Instance Labels",
				ElementType: types.StringType,
				Optional:    true,
			},
			"acl": schema.ListAttribute{
				Description: "Access Control rules to whitelist IP addresses",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}
