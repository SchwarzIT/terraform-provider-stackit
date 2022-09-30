package instance

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

// Instance is the schema model
type Instance struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	ProjectID                   types.String `tfsdk:"project_id"`
	Plan                        types.String `tfsdk:"plan"`
	Grafana                     *Grafana     `tfsdk:"grafana"`
	Metrics                     *Metrics     `tfsdk:"metrics"`
	PlanID                      types.String `tfsdk:"plan_id"`
	DashboardURL                types.String `tfsdk:"dashboard_url"`
	IsUpdatable                 types.Bool   `tfsdk:"is_updatable"`
	GrafanaURL                  types.String `tfsdk:"grafana_url"`
	GrafanaInitialAdminPassword types.String `tfsdk:"grafana_initial_admin_password"`
	GrafanaInitialAdminUser     types.String `tfsdk:"grafana_initial_admin_user"`
	MetricsURL                  types.String `tfsdk:"metrics_url"`
	MetricsPushURL              types.String `tfsdk:"metrics_push_url"`
	TargetsURL                  types.String `tfsdk:"targets_url"`
	AlertingURL                 types.String `tfsdk:"alerting_url"`
	LogsURL                     types.String `tfsdk:"logs_url"`
	LogsPushURL                 types.String `tfsdk:"logs_push_url"`
	JaegerTracesURL             types.String `tfsdk:"jaeger_traces_url"`
	JaegerUIURL                 types.String `tfsdk:"jaeger_ui_url"`
	OtlpTracesURL               types.String `tfsdk:"otlp_traces_url"`
	ZipkinSpansURL              types.String `tfsdk:"zipkin_spans_url"`
}

type Grafana struct {
	EnablePublicAccess types.Bool `tfsdk:"enable_public_access"`
}

type Metrics struct {
	RetentionDays               types.Int64 `tfsdk:"retention_days"`
	RetentionDays5mDownsampling types.Int64 `tfsdk:"retention_days_5m_downsampling"`
	RetentionDays1hDownsampling types.Int64 `tfsdk:"retention_days_1h_downsampling"`
}

// GetSchema returns the terraform schema structure
func (r Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages Argus Instances",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the Argus instance ID",
				Type:        types.StringType,
				Computed:    true,
			},

			"name": {
				Description: "Specifies the name of the Argus instance",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
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
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},

			"plan": {
				Description: "Specifies the Argus plan. Available options are: `Monitoring-Medium-EU01`, `Monitoring-Large-EU01`, `Frontend-Starter-EU01`, `Monitoring-XL-EU01`, `Monitoring-XXL-EU01`, `Monitoring-Starter-EU01`, `Monitoring-Basic-EU01`, `Observability-Medium-EU01`, `Observability-Large-EU01 `, `Observability-XL-EU01`, `Observability-Starter-EU01`, `Observability-Basic-EU01`, `Observability-XXL-EU01`.",
				Type:        types.StringType,
				Required:    true,
			},

			"grafana": {
				Description: "A Grafana configuration block",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enable_public_access": {
						Description: "If true, anyone can access Grafana dashboards without logging in. Default is set to `false`.",
						Type:        types.BoolType,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.BoolDefault(default_grafana_enable_public_access),
						},
					},
				}),
			},

			"metrics": {
				Description: "Metrics configuration block",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"retention_days": {
						Description: "Specifies for how many days the raw metrics are kept. Default is set to `90`",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_metrics_retention_days),
						},
					},
					"retention_days_5m_downsampling": {
						Description: "Specifies for how many days the 5m downsampled metrics are kept. must be less than the value of the general retention. Default is set to `0` (disabled).",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_metrics_retention_days_5m_downsampling),
						},
					},
					"retention_days_1h_downsampling": {
						Description: "Specifies for how many days the 1h downsampled metrics are kept. must be less than the value of the 5m downsampling retention. Default is set to `0` (disabled).",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_metrics_retention_days_1h_downsampling),
						},
					},
				}),
			},

			// Read only:

			"plan_id": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Argus Plan ID.",
			},

			"dashboard_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Argus instance dashboard URL.",
			},

			"is_updatable": {
				Type:        types.BoolType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies if the instance can be updated.",
			},

			"grafana_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Grafana URL.",
			},

			"grafana_initial_admin_password": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Sensitive:   true,
				Description: "Specifies an initial Grafana admin password.",
			},

			"grafana_initial_admin_user": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies an initial Grafana admin username.",
			},

			"metrics_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies metrics URL.",
			},

			"metrics_push_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies URL for pushing metrics.",
			},

			"targets_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Targets URL.",
			},

			"alerting_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Alerting URL.",
			},

			"logs_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies Logs URL.",
			},

			"logs_push_url": {
				Type:        types.StringType,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "Specifies URL for pushing logs.",
			},

			"jaeger_traces_url": {
				Type:     types.StringType,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"jaeger_ui_url": {
				Type:     types.StringType,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"otlp_traces_url": {
				Type:     types.StringType,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"zipkin_spans_url": {
				Type:     types.StringType,
				Computed: true,
				Required: false,
				Optional: false,
			},
		},
	}, nil
}
