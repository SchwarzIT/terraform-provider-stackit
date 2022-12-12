package instance

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/instances"
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
	c := r.client

	res, err := c.Argus.Plans.List(ctx, s.ProjectID.Value)
	if err != nil {
		diags.AddError("failed to list plans", err.Error())
		return
	}

	avl := ""
	for _, v := range res.Plans {
		if v.Name == s.Plan.Value {
			s.PlanID = types.StringValue(v.PlanID)
			break
		}
		avl = fmt.Sprintf("%s\n- %s", avl, v.Name)
	}
	if s.PlanID.Value == "" {
		diags.AddError("invalid plan", fmt.Sprintf("couldn't find plan '%s'.\navailable names are:%s", s.Plan.Value, avl))
		return
	}
}

func (l Instance) isEqual(got instances.Instance) bool {
	if l.Name.Value == got.Name &&
		l.Plan.Value == got.PlanName &&
		l.PlanID.Value == got.PlanID {
		return true
	}
	return false
}

func updateByAPIResult(s *Instance, res instances.Instance) {
	s.ID = types.StringValue(res.ID)
	s.Plan = types.StringValue(res.PlanName)
	s.PlanID = types.StringValue(res.PlanID)
	s.Name = types.StringValue(res.Name)
	s.DashboardURL = types.StringValue(res.DashboardURL)
	s.IsUpdatable = types.Bool{Value: res.IsUpdatable}
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
	s.JaegerUIURL = types.StringValue(res.Instance.JaegerUIURL)
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
