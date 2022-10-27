package job

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/instances"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Job is the schema model
type Job struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	ProjectID       types.String `tfsdk:"project_id"`
	ArgusInstanceID types.String `tfsdk:"argus_instance_id"`
	MetricsPath     types.String `tfsdk:"metrics_path"`
	Scheme          types.String `tfsdk:"scheme"`
	ScrapeInterval  types.String `tfsdk:"scrape_interval"`
	ScrapeTimeout   types.String `tfsdk:"scrape_timeout"`
	SAML2           *SAML2       `tfsdk:"saml2"`
	BasicAuth       *BasicAuth   `tfsdk:"basic_auth"`
	Targets         []Target     `tfsdk:"targets"`
}

// SAML2 holds saml configuration
type SAML2 struct {
	EnableURLParameters types.Bool `tfsdk:"enable_url_parameters"`
}

// Target holds targets for scraping
type Target struct {
	URLs   []types.String `tfsdk:"urls"`
	Labels types.Map      `tfsdk:"labels"`
}

// BasicAuth holds basic auth data
type BasicAuth struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages Argus Instance Jobs",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the Argus Job ID",
				Type:        types.StringType,
				Computed:    true,
			},

			"name": {
				Description: "Specifies the name of the scraping job",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(
						instances.ValidateInstanceName,
						"validate argus instance name",
					),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"project_id": {
				Description: "Specifies the Project ID the Argus instance belongs to",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"argus_instance_id": {
				Description: "Specifies the Argus Instance ID the job belongs to",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"metrics_path": {
				Description: "Specifies the job scraping path. Defaults to `/metrics`",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.StringDefault(default_metrics_path),
				},
			},

			"scheme": {
				Description: "Specifies the scheme. Default is `https`.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.StringDefault(default_scheme),
				},
			},

			"scrape_interval": {
				Description: "Specifies the scrape interval as duration string. Default is `5m`.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.StringDefault(default_scrape_interval),
				},
			},

			"scrape_timeout": {
				Description: "Specifies the scrape timeout as duration string. Default is `2m`.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.StringDefault(default_scrape_timeout),
				},
			},

			"saml2": {
				Description: "A saml2 configuration block",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enable_url_parameters": {
						Description: "Should URL parameters be enabled? Default is `true`",
						Type:        types.BoolType,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.BoolDefault(default_saml2_enable_url_parameters),
						},
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
						Required:    true,
					},
					"password": {
						Description: "Specifies basic auth password",
						Type:        types.StringType,
						Required:    true,
					},
				}),
			},

			"targets": {
				Description: "targets list",
				Required:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"urls": {
						Description: "Specifies basic auth username",
						Type:        types.ListType{ElemType: types.StringType},
						Required:    true,
					},
					"labels": {
						Description: "Specifies basic auth password",
						Type:        types.MapType{ElemType: types.StringType},
						Optional:    true,
					},
				}),
			},
		},
	}, nil
}
