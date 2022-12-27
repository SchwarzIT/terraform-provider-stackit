package bucket

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/buckets"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Data source for Object Storage buckets",
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
			},

			"object_storage_project_id": {
				Description: "The ID returned from `stackit_object_storage_project`",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			"region": {
				Description: "the region where the bucket was created",
				Type:        types.StringType,
				Computed:    true,
			},

			"host_style_url": {
				Description: "url with dedicated host name",
				Type:        types.StringType,
				Computed:    true,
			},

			"path_style_url": {
				Description: "url with path to the bucket",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}
}
