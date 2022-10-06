package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Data source for Argus Instances",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the Argus instance ID",
				Type:        types.StringType,
				Required:    true,
			},

			"project_id": {
				Description: "Specifies the Project ID the Argus instance belongs to",
				Type:        types.StringType,
				Required:    true,
			},

			"name": {
				Description: "Specifies the name of the Argus instance",
				Type:        types.StringType,
				Computed:    true,
			},

			"plan": {
				Description: "Specifies the Argus plan. Available options are: `Monitoring-Medium-EU01`, `Monitoring-Large-EU01`, `Frontend-Starter-EU01`, `Monitoring-XL-EU01`, `Monitoring-XXL-EU01`, `Monitoring-Starter-EU01`, `Monitoring-Basic-EU01`, `Observability-Medium-EU01`, `Observability-Large-EU01 `, `Observability-XL-EU01`, `Observability-Starter-EU01`, `Observability-Basic-EU01`, `Observability-XXL-EU01`.",
				Type:        types.StringType,
				Computed:    true,
			},

			"grafana": {
				Description: "A Grafana configuration block",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enable_public_access": {
						Description: "If true, anyone can access Grafana dashboards without logging in. Default is set to `false`.",
						Type:        types.BoolType,
						Computed:    true,
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
						Computed:    true,
					},
					"retention_days_5m_downsampling": {
						Description: "Specifies for how many days the 5m downsampled metrics are kept. must be less than the value of the general retention. Default is set to `0` (disabled).",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"retention_days_1h_downsampling": {
						Description: "Specifies for how many days the 1h downsampled metrics are kept. must be less than the value of the 5m downsampling retention. Default is set to `0` (disabled).",
						Type:        types.Int64Type,
						Computed:    true,
					},
				}),
			},

			// Read only:

			"plan_id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies Argus Plan ID.",
			},

			"dashboard_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies Argus instance dashboard URL.",
			},

			"is_updatable": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Specifies if the instance can be updated.",
			},

			"grafana_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies Grafana URL.",
			},

			"grafana_initial_admin_password": {
				Type:        types.StringType,
				Computed:    true,
				Sensitive:   true,
				Description: "Specifies an initial Grafana admin password.",
			},

			"grafana_initial_admin_user": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies an initial Grafana admin username.",
			},

			"metrics_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies metrics URL.",
			},

			"metrics_push_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies URL for pushing metrics.",
			},

			"targets_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies Targets URL.",
			},

			"alerting_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies Alerting URL.",
			},

			"logs_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies Logs URL.",
			},

			"logs_push_url": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Specifies URL for pushing logs.",
			},

			"jaeger_traces_url": {
				Type:     types.StringType,
				Computed: true,
			},

			"jaeger_ui_url": {
				Type:     types.StringType,
				Computed: true,
			},

			"otlp_traces_url": {
				Type:     types.StringType,
				Computed: true,
			},

			"zipkin_spans_url": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}
