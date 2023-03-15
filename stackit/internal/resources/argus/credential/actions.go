package credential

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var cred Credential
	diags := req.Plan.Get(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Instances.InstanceCredentialsCreateWithResponse(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString())
	if agg := validate.Response(res, err, "JSON201"); agg != nil {
		resp.Diagnostics.AddError("failed to create argus instance credentials", agg.Error())
		return
	}

	if err := r.applyClientResponse(ctx, &cred, res.JSON201.Credentials); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	// update state
	diags = resp.State.Set(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var cred Credential
	diags := req.State.Get(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read instance credential
	res, err := r.client.Instances.InstanceCredentialsReadWithResponse(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString(), cred.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read credential", agg.Error())
		return
	}

	// update state
	diags = resp.State.Set(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function - not used for this resource
func (r Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var cred Credential
	resp.Diagnostics.Append(req.State.Get(ctx, &cred)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Instances.InstanceCredentialsDeleteWithResponse(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString(), cred.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if !strings.Contains(agg.Error(), "EOF") {
			resp.Diagnostics.AddError("failed to delete credential", agg.Error())
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,instance_id,username`.\nInstead got: %q", req.ID),
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

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), idParts[2])...)

	if resp.Diagnostics.HasError() {
		return
	}
}
