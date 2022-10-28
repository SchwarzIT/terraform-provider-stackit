package bucket

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/buckets"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Bucket is the schema model
type Bucket struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ProjectID    types.String `tfsdk:"project_id"`
	Region       types.String `tfsdk:"region"`
	HostStyleURL types.String `tfsdk:"host_style_url"`
	PathStyleURL types.String `tfsdk:"path_style_url"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages Object Storage buckets",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "the bucket name",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(buckets.ValidateBucketName, "validate bucket name"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"project_id": {
				Description: "project ID the bucket belongs to",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"region": {
				Description: "the region where the bucket was created",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"host_style_url": {
				Description: "url with dedicated host name",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"path_style_url": {
				Description: "url with path to the bucket",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}, nil
}
