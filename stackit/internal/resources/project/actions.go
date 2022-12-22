package project

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v2/resource-management/projects"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
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
	}

	// update state
	diags = resp.State.Set(ctx, p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createProject(ctx context.Context, resp *resource.CreateResponse, plan Project) Project {
	labels := map[string]string{
		"billingReference": plan.BillingRef.ValueString(),
		"scope":            "PUBLIC",
	}

	members := []projects.ProjectMember{
		{
			Subject: r.client.GetConfig().ServiceAccountEmail,
			Role:    consts.ROLE_PROJECT_OWNER,
		},
		{
			Subject: plan.OwnerEmail.ValueString(),
			Role:    consts.ROLE_PROJECT_OWNER,
		},
	}

	c := r.client
	project, process, err := c.ResourceManagement.Projects.Create(ctx, plan.ParentContainerID.ValueString(), plan.Name.ValueString(), labels, members...)
	if err != nil {
		resp.Diagnostics.AddError("failed to create project", err.Error())
		return plan
	}

	if _, err := process.Wait(); err != nil {
		resp.Diagnostics.AddError("failed to verify project is active", err.Error())
		return plan
	}
	plan.ID = types.StringValue(project.ProjectID)
	plan.ContainerID = types.StringValue(project.ContainerID)
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

	if p.ContainerID.ValueString() == "" && p.ID.ValueString() != "" {
		res, err := c.ResourceManagement.Projects.Get(ctx, p.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("failed to fetch container ID", err.Error())
			return
		}
		p.ContainerID = types.StringValue(res.ContainerID)
	}

	project, err := c.ResourceManagement.Projects.Get(ctx, p.ContainerID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read project", err.Error())
		return
	}

	p.ID = types.StringValue(project.ProjectID)
	p.ContainerID = types.StringValue(project.ContainerID)
	p.ParentContainerID = types.StringValue(project.Parent.ContainerID)
	p.Name = types.StringValue(project.Name)
	p.BillingRef = types.StringValue(project.Labels["billingReference"])
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
	if plan.Name.Equal(state.Name) && plan.BillingRef.Equal(state.BillingRef) {
		return
	}
	c := r.client

	labels := map[string]string{
		"billingReference": plan.BillingRef.ValueString(),
		"scope":            "PUBLIC",
	}

	_, err := c.ResourceManagement.Projects.Update(ctx, plan.ParentContainerID.ValueString(), plan.ContainerID.ValueString(), plan.Name.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError("failed to update project", err.Error())
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
	process, err := c.ResourceManagement.Projects.Delete(ctx, state.ContainerID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete project", err.Error())
		return
	}

	if _, err := process.Wait(); err != nil {
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
