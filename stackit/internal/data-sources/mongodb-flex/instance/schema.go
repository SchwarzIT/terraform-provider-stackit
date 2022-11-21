package instance

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `Data source for MongoDB Flex instance
		
~> **Note:** MongoDB Flex is in 'beta' stage in STACKIT
`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				Description: "Specifies the instance name",
				Type:        types.StringType,
				Required:    true,
			},
			"project_id": {
				Description: "The project ID",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
			},
			"machine_type": {
				Description: "The Machine Type",
				Type:        types.StringType,
				Computed:    true,
			},
			"version": {
				Description: "MongoDB version",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"replicas": {
				Description: "Number of replicas",
				Type:        types.Int64Type,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.Int64Default(1),
					resource.RequiresReplace(),
				},
			},
			"backup_schedule": {
				Description: "Specifies the backup schedule (cron style)",
				Type:        types.StringType,
				Computed:    true,
			},
			"storage": {
				Description: "A signle `storage` block as defined below",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"class": {
						Description: "Specifies the storage class. Available option: `premium-perf6-stackit`",
						Type:        types.StringType,
						Computed:    true,
					},
					"size": {
						Description: "The storage size in GB",
						Type:        types.Int64Type,
						Computed:    true,
					},
				}),
			},
			"acl": {
				Description: "Access Control rules to whitelist IP addresses",
				Type:        types.ListType{ElemType: types.StringType},
				Computed:    true,
			},
		},
	}, nil
}
