package project

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KubernetesProject
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	c := r.client.Kubernetes.Project
	res, err := c.Create(ctx, plan.ProjectID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if !strings.Contains(agg.Error(), common.ERR_UNEXPECTED_EOF) {
			resp.Diagnostics.AddError("failed during SKE project creation", agg.Error())
			return
		}
	}

	plan.ID = plan.ProjectID

	timeout, d := plan.Timeouts.Create(ctx, 10*time.Minute)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}
	process := res.WaitHandler(ctx, c, plan.ID.ValueString()).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed verifying SKE project creation", err.Error())
		return
	}

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KubernetesProject
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read
	c := r.client.Kubernetes.Project
	res, err := c.Get(ctx, state.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed during SKE project read", agg.Error())
		return
	}

}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KubernetesProject
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle deletion
	c := r.client.Kubernetes.Project
	res, err := c.Delete(ctx, state.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed during SKE project deletion", agg.Error())
		return
	}

	timeout, d := state.Timeouts.Delete(ctx, 10*time.Minute)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}
	process := res.WaitHandler(ctx, c, state.ID.ValueString()).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
		if !strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.Diagnostics.AddError("failed to verify project deletion", err.Error())
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import id to be set, got an empty string",
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), req.ID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
