package instance

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/instances"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	DefaultGrafanaEnablePublicAccess          bool  = false
	DefaultMetricsRetentionDays               int64 = 90
	DefaultMetricsRetentionDays5mDownsampling int64 = 0
	DefaultMetricsRetentionDays1hDownsampling int64 = 0
)

func (r Resource) loadPlanID(ctx context.Context, diags *diag.Diagnostics, s *Instance) {
	c := r.client.Argus

	res, err := c.Plans.ListPlans(ctx, s.ProjectID.ValueString())
	if agg := common.Validate(diags, res, err, "JSON200"); agg != nil {
		diags.AddError("failed to list argus plans", agg.Error())
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
		s.IsUpdatable = types.BoolValue(*res.IsUpdatable)
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
