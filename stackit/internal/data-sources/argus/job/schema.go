package job

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/instances"
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
		Description: "Data source for Argus Jobs",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "Specifies the Argus Job ID",
				Type:        types.StringType,
				Computed:    true,
			},

			"name": {
				Description: "Specifies the name of the scraping job",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.StringWith(
						instances.ValidateInstanceName,
						"validate argus instance name",
					),
				},
			},

			"project_id": {
				Description: "Specifies the Project ID the Argus instance belongs to",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			"argus_instance_id": {
				Description: "Specifies the Argus Instance ID the job belongs to",
				Type:        types.StringType,
				Required:    true,
			},

			"metrics_path": {
				Description: "Specifies the job scraping path. Defaults to `/metrics`",
				Type:        types.StringType,
				Computed:    true,
			},

			"scheme": {
				Description: "Specifies the scheme. Default is `https`.",
				Type:        types.StringType,
				Computed:    true,
			},

			"scrape_interval": {
				Description: "Specifies the scrape interval as duration string. Default is `5m`.",
				Type:        types.StringType,
				Computed:    true,
			},

			"scrape_timeout": {
				Description: "Specifies the scrape timeout as duration string. Default is `2m`.",
				Type:        types.StringType,
				Computed:    true,
			},

			"saml2": {
				Description: "A saml2 configuration block",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enable_url_parameters": {
						Description: "Should URL parameters be enabled? Default is `true`",
						Type:        types.BoolType,
						Computed:    true,
					},
				}),
			},

			"basic_auth": {
				Description: "A basic_auth block",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"username": {
						Description: "Specifies basic auth username",
						Type:        types.StringType,
						Computed:    true,
					},
					"password": {
						Description: "Specifies basic auth password",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},

			"targets": {
				Description: "targets list",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"urls": {
						Description: "Specifies basic auth username",
						Type:        types.ListType{ElemType: types.StringType},
						Computed:    true,
					},
					"labels": {
						Description: "Specifies basic auth password",
						Type:        types.MapType{ElemType: types.StringType},
						Computed:    true,
					},
				}),
			},
		},
	}
}
