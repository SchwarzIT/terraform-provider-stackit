package postgresinstance

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gorm.io/gorm/schema"
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

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages Postgres Flex instances
		
~> **Note:** Postgres Flex is in 'alpha' stage in STACKIT
`,
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": {
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": {
				Description: "The project ID the instance runs in. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"machine_type": {
				Description: "The Machine Type. Available options: `c1.2` `m1.2`, `c1.3`, `m1.3`, `c1.4`, `c1.5`, `m1.5`",
				Type:        types.StringType,
				Required:    true,
			},
			"version": {
				Description: "Postgres version. Options: `13`, `14`. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					modifiers.StringDefault(default_version),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replicas": {
				Description: "Number of replicas (Default is `1`). Changing this value requires the resource to be recreated.",
				Type:        types.Int64Type,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					modifiers.Int64Default(default_replicas),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_schedule": {
				Description: "Specifies the backup schedule (cron style)",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					modifiers.StringDefault(default_backup_schedule),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage": {
				Description: "A signle `storage` block as defined below. Changing this value requires the resource to be recreated.",
				Optional:    true,
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"class": {
						Description: "Specifies the storage class. Available option: `premium-perf6-stackit`",
						Type:        types.StringType,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefault(default_storage_class),
						},
					},
					"size": {
						Description: "The storage size in GB",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							modifiers.Int64Default(default_storage_size),
						},
					},
				}),
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": {
				Description: "The databse admin user",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Description: "Specifies the user id",
						Type:        types.StringType,
						Computed:    true,
					},
					"username": {
						Description: "Specifies the user's username",
						Type:        types.StringType,
						Computed:    true,
					},
					"password": {
						Description: "Specifies the user's password",
						Type:        types.StringType,
						Computed:    true,
						Sensitive:   true,
					},
					"database": {
						Description: "Specifies the database the user can access",
						Type:        types.StringType,
						Computed:    true,
					},
					"hostname": {
						Description: "Specifies the allowed user hostname",
						Type:        types.StringType,
						Computed:    true,
					},
					"port": {
						Description: "Specifies the port",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"uri": {
						Description: "Specifies connection URI",
						Type:        types.StringType,
						Computed:    true,
						Sensitive:   true,
					},
					"roles": {
						Description: "Specifies the roles assigned to the user",
						Type:        types.ListType{ElemType: types.StringType},
						Computed:    true,
					},
				}),
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"options": {
				Description: "Specifies postgres instance options",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"labels": {
				Description: "Instance Labels",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"acl": {
				Description: "Access Control rules to whitelist IP addresses",
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
			},
		},
	}
}
