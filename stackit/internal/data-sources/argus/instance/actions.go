package instance

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/instance"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Instance

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.Argus.Instances.Get(ctx, config.ProjectID.ValueString(), config.ID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed instance read", agg.Error())
		return
	}
	b := res.JSON200
	config.ID = types.StringValue(b.ID)
	config.Name = types.StringNull()
	if b.Name != nil {
		config.Name = types.StringValue(*b.Name)
	}
	config.Plan = types.StringValue(b.PlanName)
	config.PlanID = types.StringValue(b.PlanID)
	config.DashboardURL = types.StringValue(b.DashboardURL)
	config.IsUpdatable = types.BoolNull()
	if b.IsUpdatable != nil {
		config.IsUpdatable = types.BoolValue(*b.IsUpdatable)
	}
	config.GrafanaURL = types.StringValue(b.Instance.GrafanaURL)
	config.GrafanaInitialAdminPassword = types.StringValue(b.Instance.GrafanaAdminPassword)
	config.GrafanaInitialAdminUser = types.StringValue(b.Instance.GrafanaAdminUser)
	config.MetricsURL = types.StringValue(b.Instance.MetricsURL)
	config.MetricsPushURL = types.StringValue(b.Instance.PushMetricsURL)
	config.TargetsURL = types.StringValue(b.Instance.TargetsURL)
	config.AlertingURL = types.StringValue(b.Instance.AlertingURL)
	config.LogsURL = types.StringValue(b.Instance.LogsURL)
	config.LogsPushURL = types.StringValue(b.Instance.LogsPushURL)
	config.JaegerTracesURL = types.StringValue(b.Instance.JaegerTracesURL)
	config.JaegerUIURL = types.StringValue(b.Instance.JaegerUiURL)
	config.OtlpTracesURL = types.StringValue(b.Instance.OtlpTracesURL)
	config.ZipkinSpansURL = types.StringValue(b.Instance.ZipkinSpansURL)
	config.Grafana = &instance.Grafana{
		EnablePublicAccess: types.BoolValue(b.Instance.GrafanaPublicReadAccess),
	}
	config.Metrics = &instance.Metrics{
		RetentionDays:               types.Int64Value(int64(b.Instance.MetricsRetentionTimeRaw)),
		RetentionDays1hDownsampling: types.Int64Value(int64(b.Instance.MetricsRetentionTime1h)),
		RetentionDays5mDownsampling: types.Int64Value(int64(b.Instance.MetricsRetentionTime5m)),
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
