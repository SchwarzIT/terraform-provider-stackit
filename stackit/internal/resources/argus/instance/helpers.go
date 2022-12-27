package instance

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/generated/instances"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	default_grafana_enable_public_access                 = false
	default_metrics_retention_days                 int64 = 90
	default_metrics_retention_days_5m_downsampling int64 = 0
	default_metrics_retention_days_1h_downsampling int64 = 0
)

func (r Resource) loadPlanID(ctx context.Context, diags *diag.Diagnostics, s *Instance) {
	c := r.client.Services.Argus

	res, err := c.Plans.ListPlansWithResponse(ctx, s.ProjectID.ValueString())
	if err != nil {
		diags.AddError("failed to prepare list plans request", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed to make list plans request", res.HasError.Error())
		return
	}
	if res.JSON200 == nil {
		diags.AddError("received an empty response", "JSON200 == nil")
		return
	}

	avl := ""
	for _, v := range res.JSON200.Plans {
		if v.Name == nil {
			continue
		}
		if *v.Name == s.Plan.ValueString() {
			s.PlanID = types.StringValue(v.PlanID.String())
			break
		}
		avl = fmt.Sprintf("%s\n- %s", avl, *v.Name)
	}
	if s.PlanID.ValueString() == "" {
		diags.AddError("invalid plan", fmt.Sprintf("couldn't find plan '%s'.\navailable names are:%s", s.Plan.ValueString(), avl))
		return
	}
}

func (l Instance) isEqual(got instances.ProjectInstanceUI) bool {
	if got.Name != nil && l.Name.ValueString() == *got.Name &&
		l.Plan.ValueString() == got.PlanName &&
		l.PlanID.ValueString() == got.PlanID {
		return true
	}
	return false
}

func updateByAPIResult(s *Instance, res *instances.ProjectInstanceUI) {
	s.ID = types.StringValue(res.ID)
	s.Plan = types.StringValue(res.PlanName)
	s.PlanID = types.StringValue(res.PlanID)
	if res.Name != nil {
		s.Name = types.StringValue(*res.Name)
	}
	s.DashboardURL = types.StringValue(res.DashboardURL)
	if res.IsUpdatable != nil {
		s.IsUpdatable = types.Bool{Value: *res.IsUpdatable}
	}
	s.GrafanaURL = types.StringValue(res.Instance.GrafanaURL)
	s.GrafanaInitialAdminPassword = types.StringValue(res.Instance.GrafanaAdminPassword)
	s.GrafanaInitialAdminUser = types.StringValue(res.Instance.GrafanaAdminUser)
	s.MetricsURL = types.StringValue(res.Instance.MetricsURL)
	s.MetricsPushURL = types.StringValue(res.Instance.PushMetricsURL)
	s.TargetsURL = types.StringValue(res.Instance.TargetsURL)
	s.AlertingURL = types.StringValue(res.Instance.AlertingURL)
	s.LogsURL = types.StringValue(res.Instance.LogsURL)
	s.LogsPushURL = types.StringValue(res.Instance.LogsPushURL)
	s.JaegerTracesURL = types.StringValue(res.Instance.JaegerTracesURL)
	s.JaegerUIURL = types.StringValue(res.Instance.JaegerUiURL)
	s.OtlpTracesURL = types.StringValue(res.Instance.OtlpTracesURL)
	s.ZipkinSpansURL = types.StringValue(res.Instance.ZipkinSpansURL)
}

func transformDayMetric(days interface{}) int64 {
	t := strings.TrimSuffix(days.(string), "d")
	if t == "" {
		t = "0"
	}
	r, _ := strconv.Atoi(t)
	return int64(r)
}
