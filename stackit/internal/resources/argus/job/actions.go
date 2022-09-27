package job

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

	var plan Job
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.Provider.Client()
	job := plan.ToClientJob()
	created := false
	if err := helper.RetryContext(ctx, common.DURATION_20M, func() *helper.RetryError {
		if !created {
			if _, err := c.Argus.Jobs.Create(ctx, plan.ProjectID.Value, plan.ArgusInstanceID.Value, job); err != nil {
				if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
					return helper.NonRetryableError(err)
				}
				return helper.RetryableError(err)
			}
			created = true
		}

		// Validate
		res, err := c.Argus.Jobs.Get(ctx, plan.ProjectID.Value, plan.ArgusInstanceID.Value, plan.Name.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		plan.FromClientJob(res.Data)
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to create job", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Job

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.Provider.Client()
	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		res, err := c.Argus.Jobs.Get(ctx, state.ProjectID.Value, state.ArgusInstanceID.Value, state.Name.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		state.FromClientJob(res.Data)
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read job", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Job
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.Provider.Client()
	job := plan.ToClientJob()
	updated := false
	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		if !updated {
			if _, err := c.Argus.Jobs.Update(ctx, plan.ProjectID.Value, plan.ArgusInstanceID.Value, job); err != nil {
				return helper.RetryableError(err)
			}
			updated = true
		}

		// read
		res, err := c.Argus.Jobs.Get(ctx, plan.ProjectID.Value, plan.ArgusInstanceID.Value, plan.Name.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		plan.FromClientJob(res.Data)
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to update job", err.Error())
		return
	}
	// update state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Job
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.Provider.Client()
	job := state.ToClientJob()
	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		if _, err := c.Argus.Jobs.Update(ctx, state.ProjectID.Value, state.ArgusInstanceID.Value, job); err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to delete job", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,id,name` where `id` is the instance id and `name` is the job name.\nInstead got: %q", req.ID),
		)
		return
	}

	projectID := idParts[0]
	instanceID := idParts[1]
	name := idParts[2]

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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("argus_instance_id"), instanceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)

}
