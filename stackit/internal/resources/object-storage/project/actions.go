package project

import (
	"context"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ObjectStorageProject
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	c := r.client.ObjectStorage.Project
	res, err := c.Create(ctx, plan.ProjectID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if !strings.Contains(agg.Error(), common.ERR_UNEXPECTED_EOF) {
			resp.Diagnostics.AddError("failed ObjectStorage project creation", agg.Error())
			return
		}
	}

	plan.ID = types.StringValue(plan.ProjectID.ValueString())

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ObjectStorageProject
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read
	c := r.client.ObjectStorage.Project
	res, err := c.Get(ctx, state.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		resp.Diagnostics.AddError("failed ObjectStorage project read", agg.Error())
		return
	}
	if res.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ObjectStorageProject
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	c := r.client.ObjectStorage.Project
	res, err := c.Delete(ctx, state.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		resp.Diagnostics.AddError("failed ObjectStorage project deletion", agg.Error())
		return
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
