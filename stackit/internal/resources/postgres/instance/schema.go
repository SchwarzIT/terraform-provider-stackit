package postgresinstance

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PostgresInstance is the schema model
type PostgresInstance struct {
	ID             types.String      `tfsdk:"id"`
	Name           types.String      `tfsdk:"name"`
	ProjectID      types.String      `tfsdk:"project_id"`
	FlavorID       types.String      `tfsdk:"flavor_id"`
	Version        types.String      `tfsdk:"version"`
	Replicas       types.Int64       `tfsdk:"version"`
	BackupSchedule types.String      `tfsdk:"backup_schedule"`
	Options        map[string]string `tfsdk:"options"`
	Labels         map[string]string `tfsdk:"labels"`
	ACL            []string          `tfsdk:"labels"`
	Storage        Storage           `tfsdk:"storage"`
	Users          []User            `tfsdk:"users"`
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
	Roles    []string     `tfsdk:"roles"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages kubernetes clusters",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(clusters.ValidateClusterName, "validate cluster name"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"project_id": {
				Description: "The project ID the cluster runs in. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"flavor_id": {
				Description: "The Flavor ID",
				Type:        types.StringType,
				Required:    true,
			},
			"version": {
				Description: "Postgres version. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(clientValidate.SemVer, "validate postgres version"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"replicas": {
				Description: "How many replicas of the database should exist. Changing this value requires the resource to be recreated.",
				Type:        types.Int64Type,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.Int64Default(1),
					resource.RequiresReplace(),
				},
			},
			"backup_schedule": {
				Description: "Specifies the backup schedule (cron style)",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"storage": {
				Description: "A signle `storage` block as defined below. Changing this value requires the resource to be recreated.",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"class": {
						Description: "Specifies the storage class",
						Type:        types.StringType,
						Required:    true,
					},
					"size": {
						Description: "The storage size in GB",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(20),
						},
					},
				}),
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"users": {
				Description: "One or more databse users",
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
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
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
				}),
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
				Description: "Instance Labels",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
		},
	}, nil
}
