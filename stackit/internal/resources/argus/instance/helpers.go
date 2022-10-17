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
	c := r.Provider.Client()

	res, err := c.Argus.Plans.List(ctx, s.ProjectID.Value)
	if err != nil {
		diags.AddError("failed to list plans", err.Error())
		return
	}

	avl := ""
	for _, v := range res.Plans {
		if v.Name == s.Plan.Value {
			s.PlanID = types.String{Value: v.PlanID}
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
	s.ID = types.String{Value: res.ID}
	s.Plan = types.String{Value: res.PlanName}
	s.PlanID = types.String{Value: res.PlanID}
	s.Name = types.String{Value: res.Name}
	s.DashboardURL = types.String{Value: res.DashboardURL}
	s.IsUpdatable = types.Bool{Value: res.IsUpdatable}
	s.GrafanaURL = types.String{Value: res.Instance.GrafanaURL}
	s.GrafanaInitialAdminPassword = types.String{Value: res.Instance.GrafanaAdminPassword}
	s.GrafanaInitialAdminUser = types.String{Value: res.Instance.GrafanaAdminUser}
	s.MetricsURL = types.String{Value: res.Instance.MetricsURL}
	s.MetricsPushURL = types.String{Value: res.Instance.PushMetricsURL}
	s.TargetsURL = types.String{Value: res.Instance.TargetsURL}
	s.AlertingURL = types.String{Value: res.Instance.AlertingURL}
	s.LogsURL = types.String{Value: res.Instance.LogsURL}
	s.LogsPushURL = types.String{Value: res.Instance.LogsPushURL}
	s.JaegerTracesURL = types.String{Value: res.Instance.JaegerTracesURL}
	s.JaegerUIURL = types.String{Value: res.Instance.JaegerUIURL}
	s.OtlpTracesURL = types.String{Value: res.Instance.OtlpTracesURL}
	s.ZipkinSpansURL = types.String{Value: res.Instance.ZipkinSpansURL}
}

func transformDayMetric(days interface{}) int64 {
	t := strings.TrimSuffix(days.(string), "d")
	if t == "" {
		t = "0"
	}
	r, _ := strconv.Atoi(t)
	return int64(r)
}
