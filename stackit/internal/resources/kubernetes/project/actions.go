package project

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KubernetesProject
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	c := r.client.Services.Kubernetes.Project
	res, err := c.CreateProjectWithResponse(ctx, plan.ProjectID.Value)
	if err != nil {
		diags.AddError("failed initiating SKE project creation", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed during SKE project creation", res.HasError.Error())
		return
	}

	plan.ID = types.StringNull()
	if res.JSON200.ProjectID != nil {
		plan.ID = types.StringValue(*res.JSON200.ProjectID)
	}

	// update state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	process := res.WaitHandler(ctx, c, plan.ProjectID.Value)
	if _, err := process.Wait(); err != nil {
		diags.AddError("failed verifying SKE project creation", err.Error())
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
	c := r.client.Services.Kubernetes.Project
	res, err := c.GetProjectWithResponse(ctx, state.ID.Value)
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		diags.AddError("failed requesting SKE project read", err.Error())
		return
	}
	if res.HasError != nil {
		if strings.Contains(res.HasError.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		diags.AddError("failed during SKE project read", res.HasError.Error())
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

	// handle creation
	c := r.client.Services.Kubernetes.Project
	res, err := c.DeleteProjectWithResponse(ctx, state.ID.Value)
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed initiating SKE project deletion", err.Error())
		return

	}
	if res.HasError != nil {
		if strings.Contains(res.HasError.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed during SKE project deletion", res.HasError.Error())
		return

	}
	process := res.WaitHandler(ctx, c, state.ID.Value)
	if _, err := process.Wait(); err != nil {
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
