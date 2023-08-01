package job

import (
	"context"
	"fmt"
	"strings"

	scrapeconfig "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/scrape-config"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Job
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	job := scrapeconfig.CreateJSONRequestBody(plan.ToClientJob())
	res, err := c.Argus.ScrapeConfig.Create(ctx, plan.ProjectID.ValueString(), plan.ArgusInstanceID.ValueString(), job)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON202"); agg != nil {
		resp.Diagnostics.AddError("failed to create argus job", agg.Error())
		return
	}

	data := scrapeconfig.Job{}
	found := false
	for _, v := range res.JSON202.Data {
		if v.JobName == job.JobName {
			data = v
			found = true
			break
		}
	}
	if !found {
		resp.Diagnostics.AddError("failed to find job name", "no job by that name was found in create response")
		return
	}
	plan.FromClientJob(data)
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

	c := r.client

	res, err := c.Argus.ScrapeConfig.Get(ctx, state.ProjectID.ValueString(), state.ArgusInstanceID.ValueString(), state.Name.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		diags.AddError("failed to read argus job", agg.Error())
		return
	}

	state.FromClientJob(res.JSON200.Data)
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

	c := r.client
	job := scrapeconfig.UpdateJSONRequestBody(plan.ToClientUpdateJob())
	ures, err := c.Argus.ScrapeConfig.Update(ctx, plan.ProjectID.ValueString(), plan.ArgusInstanceID.ValueString(), plan.Name.ValueString(), job)
	if agg := common.Validate(&resp.Diagnostics, ures, err); agg != nil {
		resp.Diagnostics.AddError("failed to update argus job", agg.Error())
		return
	}

	// read job to verify update
	res, err := c.Argus.ScrapeConfig.Get(ctx, plan.ProjectID.ValueString(), plan.ArgusInstanceID.ValueString(), plan.Name.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to read argus job", agg.Error())
		return
	}
	plan.FromClientJob(res.JSON200.Data)

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

	c := r.client
	job := state.ToClientJob()
	params := &scrapeconfig.DeleteParams{
		JobName: []string{job.JobName},
	}
	res, err := c.Argus.ScrapeConfig.Delete(ctx, state.ProjectID.ValueString(), state.ArgusInstanceID.ValueString(), params)
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		resp.Diagnostics.AddError("failed to delete argus job", agg.Error())
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
