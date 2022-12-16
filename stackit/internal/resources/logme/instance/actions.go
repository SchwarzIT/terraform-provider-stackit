package instance

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/instances"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load plan
	var plan Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// validate
	if err := r.validate(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("failed instance validation", err.Error())
		return
	}

	acl := []string{}
	for _, v := range plan.ACL.Elems {
		nv, err := common.ToString(ctx, v)
		if err != nil {
			continue
		}
		acl = append(acl, nv)
	}

	lm := r.client.DataServices.LogMe

	// handle creation
	res, wait, err := lm.Instances.Create(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), plan.PlanID.ValueString(), map[string]string{
		"sgw_acl": strings.Join(acl, ","),
	})

	if err != nil {
		resp.Diagnostics.AddError("failed instance creation", err.Error())
		return
	}

	// set state
	plan.ID = types.StringValue(res.InstanceID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), res.InstanceID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := wait.Wait()
	if err != nil {
		resp.Diagnostics.AddError("failed instance `create` validation", err.Error())
		return
	}

	i, ok := instance.(instances.Instance)
	if !ok {
		resp.Diagnostics.AddError("failed to parse client response", "response is not of instances.Instance")
		return
	}

	if err := r.applyClientResponse(ctx, &plan, i); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	// update state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	lm := r.client.DataServices.LogMe

	// read instance
	i, err := lm.Instances.Get(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read instance", err.Error())
		return
	}

	if err := r.applyClientResponse(ctx, &state, i); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	// update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state Instance
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	if plan.ACL.IsUnknown() || plan.ACL.IsNull() {
		plan.ACL = state.ACL
	}

	// validate
	if err := r.validate(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("failed validation", err.Error())
		return
	}

	acl := []string{}
	for _, v := range plan.ACL.Elems {
		nv, err := common.ToString(ctx, v)
		if err != nil {
			continue
		}
		acl = append(acl, nv)
	}
	lm := r.client.DataServices.LogMe

	// handle update
	_, process, err := lm.Instances.Update(ctx, state.ProjectID.ValueString(), state.ID.ValueString(), plan.PlanID.ValueString(), map[string]string{
		"sgw_acl": strings.Join(acl, ","),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed instance update", err.Error())
		return
	}

	instance, err := process.Wait()
	if err != nil {
		elaborate := ""
		if i, ok := instance.(instances.Instance); ok {
			elaborate = "\n- type: " + i.LastOperation.Type + "\n- state: " + i.LastOperation.State
		} else {
			elaborate = "\n- couldn't parst response as instances.Instance"
		}
		resp.Diagnostics.AddError("failed instance update validation"+elaborate, err.Error())
		return
	}

	i, ok := instance.(instances.Instance)
	if !ok {
		resp.Diagnostics.AddError("failed to parse client response", "response is not of instances.Instance")
		return
	}

	planID := plan.PlanID
	if err := r.applyClientResponse(ctx, &plan, i); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	if !plan.PlanID.Equal(planID) {
		resp.Diagnostics.AddError("server returned wrong plan ID after update", fmt.Sprintf("expected plan ID %s but received %s", planID.ValueString(), plan.PlanID.ValueString()))
		return
	}

	// update state
	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Instance
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lm := r.client.DataServices.LogMe

	// handle deletion
	process, err := lm.Instances.Delete(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete instance", err.Error())
		return
	}

	if instance, err := process.Wait(); err != nil {
		resp.Diagnostics.AddError("failed to verify instance deletion", err.Error())
		if i, ok := instance.(instances.Instance); ok {
			resp.Diagnostics.AddError("instance delete response", "- type: "+i.LastOperation.Type+"\n- state: "+i.LastOperation.State)
		} else {
			resp.Diagnostics.AddError("instance delete response", "- couldn't parst response as instances.Instance")
		}
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,instance_id`.\nInstead got: %q", req.ID),
		)
		return
	}

	// validate project id
	if err := clientValidate.ProjectID(idParts[0]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate project_id.\n%s", err.Error()),
		)
		return
	}

	plan, version, err := r.getPlanAndVersion(ctx, idParts[0], idParts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Error during import",
			err.Error(),
		)
		return
	}
	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan"), plan)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("version"), version)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
