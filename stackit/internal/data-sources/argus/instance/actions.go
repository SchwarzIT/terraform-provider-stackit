package instance

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/instance"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config instance.Instance

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	b, err := c.Argus.Instances.Get(ctx, config.ProjectID.Value, config.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to read instance", err.Error())
		return
	}

	config.ID = types.String{Value: b.ID}
	config.Name = types.String{Value: b.Name}
	config.Plan = types.String{Value: b.PlanName}
	config.PlanID = types.String{Value: b.PlanID}
	config.DashboardURL = types.String{Value: b.DashboardURL}
	config.IsUpdatable = types.Bool{Value: b.IsUpdatable}
	config.GrafanaURL = types.String{Value: b.Instance.GrafanaURL}
	config.GrafanaInitialAdminPassword = types.String{Value: b.Instance.GrafanaAdminPassword}
	config.GrafanaInitialAdminUser = types.String{Value: b.Instance.GrafanaAdminUser}
	config.MetricsURL = types.String{Value: b.Instance.MetricsURL}
	config.MetricsPushURL = types.String{Value: b.Instance.PushMetricsURL}
	config.TargetsURL = types.String{Value: b.Instance.TargetsURL}
	config.AlertingURL = types.String{Value: b.Instance.AlertingURL}
	config.LogsURL = types.String{Value: b.Instance.LogsURL}
	config.LogsPushURL = types.String{Value: b.Instance.LogsPushURL}
	config.JaegerTracesURL = types.String{Value: b.Instance.JaegerTracesURL}
	config.JaegerUIURL = types.String{Value: b.Instance.JaegerUIURL}
	config.OtlpTracesURL = types.String{Value: b.Instance.OtlpTracesURL}
	config.ZipkinSpansURL = types.String{Value: b.Instance.ZipkinSpansURL}

	config.Grafana = &instance.Grafana{
		EnablePublicAccess: types.Bool{Value: b.Instance.GrafanaPublicReadAccess},
	}
	config.Metrics = &instance.Metrics{
		RetentionDays:               types.Int64{Value: int64(b.Instance.MetricsRetentionTimeRaw)},
		RetentionDays1hDownsampling: types.Int64{Value: int64(b.Instance.MetricsRetentionTime1h)},
		RetentionDays5mDownsampling: types.Int64{Value: int64(b.Instance.MetricsRetentionTime5m)},
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
