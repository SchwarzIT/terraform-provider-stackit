package instance

import (
	"context"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/instances"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/instance"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"net/http"
	"strings"
)

func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.Provider.Client()
	var config instance.Instance
	var b instances.Instance

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		var err error
		b, err = c.Argus.Instances.Get(ctx, config.ProjectID.Value, config.ID.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read instance", err.Error())
		return
	}
	config.ID = types.String{Value: b.ID}
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
	config.Grafana.EnablePublicAccess = types.Bool{Value: b.Instance.GrafanaPublicReadAccess}
	config.Metrics.RetentionDays1hDownsampling = types.Int64{Value: int64(b.Instance.MetricsRetentionTime1h)}
	config.Metrics.RetentionDays5mDownsampling = types.Int64{Value: int64(b.Instance.MetricsRetentionTime5m)}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
