package bucket

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/buckets"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gorm.io/gorm/schema"
)

// Bucket is the schema model
type Bucket struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ObjectStorageProjectID types.String `tfsdk:"object_storage_project_id"`
	Region                 types.String `tfsdk:"region"`
	HostStyleURL           types.String `tfsdk:"host_style_url"`
	PathStyleURL           types.String `tfsdk:"path_style_url"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Object Storage buckets",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "the bucket name",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.StringWith(buckets.ValidateBucketName, "validate bucket name"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"object_storage_project_id": {
				Description: "The ID returned from `stackit_object_storage_project`",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
	}
}
