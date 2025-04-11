package project

import (
	"context"
	"fmt"
	"reflect"
	"time"

	rmv2 "github.com/SchwarzIT/community-stackit-go-client/pkg/services/resource-management/v2.0"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Project
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	plan = r.createProject(ctx, resp, plan)
	if resp.Diagnostics.HasError() {
		return
	}

	p := Project{
		ID:                types.StringValue(plan.ID.ValueString()),
		ContainerID:       types.StringValue(plan.ContainerID.ValueString()),
		ParentContainerID: types.StringValue(plan.ParentContainerID.ValueString()),
		Name:              types.StringValue(plan.Name.ValueString()),
		BillingRef:        types.StringValue(plan.BillingRef.ValueString()),
		OwnerEmail:        types.StringValue(plan.OwnerEmail.ValueString()),
		Timeouts:          plan.Timeouts,
		Labels:            plan.Labels,
	}
	// update state
	diags = resp.State.Set(ctx, p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createProject(ctx context.Context, resp *resource.CreateResponse, plan Project) Project {
	labels := rmv2.ResourceLabels{
		"billingReference": plan.BillingRef.ValueString(),
		"scope":            "PUBLIC",
	}

	for k, v := range plan.Labels {
		labels[k] = v
	}

	owner := rmv2.PROJECT_OWNER
	subj1 := r.client.Client.GetServiceAccountEmail()
	subj2 := plan.OwnerEmail.ValueString()
	members := []rmv2.ProjectMember{
		{
			Subject: &subj1,
			Role:    &owner,
		},
		{
			Subject: &subj2,
			Role:    &owner,
		},
	}

	body := rmv2.ProjectRequestBody{
		ContainerParentID: plan.ParentContainerID.ValueString(),
		Labels:            &labels,
		Members:           members,
		Name:              plan.Name.ValueString(),
	}
	res, err := r.client.ResourceManagement.Create(ctx, body)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON201"); agg != nil {
		resp.Diagnostics.AddError("failed creating project", agg.Error())
		return plan
	}

	// give API a bit of time to process request
	time.Sleep(30 * time.Second)

	timeout, d := plan.Timeouts.Create(ctx, 30*time.Minute)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return plan
	}
	process := res.WaitHandler(ctx, r.client.ResourceManagement, res.JSON201.ContainerID).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed validating project %s creation", res.JSON201.ProjectID), err.Error())
		return plan
	}

	plan.ID = types.StringValue(res.JSON201.ProjectID.String())
	plan.ContainerID = types.StringValue(res.JSON201.ContainerID)
	return plan
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var p Project

	diags := req.State.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.ResourceManagement.Get(ctx, p.ID.ValueString(), &rmv2.GetParams{})
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed reading project", agg.Error())
		return
	}

	p.ID = types.StringValue(res.JSON200.ProjectID.String())
	p.ContainerID = types.StringValue(res.JSON200.ContainerID)
	p.ParentContainerID = types.StringValue(res.JSON200.Parent.ContainerID)
	p.Name = types.StringValue(res.JSON200.Name)

	if res.JSON200.Labels != nil {
		l := *res.JSON200.Labels
		if v, ok := l["billingReference"]; ok {
			p.BillingRef = types.StringValue(v)
		}

		p.Labels = l

		// delete "hidden" labels which we assign via attribute
		// or similar to stay compatible with existing terraform code
		delete(p.Labels, "billingReference")
		delete(p.Labels, "scope")
	}

	diags = resp.State.Set(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Project
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state Project
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsUnknown() {
		plan.ID = state.ID
	}

	if plan.ContainerID.IsUnknown() {
		plan.ID = state.ContainerID
	}

	r.updateProject(ctx, plan, state, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) updateProject(ctx context.Context, plan, state Project, resp *resource.UpdateResponse) {
	if plan.Name.Equal(state.Name) && plan.BillingRef.Equal(state.BillingRef) && reflect.DeepEqual(plan.Labels, state.Labels) {
		return
	}

	labels := rmv2.ResourceLabels{
		"billingReference": plan.BillingRef.ValueString(),
		"scope":            "PUBLIC",
	}

	for k, v := range plan.Labels {
		// skip these internally / hidden reserved ones
		// this is to ensure backwards compatibility
		if k == "billingReference" || k == "scope" {
			continue
		}

		labels[k] = v
	}

	name := plan.Name.ValueString()
	parent := plan.ParentContainerID.ValueString()
	body := rmv2.UpdateJSONRequestBody{
		ContainerParentID: &parent,
		Labels:            &labels,
		Name:              &name,
	}
	res, err := r.client.ResourceManagement.Update(ctx, plan.ContainerID.ValueString(), body)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed updating project", agg.Error())
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Project
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	res, err := c.ResourceManagement.Delete(ctx, state.ContainerID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		resp.Diagnostics.AddError("failed deleting project", agg.Error())
		return
	}

	timeout, d := state.Timeouts.Delete(ctx, 30*time.Minute)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}
	process := res.WaitHandler(ctx, c.ResourceManagement, state.ContainerID.ValueString()).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed to verify project deletion", err.Error())
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// validate project id
	if err := clientValidate.ProjectID(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
