package job

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages Argus Instance Jobs\n%s",
			common.EnvironmentInfo(r.urls),
		),
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"project_id": schema.StringAttribute{
				Description: "Specifies the Project ID the Argus instance belongs to",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"argus_instance_id": schema.StringAttribute{
				Description: "Specifies the Argus Instance ID the job belongs to",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"metrics_path": schema.StringAttribute{
				Description: "Specifies the job scraping path. Defaults to `/metrics`",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(DefaultMetricsPath),
			},

			"scheme": schema.StringAttribute{
				Description: "Specifies the scheme. Default is `https`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(DefaultScheme),
			},

			"scrape_interval": schema.StringAttribute{
				Description: "Specifies the scrape interval as duration string. Default is `5m`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(DefaultScrapeInterval),
			},

			"scrape_timeout": schema.StringAttribute{
				Description: "Specifies the scrape timeout as duration string. Default is `2m`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(DefaultScrapeTimeout),
			},

			"saml2": schema.SingleNestedAttribute{
				Description: "A saml2 configuration block",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enable_url_parameters": schema.BoolAttribute{
						Description: "Should URL parameters be enabled? Default is `true`",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(DefaultSAML2EnableURLParameters),
					},
				},
			},

			"basic_auth": schema.SingleNestedAttribute{
				Description: "A basic_auth block",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Description: "Specifies basic auth username",
						Required:    true,
					},
					"password": schema.StringAttribute{
						Description: "Specifies basic auth password",
						Required:    true,
					},
				},
			},

			"targets": schema.ListNestedAttribute{
				Description: "targets list",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"urls": schema.ListAttribute{
							Description: "Specifies target URLs",
							ElementType: types.StringType,
							Required:    true,
						},
						"labels": schema.MapAttribute{
							Description: "Specifies labels",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}
