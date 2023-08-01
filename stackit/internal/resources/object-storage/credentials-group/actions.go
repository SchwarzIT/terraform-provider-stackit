package credentialsgroup

import (
	"context"
	"fmt"
	"strings"

	credentialsgroup "github.com/SchwarzIT/community-stackit-go-client/pkg/services/object-storage/v1.0.1/credentials-group"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
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

	err := r.createCredentialGroup(ctx, &resp.Diagnostics, &data)
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

func (r Resource) createCredentialGroup(ctx context.Context, diags *diag.Diagnostics, data *CredentialsGroup) diag.Diagnostic {
	c := r.client
	body := credentialsgroup.CreateJSONRequestBody{
		DisplayName: data.Name.ValueString(),
	}
	cres, err := c.ObjectStorage.CredentialsGroup.Create(ctx, data.ObjectStorageProjectID.ValueString(), body)
	if agg := common.Validate(diags, cres, err); agg != nil {
		return diag.NewErrorDiagnostic("failed to create credential group", agg.Error())

	}

	res, err := c.ObjectStorage.CredentialsGroup.Get(ctx, data.ObjectStorageProjectID.ValueString())
	if agg := common.Validate(diags, res, err, "JSON200.CredentialsGroups"); agg != nil {
		return diag.NewErrorDiagnostic("failed to list credential groups", agg.Error())
	}

	for _, group := range res.JSON200.CredentialsGroups {
		if group.DisplayName == data.Name.ValueString() {
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

	res, err := c.ObjectStorage.CredentialsGroup.Get(ctx, state.ObjectStorageProjectID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200.CredentialsGroups"); agg != nil {
		resp.Diagnostics.AddError("failed to read credential groups", agg.Error())
		return
	}

	found := false
	for _, group := range res.JSON200.CredentialsGroups {
		if group.DisplayName == state.Name.ValueString() {
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

	c := r.client.ObjectStorage.CredentialsGroup
	res, err := c.Delete(ctx, state.ObjectStorageProjectID.ValueString(), state.ID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		resp.Diagnostics.AddError("failed to delete credential groups", agg.Error())
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
			fmt.Sprintf("Expected import identifier with format: `object_storage_project_id,id` where `name` is the credentials group name.\nInstead got: %q", req.ID),
		)
		return
	}

	// validate project id
	if err := clientValidate.ProjectID(idParts[0]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate object_storage_project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_storage_project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
