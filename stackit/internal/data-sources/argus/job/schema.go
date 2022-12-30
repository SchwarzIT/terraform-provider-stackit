package job

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for Argus Instance Jobs",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the Argus Job ID",
				Computed:    true,
			},

			"name": schema.StringAttribute{
				Description: "Specifies the name of the scraping job",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 200),
				},
			},

			"project_id": schema.StringAttribute{
				Description: "Specifies the Project ID the Argus instance belongs to",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			"argus_instance_id": schema.StringAttribute{
				Description: "Specifies the Argus Instance ID the job belongs to",
				Required:    true,
			},

			"metrics_path": schema.StringAttribute{
				Description: "Specifies the job scraping path.",
				Computed:    true,
			},

			"scheme": schema.StringAttribute{
				Description: "Specifies the scheme.",
				Computed:    true,
			},

			"scrape_interval": schema.StringAttribute{
				Description: "Specifies the scrape interval as duration string.",
				Computed:    true,
			},

			"scrape_timeout": schema.StringAttribute{
				Description: "Specifies the scrape timeout as duration string.",
				Computed:    true,
			},

			"saml2": schema.SingleNestedAttribute{
				Description: "A saml2 configuration block",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enable_url_parameters": schema.BoolAttribute{
						Description: "Should URL parameters be enabled?",
						Computed:    true,
					},
				},
			},

			"basic_auth": schema.SingleNestedAttribute{
				Description: "A basic_auth block",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Description: "Specifies basic auth username",
						Computed:    true,
					},
					"password": schema.StringAttribute{
						Description: "Specifies basic auth password",
						Computed:    true,
					},
				},
			},

			"targets": schema.ListNestedAttribute{
				Description: "targets list",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"urls": schema.ListAttribute{
							Description: "Specifies target URLs",
							ElementType: types.StringType,
							Computed:    true,
						},
						"labels": schema.MapAttribute{
							Description: "Specifies labels",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
