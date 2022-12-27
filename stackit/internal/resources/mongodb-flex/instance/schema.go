package instance

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
	MachineType    types.String      `tfsdk:"machine_type"` // aka FlavorID
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
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
	URI      types.String `tfsdk:"uri"`
	Roles    types.List   `tfsdk:"roles"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages MongoDB Flex instances
		
~> **Note:** MongoDB Flex is in 'beta' stage in STACKIT
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
				Description: "The Machine Type. Available options: `T1.2`, `C1.1`, `G1.1`, `M1.1`, `C1.2`, `G1.2`, `M1.2`, `C1.3`, `G1.3`, `M1.3`, `C1.4`, `G1.4`, `M1.4`, `C1.5`, `G1.5`",
				Required:    true,
			},
			"version": schema.StringAttribute{
				Description: "MongoDB version. Version `5.0` and `6.0` are supported. Changing this value requires the resource to be recreated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					modifiers.StringDefault(default_version),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replicas": schema.Int64Attribute{
				Description: "Number of replicas (Default is `1`)",
				Optional:    true,
				Computed:    true,
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
						Description: "Specifies the storage class. Available option: `premium-perf2-mongodb`",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefault(default_storage_class),
						},
					},
					"size": schema.Int64Attribute{
						Description: "The storage size in GB. Default is `10`.",
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
				Optional:    true,
			},
			"labels": schema.MapAttribute{
				Description: "Instance Labels",
				ElementType: types.StringType,

				Optional: true,
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
