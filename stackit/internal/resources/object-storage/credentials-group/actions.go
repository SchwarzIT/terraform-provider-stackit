package credentialsgroup

import (
	"context"
	"fmt"
	"strings"

	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CredentialsGroup
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.createCredentialGroup(ctx, &data)
	if err != nil {
		resp.Diagnostics.Append(err)
		return
	}

	// update state
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createCredentialGroup(ctx context.Context, data *CredentialsGroup) diag.Diagnostic {
	c := r.client
	err := c.ObjectStorage.CredentialsGroup.Create(ctx, data.ProjectID.Value, data.Name.Value)
	if err != nil {
		return diag.NewErrorDiagnostic("failed to create credential group", err.Error())

	}

	res, err := c.ObjectStorage.CredentialsGroup.List(ctx, data.ProjectID.Value)
	if err != nil {
		return diag.NewErrorDiagnostic("failed to list credential groups", err.Error())
	}

	for _, group := range res.CredentialsGroups {
		if group.DisplayName == data.Name.Value {
			data.ID = types.StringValue(group.CredentialsGroupID)
			data.URN = types.StringValue(group.URN)
			break
		}
	}
	return nil
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var state CredentialsGroup

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.ObjectStorage.CredentialsGroup.List(ctx, state.ProjectID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to read credential group", err.Error())
		return
	}

	found := false
	for _, group := range res.CredentialsGroups {
		if group.DisplayName == state.Name.Value {
			found = true
			state.ID = types.StringValue(group.CredentialsGroupID)
			state.URN = types.StringValue(group.URN)
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function - not used for this resource
func (r Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CredentialsGroup
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	if err := c.ObjectStorage.CredentialsGroup.Delete(ctx, state.ProjectID.Value, state.ID.Value); err != nil {
		resp.Diagnostics.AddError("failed to delete credential group", err.Error())
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
			fmt.Sprintf("Expected import identifier with format: `project_id,id` where `name` is the credentials group name.\nInstead got: %q", req.ID),
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
