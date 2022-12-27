package instance

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	grafanaConfigs "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/generated/grafana-configs"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/generated/instances"
	metricsStorageRetention "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/generated/metrics-storage-retention"
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
	c := r.client.Services.Argus
	n := plan.Name.ValueString()
	pa := map[string]interface{}{}
	body := instances.InstanceCreateJSONRequestBody{
		Name:      &n,
		PlanID:    plan.PlanID.ValueString(),
		Parameter: &pa,
	}
	res, err := c.Instances.InstanceCreateWithResponse(ctx, plan.ProjectID.ValueString(), body)
	if err != nil {
		diags.AddError("failed preparing instance creation request", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed making instance creation request", res.HasError.Error())
		return
	}
	if res.JSON202 == nil {
		diags.AddError("got an empty response", "JSON202 == nil")
		return
	}

	process := res.WaitHandler(ctx, c.Instances, plan.ProjectID.ValueString(), res.JSON202.InstanceID)
	wr, err := process.Wait()
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

	c := r.client.Services
	epa := s.Grafana.EnablePublicAccess.ValueBool()
	cfg := grafanaConfigs.UpdateJSONRequestBody{
		PublicReadAccess: &epa,
	}

	res, err := c.Argus.GrafanaConfigs.UpdateWithResponse(ctx, s.ProjectID.ValueString(), s.ID.ValueString(), cfg)
	if err != nil {
		diags.AddError("failed to prepare grafana config request", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed to make grafana config request", res.HasError.Error())
		return
	}
}

func (r Resource) setMetricsConfig(ctx context.Context, diags *diag.Diagnostics, s *Instance, ref *Instance) {
	if s.Metrics == nil && ref == nil {
		return
	}

	if ref != nil && ref.Metrics != nil {
		if s.Metrics == nil {
			s.Metrics = &Metrics{
				RetentionDays:               types.Int64Value(default_metrics_retention_days),
				RetentionDays5mDownsampling: types.Int64Value(default_metrics_retention_days_5m_downsampling),
				RetentionDays1hDownsampling: types.Int64Value(default_metrics_retention_days_1h_downsampling),
			}
		} else if ref.Metrics.RetentionDays.Equal(s.Metrics.RetentionDays) &&
			ref.Metrics.RetentionDays5mDownsampling.Equal(s.Metrics.RetentionDays5mDownsampling) &&
			ref.Metrics.RetentionDays1hDownsampling.Equal(s.Metrics.RetentionDays1hDownsampling) {
			return
		}
	}

	c := r.client.Services
	cfg := metricsStorageRetention.UpdateJSONRequestBody{
		MetricsRetentionTimeRaw: fmt.Sprintf("%dd", s.Metrics.RetentionDays.ValueInt64()),
		MetricsRetentionTime5m:  fmt.Sprintf("%dd", s.Metrics.RetentionDays5mDownsampling.ValueInt64()),
		MetricsRetentionTime1h:  fmt.Sprintf("%dd", s.Metrics.RetentionDays1hDownsampling.ValueInt64()),
	}

	res, err := c.Argus.MetricsStorageRetention.UpdateWithResponse(ctx, s.ProjectID.ValueString(), s.ID.ValueString(), cfg)
	if err != nil {
		diags.AddError("failed to prepare metrics config request", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed to make metrics config request", res.HasError.Error())
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
	c := r.client.Services
	res, err := c.Argus.Instances.InstanceReadWithResponse(ctx, s.ProjectID.ValueString(), s.ID.ValueString())
	if err != nil {
		diags.AddError("failed to prepare read instance request", err.Error())
		return
	}
	if res.HasError != nil {
		if res.StatusCode() == http.StatusNotFound {
			s.ID = types.StringValue("")
			return
		}
		diags.AddError("failed to read instance", res.HasError.Error())
		return
	}
	if res.JSON200 == nil {
		diags.AddError("read instance returned an empty response", "JSON200 == nil")
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

	c := r.client.Services
	res, err := c.Argus.GrafanaConfigs.ListWithResponse(ctx, s.ProjectID.ValueString(), s.ID.ValueString())
	if err != nil {
		diags.AddError("failed to prepare read grafana configs request", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed to make read grafana configs request", res.HasError.Error())
		return
	}
	if res.JSON200 == nil {
		diags.AddError("read grafana configs returned an empty response", "JSON200 == nil")
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

	c := r.client.Services
	res, err := c.Argus.MetricsStorageRetention.ListWithResponse(ctx, s.ProjectID.ValueString(), s.ID.ValueString())
	if err != nil {
		diags.AddError("failed to prepare read metrics storage retention request", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed to make read metrics storage retention request", res.HasError.Error())
		return
	}
	if res.JSON200 == nil {
		diags.AddError("read metrics storage retention returned an empty response", "JSON200 == nil")
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

	c := r.client.Services
	n := plan.Name.ValueString()
	p := map[string]interface{}{}
	body := instances.InstanceUpdateJSONRequestBody{
		Name:      &n,
		Parameter: &p,
		PlanID:    plan.PlanID.ValueString(),
	}
	res, err := c.Argus.Instances.InstanceUpdateWithResponse(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		diags.AddError("failed preparing instance update request", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed during instance update", res.HasError.Error())
		return
	}
	if res.JSON202 == nil {
		diags.AddError("read instance returned an empty response", "JSON202 == nil")
		return
	}
	process := res.WaitHandler(ctx, c.Argus.Instances, plan.ProjectID.ValueString(), plan.ID.ValueString())
	process.SetTimeout(2 * time.Hour)
	wr, err := process.Wait()
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

	c := r.client.Services
	res, err := c.Argus.Instances.InstanceDeleteWithResponse(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed preparing instance delete request", err.Error())
		return
	}
	if res.HasError != nil {
		resp.Diagnostics.AddError("failed during instance delete", res.HasError.Error())
		return
	}
	process := res.WaitHandler(ctx, r.client.Services.Argus.Instances, state.ProjectID.ValueString(), state.ID.ValueString())
	if _, err := process.Wait(); err != nil {
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
