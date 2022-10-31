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
}

// Storage represent instance storage
type Storage struct {
	Class types.String `tfsdk:"class"`
	Size  types.Int64  `tfsdk:"size"`
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
				Description: "Specifies the cluster name (lower case, alphanumeric, hypens allowed, up to 11 chars)",
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
				Description: "The project ID the cluster runs in",
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
				Description: "Postgres version",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(clientValidate.SemVer, "validate postgres version"),
				},
			},
			"replicas": {
				Description: "How many replicas of the database should exist",
				Type:        types.Int64Type,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.Int64Default(1),
				},
			},
			"backup_schedule": {
				Description: "Specifies the backup schedule (cron style)",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"storage": {
				Description: "One or more `node_pool` block as defined below",
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
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
