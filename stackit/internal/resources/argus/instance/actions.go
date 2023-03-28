package instance

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	grafanaConfigs "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/grafana-configs"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/instances"
	metricsStorageRetention "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/metrics-storage-retention"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("cannot get plan", "failed trying to get plan")
		return
	}

	r.loadPlanID(ctx, &resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	r.createInstance(ctx, &resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// artificial wait for instance to be ready
	time.Sleep(1 * time.Minute)

	r.setGrafanaConfig(ctx, &resp.Diagnostics, &plan, nil)
	if resp.Diagnostics.HasError() {
		return
	}

	r.setMetricsConfig(ctx, &resp.Diagnostics, &plan, nil)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createInstance(ctx context.Context, diags *diag.Diagnostics, plan *Instance) {
	c := r.client.Argus
	n := plan.Name.ValueString()
	pa := map[string]interface{}{}
	body := instances.CreateJSONRequestBody{
		Name:      &n,
		PlanID:    plan.PlanID.ValueString(),
		Parameter: &pa,
	}
	res, err := c.Instances.Create(ctx, plan.ProjectID.ValueString(), body)
	if agg := validate.Response(res, err, "JSON202"); agg != nil {
		diags.AddError("failed to create argus instance", agg.Error())
		if res != nil && res.JSON400 != nil {
			diags.AddError("bad request", fmt.Sprintf("%+v", *res.JSON400))
		}
		return
	}

	process := res.WaitHandler(ctx, c.Instances, plan.ProjectID.ValueString(), res.JSON202.InstanceID).SetTimeout(1 * time.Hour)
	wr, err := process.WaitWithContext(ctx)
	if err != nil {
		diags.AddError("failed validating instance creation", err.Error())
		return
	}

	got, ok := wr.(*instances.ProjectInstanceUI)
	if !ok {
		diags.AddError("failed wait result conversion", "result is not of *instances.ProjectInstanceUI")
		return
	}

	updateByAPIResult(plan, got)
}

func (r Resource) setGrafanaConfig(ctx context.Context, diags *diag.Diagnostics, s *Instance, ref *Instance) {
	if s.Grafana == nil && ref == nil {
		return
	}

	if ref != nil && ref.Grafana != nil {
		if s.Grafana == nil {
			s.Grafana = &Grafana{EnablePublicAccess: types.BoolValue(default_grafana_enable_public_access)}
		} else if ref.Grafana.EnablePublicAccess.Equal(s.Grafana.EnablePublicAccess) {
			return
		}
	}

	c := r.client
	cfg := grafanaConfigs.UpdateJSONRequestBody{}
	if s.Grafana != nil {
		epa := s.Grafana.EnablePublicAccess.ValueBool()
		cfg.PublicReadAccess = &epa
	}

	res, err := c.Argus.GrafanaConfigs.Update(ctx, s.ProjectID.ValueString(), s.ID.ValueString(), cfg)
	if agg := validate.Response(res, err); agg != nil {
		diags.AddError("failed to make grafana config request", agg.Error())
		return
	}
}

func (r Resource) setMetricsConfig(ctx context.Context, diags *diag.Diagnostics, s *Instance, ref *Instance) {
	if s.Metrics == nil && ref == nil {
		return
	}
	m := s.Metrics
	if m == nil {
		m = &Metrics{
			RetentionDays:               types.Int64Value(default_metrics_retention_days),
			RetentionDays5mDownsampling: types.Int64Value(default_metrics_retention_days_5m_downsampling),
			RetentionDays1hDownsampling: types.Int64Value(default_metrics_retention_days_1h_downsampling),
		}
	}
	if ref != nil && ref.Metrics != nil {
		if ref.Metrics.RetentionDays.Equal(m.RetentionDays) &&
			ref.Metrics.RetentionDays5mDownsampling.Equal(m.RetentionDays5mDownsampling) &&
			ref.Metrics.RetentionDays1hDownsampling.Equal(m.RetentionDays1hDownsampling) {
			return
		}
	}
	cfg := metricsStorageRetention.UpdateJSONRequestBody{
		MetricsRetentionTimeRaw: fmt.Sprintf("%dd", m.RetentionDays.ValueInt64()),
		MetricsRetentionTime5m:  fmt.Sprintf("%dd", m.RetentionDays5mDownsampling.ValueInt64()),
		MetricsRetentionTime1h:  fmt.Sprintf("%dd", m.RetentionDays1hDownsampling.ValueInt64()),
	}
	res, err := r.client.Argus.MetricsStorageRetention.Update(ctx, s.ProjectID.ValueString(), s.ID.ValueString(), cfg)
	if agg := validate.Response(res, err); agg != nil {
		diags.AddError("failed to make metrics config request", agg.Error())
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Instance

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readInstance(ctx, &resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	r.readGrafana(ctx, &resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readMetrics(ctx, &resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) readInstance(ctx context.Context, diags *diag.Diagnostics, s *Instance) {
	c := r.client
	res, err := c.Argus.Instances.Get(ctx, s.ProjectID.ValueString(), s.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			s.ID = types.StringValue("")
			return
		}
		diags.AddError("failed to read instance", agg.Error())
		return
	}
	updateByAPIResult(s, res.JSON200)
}

func (r Resource) readGrafana(ctx context.Context, diags *diag.Diagnostics, s *Instance) {
	if s.Grafana == nil {
		return
	}
	if s.ID.ValueString() == "" {
		diags.AddError("missing instance ID", "not instance ID specified when reading grafana config")
		return
	}

	c := r.client
	res, err := c.Argus.GrafanaConfigs.List(ctx, s.ProjectID.ValueString(), s.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		diags.AddError("failed to read grafana configs", agg.Error())
		return
	}

	s.Grafana.EnablePublicAccess = types.BoolValue(*res.JSON200.PublicReadAccess)
}

func (r Resource) readMetrics(ctx context.Context, diags *diag.Diagnostics, s *Instance) {
	if s.Metrics == nil {
		return
	}
	if s.ID.ValueString() == "" {
		diags.AddError("missing instance ID", "not instance ID specified when reading metrics config")
		return
	}

	c := r.client
	res, err := c.Argus.MetricsStorageRetention.List(ctx, s.ProjectID.ValueString(), s.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		diags.AddError("failed to read metrics storage retention", agg.Error())
		return
	}
	s.Metrics.RetentionDays = types.Int64Value(transformDayMetric(res.JSON200.MetricsRetentionTimeRaw))
	s.Metrics.RetentionDays5mDownsampling = types.Int64Value(transformDayMetric(res.JSON200.MetricsRetentionTime5m))
	s.Metrics.RetentionDays1hDownsampling = types.Int64Value(transformDayMetric(res.JSON200.MetricsRetentionTime1h))
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state Instance
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// check if computed attributes are in the plan
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	if plan.PlanID.ValueString() == "" {
		if plan.Plan.Equal(state.Plan) {
			plan.PlanID = state.PlanID
		} else {
			r.loadPlanID(ctx, &resp.Diagnostics, &plan)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// update corresponding APIs
	r.setGrafanaConfig(ctx, &resp.Diagnostics, &plan, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	r.setMetricsConfig(ctx, &resp.Diagnostics, &plan, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	// update using instance API if needed
	r.updateInstance(ctx, &resp.Diagnostics, &plan, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readInstance(ctx, &resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state with plan
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) updateInstance(ctx context.Context, diags *diag.Diagnostics, plan, ref *Instance) {
	// skip update if there's nothing to update
	if ref != nil &&
		plan.Name.Equal(ref.Name) &&
		plan.Plan.Equal(ref.Plan) &&
		plan.PlanID.Equal(ref.PlanID) {
		return
	}

	c := r.client
	n := plan.Name.ValueString()
	p := map[string]interface{}{}
	body := instances.UpdateJSONRequestBody{
		Name:      &n,
		Parameter: &p,
		PlanID:    plan.PlanID.ValueString(),
	}
	res, err := c.Argus.Instances.Update(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), body)
	if agg := validate.Response(res, err, "JSON202"); agg != nil {
		diags.AddError("failed to update instance", agg.Error())
		return
	}
	process := res.WaitHandler(ctx, c.Argus.Instances, plan.ProjectID.ValueString(), plan.ID.ValueString())
	process.SetTimeout(2 * time.Hour)
	wr, err := process.WaitWithContext(ctx)
	if err != nil {
		diags.AddError("failed validating instance update", err.Error())
		return
	}

	got, ok := wr.(*instances.ProjectInstanceUI)
	if !ok || got == nil {
		diags.AddError("failed wait result conversion", "response is not of *instances.ProjectInstanceUI or nil")
		return
	}

	if !plan.isEqual(*got) {
		updateByAPIResult(plan, got)
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Instance
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() {
		resp.Diagnostics.AddError("can't perform deletion", "argus instance id is unknown or null")
	}

	c := r.client
	res, err := c.Argus.Instances.Delete(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		resp.Diagnostics.AddError("failed to delete instance", agg.Error())
		return
	}
	process := res.WaitHandler(ctx, r.client.Argus.Instances, state.ProjectID.ValueString(), state.ID.ValueString())
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed verify instance deletion", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,id` where `id` is the instance id.\nInstead got: %q", req.ID),
		)
		return
	}
	projectID := idParts[0]
	instanceID := idParts[1]

	// validate project id
	if err := clientValidate.ProjectID(projectID); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), instanceID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// pre-read imports
	inst := Instance{
		ID:        types.StringValue(instanceID),
		ProjectID: types.StringValue(projectID),
		Grafana:   &Grafana{},
		Metrics:   &Metrics{},
	}

	r.readGrafana(ctx, &resp.Diagnostics, &inst)
	if inst.Grafana.EnablePublicAccess.ValueBool() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("grafana"), &Grafana{
			EnablePublicAccess: types.BoolValue(true),
		})...)
	}

	r.readMetrics(ctx, &resp.Diagnostics, &inst)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metrics"), &Metrics{
		RetentionDays:               inst.Metrics.RetentionDays,
		RetentionDays5mDownsampling: inst.Metrics.RetentionDays5mDownsampling,
		RetentionDays1hDownsampling: inst.Metrics.RetentionDays1hDownsampling,
	})...)

}
