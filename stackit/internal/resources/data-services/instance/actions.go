package instance

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0/instances"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
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
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("failed to process plan", "failed during plan processing")
		return
	}

	// validate
	if err := r.validate(ctx, &resp.Diagnostics, &plan); err != nil {
		resp.Diagnostics.AddError("failed instance validation", err.Error())
		return
	}

	acl := []string{}
	for _, v := range plan.ACL.Elements() {
		nv, err := common.ToString(ctx, v)
		if err != nil {
			continue
		}
		acl = append(acl, nv)
	}

	// handle creation
	aclString := strings.Join(acl, ",")
	params := instances.InstanceParameters{
		SgwAcl: &aclString,
	}
	body := instances.InstanceProvisionRequest{
		InstanceName: plan.Name.ValueString(),
		PlanID:       plan.PlanID.ValueString(),
		Parameters:   &params,
	}
	res, err := r.client.Instances.Provision(ctx, plan.ProjectID.ValueString(), body)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON202"); agg != nil {
		resp.Diagnostics.AddError("failed instance provisioning", agg.Error())
		return
	}

	// set state
	plan.ID = types.StringValue(res.JSON202.InstanceID)
	if res.JSON202.InstanceID == "" {
		resp.Diagnostics.AddError("received an empty instance ID", fmt.Sprintf("invalid instance id: %+v", *res.JSON202))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), res.JSON202.InstanceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), plan.ProjectID.ValueString())...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Create(ctx, 40*time.Minute)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	process := res.WaitHandler(ctx, r.client.Instances, plan.ProjectID.ValueString(), plan.ID.ValueString()).SetTimeout(timeout)
	instance, err := process.WaitWithContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed instance `create` validation", err.Error())
		return
	}

	i, ok := instance.(*instances.GetResponse)
	if !ok {
		resp.Diagnostics.AddError("failed to parse client response", "response is not of *instances.GetResponse")
		return
	}

	if i.JSON200 == nil {
		resp.Diagnostics.AddError("failed to parse client response", "JSON200 == nil")
		return
	}

	if err := r.applyClientResponse(ctx, &plan, i.JSON200); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
	// read instance
	res, err := r.client.Instances.Get(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound, http.StatusGone) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to get instance", agg.Error())
		return
	}

	if err := r.applyClientResponse(ctx, &state, res.JSON200); err != nil {
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
	if err := r.validate(ctx, &resp.Diagnostics, &plan); err != nil {
		resp.Diagnostics.AddError("failed validation", err.Error())
		return
	}

	acl := []string{}
	for _, v := range plan.ACL.Elements() {
		nv, err := common.ToString(ctx, v)
		if err != nil {
			continue
		}
		acl = append(acl, nv)
	}

	// handle update
	aclString := strings.Join(acl, ",")
	params := instances.InstanceParameters{
		SgwAcl: &aclString,
	}
	body := instances.UpdateJSONRequestBody{
		PlanID:     plan.PlanID.ValueString(),
		Parameters: &params,
	}
	res, err := r.client.Instances.Update(ctx, state.ProjectID.ValueString(), state.ID.ValueString(), body)
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		resp.Diagnostics.AddError("failed instance update", agg.Error())
		return
	}

	timeout, d := plan.Timeouts.Update(ctx, 40*time.Minute)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	process := res.WaitHandler(ctx, r.client.Instances, state.ProjectID.ValueString(), state.ID.ValueString()).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed instance update validation", err.Error())
		return
	}

	// mitigate an API bug that returns old data after an update completed
	time.Sleep(1 * time.Minute)

	newRes, err := r.client.Instances.Get(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, newRes, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to read after update", agg.Error())
		return
	}
	if err := r.applyClientResponse(ctx, &plan, newRes.JSON200); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
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

	// handle deletion
	res, err := r.client.Instances.Deprovision(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound, http.StatusForbidden) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to deprovision instance", agg.Error())
		return
	}

	timeout, d := state.Timeouts.Delete(ctx, 40*time.Minute)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	process := res.WaitHandler(ctx, r.client.Instances, state.ProjectID.ValueString(), state.ID.ValueString()).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed to verify instance deprovision", err.Error())
	}

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

	plan, version, err := r.getPlanAndVersion(ctx, &resp.Diagnostics, idParts[0], idParts[1])
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
