package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID             types.String      `tfsdk:"id"`
	Name           types.String      `tfsdk:"name"`
	ProjectID      types.String      `tfsdk:"project_id"`
	Type           types.String      `tfsdk:"type"`
	MachineType    types.String      `tfsdk:"machine_type"` // aka FlavorID
	Version        types.String      `tfsdk:"version"`
	Replicas       types.Int64       `tfsdk:"replicas"`
	BackupSchedule types.String      `tfsdk:"backup_schedule"`
	Labels         map[string]string `tfsdk:"labels"`
	ACL            types.List        `tfsdk:"acl"`
	Storage        types.Object      `tfsdk:"storage"`
	Timeouts       timeouts.Value    `tfsdk:"timeouts"`
}

// Storage represent instance storage
type Storage struct {
	Class types.String `tfsdk:"class"`
	Size  types.Int64  `tfsdk:"size"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages MongoDB Flex instances\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name.",
				Required:    true,
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
				Description: "The Machine Type. Available options: `1.1` (1 CPU, 1 Memory), `1.2` (1 CPU, 2 Memory), `1.4` (1 CPU, 4 Memory), `1.8` (1 CPU, 8 Memory), `2.4` (2 CPU, 4 Memory), `2.8` (2 CPU, 8 Memory), `2.16` (2 CPU, 16 Memory), `4.8` (4 CPU, 8 Memory), `4.16` (4 CPU, 16 Memory), `4.32` (4 CPU, 32 Memory), `8.16` (8 CPU, 16 Memory), `8.32` (8 CPU, 32 Memory), `8.64` (8 CPU, 64 Memory), `16.32` (16 CPU, 32 Memory), `16.64` (16 CPU, 64 Memory)",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The service type. Available options: `Single`, `Replica`, `Sharded`. Changing this value requires the resource to be recreated.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Single", "Replica", "Sharded"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(DefaultType),
			},
			"version": schema.StringAttribute{
				Description: "MongoDB version. Version `5.0` and `6.0` are supported. ",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString(DefaultVersion),
			},
			"replicas": schema.Int64Attribute{
				Description: fmt.Sprintf("Number of replicas (Default is `%d`).", DefaultReplicas),
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Default: int64default.StaticInt64(DefaultReplicas),
			},
			"backup_schedule": schema.StringAttribute{
				Description: "Specifies the backup schedule (cron style).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(DefaultBackupSchedule),
			},
			"storage": schema.SingleNestedAttribute{
				Description: "A single `storage` block as defined below.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"class": schema.StringAttribute{
						Description: "Specifies the storage class. Available option: `premium-perf2-mongodb`",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(DefaultStorageClass),
					},
					"size": schema.Int64Attribute{
						Description: fmt.Sprintf("The storage size in GB (Default is `%d`).", DefaultStorageSize),
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(DefaultStorageSize),
					},
				},
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
			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}
