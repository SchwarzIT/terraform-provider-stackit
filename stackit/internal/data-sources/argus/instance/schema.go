package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for Argus Instances\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the Argus instance ID.",
				Required:    true,
			},

			"project_id": schema.StringAttribute{
				Description: "Specifies the Project ID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			// Read only:

			"name": schema.StringAttribute{
				Description: "Specifies the name of the Argus instance.",
				Computed:    true,
			},

			"plan": schema.StringAttribute{
				Description: "Specifies the Argus plan.",
				Computed:    true,
			},

			"grafana": schema.SingleNestedAttribute{
				Description: "A Grafana configuration block.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"enable_public_access": schema.BoolAttribute{
						Description: "If true, anyone can access Grafana dashboards without logging in.",
						Computed:    true,
					},
				},
			},

			"metrics": schema.SingleNestedAttribute{
				Description: "Metrics configuration block",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"retention_days": schema.Int64Attribute{
						Description: "Specifies for how many days the raw metrics are kept.",
						Computed:    true,
					},
					"retention_days_5m_downsampling": schema.Int64Attribute{
						Description: "Specifies for how many days the 5m downsampled metrics are kept.",
						Computed:    true,
					},
					"retention_days_1h_downsampling": schema.Int64Attribute{
						Description: "Specifies for how many days the 1h downsampled metrics are kept.",
						Computed:    true,
					},
				},
			},

			"plan_id": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies Argus Plan ID.",
			},

			"dashboard_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies Argus instance dashboard URL.",
			},

			"is_updatable": schema.BoolAttribute{
				Computed:    true,
				Description: "Specifies if the instance can be updated.",
			},

			"grafana_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies Grafana URL.",
			},

			"grafana_initial_admin_password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Specifies an initial Grafana admin password.",
			},

			"grafana_initial_admin_user": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies an initial Grafana admin username.",
			},

			"metrics_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies metrics URL.",
			},

			"metrics_push_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies URL for pushing metrics.",
			},

			"targets_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies Targets URL.",
			},

			"alerting_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies Alerting URL.",
			},

			"logs_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies Logs URL.",
			},

			"logs_push_url": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies URL for pushing logs.",
			},

			"jaeger_traces_url": schema.StringAttribute{
				Computed: true,
			},

			"jaeger_ui_url": schema.StringAttribute{
				Computed: true,
			},

			"otlp_traces_url": schema.StringAttribute{
				Computed: true,
			},

			"zipkin_spans_url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
