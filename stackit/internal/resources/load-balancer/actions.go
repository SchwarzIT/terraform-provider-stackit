package loadbalancer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/load-balancer/1beta.0.0/instances"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/load-balancer/1beta.0.0/project"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
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
		return
	}

	r.enableProject(ctx, plan.ProjectID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.LoadBalancer.Instances.Create(ctx, plan.ProjectID.ValueString(), &instances.CreateParams{}, prepareData(plan))
	if agg := validate.Response(res, err, "JSON200.Name"); agg != nil {
		diags.AddError("Couldn't create instance", agg.Error())
		return
	}
	if res.StatusCode() != http.StatusOK {
		diags.AddError("Couldn't create instance", fmt.Sprintf("Received status code %d", res.StatusCode()))
		common.Dump(&resp.Diagnostics, res.Body)
		return
	}
	for _, e := range *res.JSON200.Errors {
		detail := ""
		if e.Type != nil {
			detail = fmt.Sprintf("Type: %s", *e.Type)
		}
		if e.Description != nil {
			detail = fmt.Sprintf("%s\nDescription: %s", detail, *e.Description)
		}
		diags.AddError("Couldn't create instance", detail)
	}
	if diags.HasError() {
		return
	}

	plan.ID = types.StringValue(*res.JSON200.Name)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.ID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	process := res.WaitHandler(ctx, r.client.LoadBalancer.Instances, plan.ProjectID.ValueString(), plan.Name.ValueString())
	wres, err := process.Wait()
	if err != nil {
		diags.AddError("Received an error while waiting for load balancer instance to be created", err.Error())
		return
	}

	if g, ok := wres.(*instances.GetResponse); ok {
		if agg := validate.Response(g, nil, "JSON200.Name"); agg != nil {
			diags.AddError("Couldn't get instance information", agg.Error())
			return
		}
		plan.ExternalAddress = resToStr(g.JSON200.ExternalAddress)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r Resource) isProjectReady(ctx context.Context, projectID string, diags *diag.Diagnostics) bool {
	res, err := r.client.LoadBalancer.Project.GetStatus(ctx, projectID)
	if agg := validate.Response(res, err, "JSON200.Status"); agg != nil {
		diags.AddError("Couldn't get project status", agg.Error())
		return false
	}
	if *res.JSON200.Status == project.STATUS_READY {
		return true
	}
	return false
}

func (r Resource) enableProject(ctx context.Context, projectID string, diags *diag.Diagnostics) {
	isReady := r.isProjectReady(ctx, projectID, diags)
	if diags.HasError() || isReady {
		return
	}
	res, err := r.client.LoadBalancer.Project.EnableProject(ctx, projectID, &project.EnableProjectParams{})
	if agg := validate.Response(res, err); agg != nil {
		diags.AddError("Couldn't enable project", agg.Error())
		return
	}
	process := res.WaitHandler(ctx, r.client.LoadBalancer.Project, projectID)
	if _, err := process.Wait(); err != nil {
		diags.AddError("Received an error while waiting for project to be enabled", err.Error())
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

	res, err := r.client.LoadBalancer.Instances.Get(ctx, state.ProjectID.ValueString(), state.Name.ValueString())
	if agg := validate.Response(res, err, "JSON200.Name"); agg != nil {
		resp.Diagnostics.AddError("Couldn't get instance information", agg.Error())
		return
	}
	if res.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Couldn't get instance information", fmt.Sprintf("Received status code %d", res.StatusCode()))
		common.Dump(&resp.Diagnostics, res.Body)
		return
	}
	for _, e := range *res.JSON200.Errors {
		detail := ""
		if e.Type != nil {
			detail = fmt.Sprintf("Type: %s", *e.Type)
		}
		if e.Description != nil {
			detail = fmt.Sprintf("%s\nDescription: %s", detail, *e.Description)
		}
		resp.Diagnostics.AddError("Couldn't get instance information", detail)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	state.ExternalAddress = resToStr(res.JSON200.ExternalAddress)

	// update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Instance
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.LoadBalancer.Instances.Delete(ctx, state.ProjectID.ValueString(), state.Name.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		resp.Diagnostics.AddError("Couldn't delete instance", agg.Error())
		return
	}
	process := res.WaitHandler(ctx, r.client.LoadBalancer.Instances, state.ProjectID.ValueString(), state.Name.ValueString())
	if _, err := process.Wait(); err != nil {
		resp.Diagnostics.AddError("Received an error while waiting for load balancer instance to be deleted", err.Error())
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
			fmt.Sprintf("Expected import identifier with format: `project_id,id` where `id` is the instance ID.\nInstead got: %q", req.ID),
		)
		return
	}

	// validate project id
	if err := clientValidate.ProjectID(idParts[0]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate kubernetes_project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}

}
