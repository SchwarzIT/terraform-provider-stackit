package instance

import (
	"context"
	"fmt"
	"regexp"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID                          types.String   `tfsdk:"id"`
	Name                        types.String   `tfsdk:"name"`
	ProjectID                   types.String   `tfsdk:"project_id"`
	Plan                        types.String   `tfsdk:"plan"`
	Grafana                     *Grafana       `tfsdk:"grafana"`
	Metrics                     *Metrics       `tfsdk:"metrics"`
	PlanID                      types.String   `tfsdk:"plan_id"`
	DashboardURL                types.String   `tfsdk:"dashboard_url"`
	IsUpdatable                 types.Bool     `tfsdk:"is_updatable"`
	GrafanaURL                  types.String   `tfsdk:"grafana_url"`
	GrafanaInitialAdminPassword types.String   `tfsdk:"grafana_initial_admin_password"`
	GrafanaInitialAdminUser     types.String   `tfsdk:"grafana_initial_admin_user"`
	MetricsURL                  types.String   `tfsdk:"metrics_url"`
	MetricsPushURL              types.String   `tfsdk:"metrics_push_url"`
	TargetsURL                  types.String   `tfsdk:"targets_url"`
	AlertingURL                 types.String   `tfsdk:"alerting_url"`
	LogsURL                     types.String   `tfsdk:"logs_url"`
	LogsPushURL                 types.String   `tfsdk:"logs_push_url"`
	JaegerTracesURL             types.String   `tfsdk:"jaeger_traces_url"`
	JaegerUIURL                 types.String   `tfsdk:"jaeger_ui_url"`
	OtlpTracesURL               types.String   `tfsdk:"otlp_traces_url"`
	ZipkinSpansURL              types.String   `tfsdk:"zipkin_spans_url"`
	Timeouts                    timeouts.Value `tfsdk:"timeouts"`
}

type Grafana struct {
	EnablePublicAccess types.Bool `tfsdk:"enable_public_access"`
}

type Metrics struct {
	RetentionDays               types.Int64 `tfsdk:"retention_days"`
	RetentionDays5mDownsampling types.Int64 `tfsdk:"retention_days_5m_downsampling"`
	RetentionDays1hDownsampling types.Int64 `tfsdk:"retention_days_1h_downsampling"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages Argus Instances\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the Argus instance ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				Description: "Specifies the name of the Argus instance",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9]+$`),
						"must contain only lowercase alphanumeric characters",
					),
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
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"plan": schema.StringAttribute{
				Description: "Specifies the Argus plan. Available options are: `Monitoring-Medium-EU01`, `Monitoring-Large-EU01`, `Frontend-Starter-EU01`, `Monitoring-XL-EU01`, `Monitoring-XXL-EU01`, `Monitoring-Starter-EU01`, `Monitoring-Basic-EU01`, `Observability-Medium-EU01`, `Observability-Large-EU01 `, `Observability-XL-EU01`, `Observability-Starter-EU01`, `Observability-Basic-EU01`, `Observability-XXL-EU01`.",
				Required:    true,
			},

			"grafana": schema.SingleNestedAttribute{
				Description: "A Grafana configuration block",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enable_public_access": schema.BoolAttribute{
						Description: "If true, anyone can access Grafana dashboards without logging in. Default is set to `false`.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(DefaultGrafanaEnablePublicAccess),
					},
				},
			},

			"metrics": schema.SingleNestedAttribute{
				Description: "Metrics configuration block",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"retention_days": schema.Int64Attribute{
						Description: "Specifies for how many days the raw metrics are kept. Default is set to `90`",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(DefaultMetricsRetentionDays),
					},
					"retention_days_5m_downsampling": schema.Int64Attribute{
						Description: "Specifies for how many days the 5m downsampled metrics are kept. must be less than the value of the general retention. Default is set to `0` (disabled).",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(DefaultMetricsRetentionDays5mDownsampling),
					},
					"retention_days_1h_downsampling": schema.Int64Attribute{
						Description: "Specifies for how many days the 1h downsampled metrics are kept. must be less than the value of the 5m downsampling retention. Default is set to `0` (disabled).",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(DefaultMetricsRetentionDays1hDownsampling),
					},
				},
			},

			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),

			// Read only:

			"plan_id": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Argus Plan ID.",
			},

			"dashboard_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Argus instance dashboard URL.",
			},

			"is_updatable": schema.BoolAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies if the instance can be updated.",
			},

			"grafana_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Grafana URL.",
			},

			"grafana_initial_admin_password": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Sensitive:   true,
				Description: "Specifies an initial Grafana admin password.",
			},

			"grafana_initial_admin_user": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies an initial Grafana admin username.",
			},

			"metrics_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies metrics URL.",
			},

			"metrics_push_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies URL for pushing metrics.",
			},

			"targets_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Targets URL.",
			},

			"alerting_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Alerting URL.",
			},

			"logs_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Logs URL.",
			},

			"logs_push_url": schema.StringAttribute{
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies URL for pushing logs.",
			},

			"jaeger_traces_url": schema.StringAttribute{
				Computed: true,
				Required: false,
				Optional: false,
			},

			"jaeger_ui_url": schema.StringAttribute{
				Computed: true,
				Required: false,
				Optional: false,
			},

			"otlp_traces_url": schema.StringAttribute{
				Computed: true,
				Required: false,
				Optional: false,
			},

			"zipkin_spans_url": schema.StringAttribute{
				Computed: true,
				Required: false,
				Optional: false,
			},
		},
	}
}
