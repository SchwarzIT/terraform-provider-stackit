package instance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/grafana"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/instances"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/metrics"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.Provider.IsConfigured() {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on another resource.",
		)
		return
	}

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
	c := r.Provider.Client()
	created := false
	var got instances.Instance
	if err := helper.RetryContext(ctx, common.DURATION_1H, func() *helper.RetryError {
		if !created {
			res, err := c.Argus.Instances.Create(ctx, plan.ProjectID.Value, plan.Name.Value, plan.PlanID.Value, map[string]string{})
			if err != nil {
				if common.IsNonRetryable(err) {
					return helper.NonRetryableError(err)
				}
				return helper.RetryableError(err)
			}
			created = true
			plan.ID = types.String{Value: res.InstanceID}
		}

		if plan.ID.Value == "" {
			return helper.NonRetryableError(errors.New("received empty id"))
		}

		var err error
		got, err = c.Argus.Instances.Get(ctx, plan.ProjectID.Value, plan.ID.Value)
		if err != nil {
			return helper.RetryableError(err)
		}

		switch got.Status {
		case consts.ARGUS_INSTANCE_STATUS_UPDATE_SUCCEEDED:
			fallthrough
		case consts.ARGUS_INSTANCE_STATUS_CREATE_SUCCEEDED:
			return nil
		case consts.ARGUS_INSTANCE_STATUS_CREATING:
			return helper.RetryableError(fmt.Errorf("received status %s", got.Status))
		}

		// fail for any other status
		return helper.NonRetryableError(fmt.Errorf("received status %s", got.Status))
	}); err != nil {
		diags.AddError("failed to verify instance creation", err.Error())
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
			s.Grafana = &Grafana{EnablePublicAccess: types.Bool{Value: default_grafana_enable_public_access}}
		} else if ref.Grafana.EnablePublicAccess.Equal(s.Grafana.EnablePublicAccess) {
			return
		}
	}

	c := r.Provider.Client()
	cfg := grafana.Config{
		PublicReadAccess: s.Grafana.EnablePublicAccess.Value,
	}

	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		_, err := c.Argus.Grafana.UpdateConfig(ctx, s.ProjectID.Value, s.ID.Value, cfg)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		diags.AddError("failed to set grafana config", err.Error())
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
				RetentionDays:               types.Int64{Value: default_metrics_retention_days},
				RetentionDays5mDownsampling: types.Int64{Value: default_metrics_retention_days_5m_downsampling},
				RetentionDays1hDownsampling: types.Int64{Value: default_metrics_retention_days_1h_downsampling},
			}
		} else if ref.Metrics.RetentionDays.Equal(s.Metrics.RetentionDays) &&
			ref.Metrics.RetentionDays5mDownsampling.Equal(s.Metrics.RetentionDays5mDownsampling) &&
			ref.Metrics.RetentionDays1hDownsampling.Equal(s.Metrics.RetentionDays1hDownsampling) {
			return
		}
	}

	c := r.Provider.Client()
	cfg := metrics.Config{
		MetricsRetentionTimeRaw: fmt.Sprintf("%dd", s.Metrics.RetentionDays.Value),
		MetricsRetentionTime5m:  fmt.Sprintf("%dd", s.Metrics.RetentionDays5mDownsampling.Value),
		MetricsRetentionTime1h:  fmt.Sprintf("%dd", s.Metrics.RetentionDays1hDownsampling.Value),
	}

	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		_, err := c.Argus.Metrics.UpdateConfig(ctx, s.ProjectID.Value, s.ID.Value, cfg)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		diags.AddError("failed to set metrics config", err.Error())
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

	if state.ID.Value == "" {
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
	c := r.Provider.Client()
	var res instances.Instance
	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		var err error
		res, err = c.Argus.Instances.Get(ctx, s.ProjectID.Value, s.ID.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				s.ID = types.String{Value: ""}
				return nil
			}
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		diags.AddError("failed to read instance", err.Error())
		return
	}

	updateByAPIResult(s, res)
}

func (r Resource) readGrafana(ctx context.Context, diags *diag.Diagnostics, s *Instance) {
	if s.Grafana == nil {
		return
	}
	if s.ID.Value == "" {
		diags.AddError("missing instance ID", "not instance ID specified when reading grafana config")
		return
	}

	c := r.Provider.Client()
	var res grafana.GetConfigResponse
	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		var err error
		res, err = c.Argus.Grafana.GetConfig(ctx, s.ProjectID.Value, s.ID.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		diags.AddError("failed to read grafana config", err.Error())
		return
	}

	s.Grafana.EnablePublicAccess = types.Bool{Value: res.PublicReadAccess}
}

func (r Resource) readMetrics(ctx context.Context, diags *diag.Diagnostics, s *Instance) {
	if s.Metrics == nil {
		return
	}
	if s.ID.Value == "" {
		diags.AddError("missing instance ID", "not instance ID specified when reading metrics config")
		return
	}

	c := r.Provider.Client()
	var res metrics.GetConfigResponse
	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		var err error
		res, err = c.Argus.Metrics.GetConfig(ctx, s.ProjectID.Value, s.ID.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		diags.AddError("failed to read grafana config", err.Error())
		return
	}

	s.Metrics.RetentionDays = types.Int64{Value: transformDayMetric(res.MetricsRetentionTimeRaw)}
	s.Metrics.RetentionDays5mDownsampling = types.Int64{Value: transformDayMetric(res.MetricsRetentionTime5m)}
	s.Metrics.RetentionDays1hDownsampling = types.Int64{Value: transformDayMetric(res.MetricsRetentionTime1h)}
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
	if plan.ID.Value == "" {
		plan.ID = state.ID
	}

	if plan.PlanID.Value == "" {
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

	updated := types.Bool{Value: false}
	if err := helper.RetryContext(ctx, common.DURATION_20M, r.updateInstanceRetry(ctx, diags, plan, &updated)); err != nil {
		diags.AddError("failed during instance update", err.Error())
		return
	}
}

func (r Resource) updateInstanceRetry(ctx context.Context, diags *diag.Diagnostics, plan *Instance, updated *types.Bool) func() *helper.RetryError {
	c := r.Provider.Client()

	return func() *helper.RetryError {
		if !updated.Value {
			_, err := c.Argus.Instances.Update(ctx, plan.ProjectID.Value, plan.ID.Value, plan.Name.Value, plan.PlanID.Value, map[string]string{})
			if err != nil {
				if common.IsNonRetryable(err) {
					if strings.Contains(err.Error(), "instance is not in the right state") ||
						strings.Contains(err.Error(), "Could not update instance") {
						return helper.RetryableError(err)
					}
					return helper.NonRetryableError(err)
				}
				return helper.RetryableError(err)
			}
			updated.Value = true
		}

		// Give API time to start updating
		time.Sleep(5 * time.Second)

		got, err := c.Argus.Instances.Get(ctx, plan.ProjectID.Value, plan.ID.Value)
		if err != nil {
			return helper.RetryableError(err)
		}

		switch got.Status {
		case consts.ARGUS_INSTANCE_STATUS_UPDATE_SUCCEEDED:
			fallthrough
		case consts.ARGUS_INSTANCE_STATUS_CREATE_SUCCEEDED:
			if plan.isEqual(got) {
				return nil
			}
			updated.Value = false
			return helper.RetryableError(fmt.Errorf("received status %s but update didn't happen for Instance ID: %s", got.Status, plan.ID.Value))
		case consts.ARGUS_INSTANCE_STATUS_UPDATING:
			return helper.RetryableError(fmt.Errorf("received status %s", got.Status))
		}

		// fail for any other status
		return helper.NonRetryableError(fmt.Errorf("received status %s", got.Status))
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

	deleted := types.Bool{Value: false}
	if err := helper.RetryContext(ctx, common.DURATION_30M, r.delete(ctx, &deleted, &state)); err != nil {
		resp.Diagnostics.AddError("failed to delete instance", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *Resource) delete(ctx context.Context, deleted *types.Bool, state *Instance) func() *helper.RetryError {
	c := r.Provider.Client()

	return func() *helper.RetryError {
		if !deleted.Value {
			_, err := c.Argus.Instances.Delete(ctx, state.ProjectID.Value, state.ID.Value)
			if err != nil {
				if strings.Contains(err.Error(), "instance is not in the right state") ||
					strings.Contains(err.Error(), "Could not delete instance") {
					return helper.RetryableError(err)
				}
				if common.IsNonRetryable(err) {
					return helper.NonRetryableError(err)
				}
				return helper.RetryableError(err)
			}
			deleted.Value = true
		}

		// give API a couple of seconds to process
		time.Sleep(2 * time.Second)

		// verify deletion
		res, err := c.Argus.Instances.List(ctx, state.ProjectID.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return nil
			}
			return helper.RetryableError(err)
		}
		for _, i := range res.Instances {
			if i.Instance != state.ID.Value {
				continue
			}
			if i.Status == consts.ARGUS_INSTANCE_STATUS_DELETE_SUCCEEDED {
				return nil
			}
			if i.Status == consts.ARGUS_INSTANCE_STATUS_CREATE_SUCCEEDED {
				deleted.Value = false
			}
			return helper.RetryableError(fmt.Errorf("deletion incomplete. instance %s has status %s", i.Name, i.Status))

		}
		return nil
	}
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
		ID:        types.String{Value: instanceID},
		ProjectID: types.String{Value: projectID},
		Grafana:   &Grafana{},
		Metrics:   &Metrics{},
	}

	r.readGrafana(ctx, &resp.Diagnostics, &inst)
	if inst.Grafana.EnablePublicAccess.Value {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("grafana"), &Grafana{
			EnablePublicAccess: types.Bool{Value: true},
		})...)
	}

	r.readMetrics(ctx, &resp.Diagnostics, &inst)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metrics"), &Metrics{
		RetentionDays:               inst.Metrics.RetentionDays,
		RetentionDays5mDownsampling: inst.Metrics.RetentionDays5mDownsampling,
		RetentionDays1hDownsampling: inst.Metrics.RetentionDays1hDownsampling,
	})...)

}
