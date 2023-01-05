package credential

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

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

	// handle creation
	res, err := r.client.Credentials.PostWithResponse(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed preparing credential creation request", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("empty response", "received an empty response during credential creation. res == nil")
		return
	}
	if res.HasError != nil {
		if strings.Contains(res.HasError.Error(), "service bind failed") {
			time.Sleep(30 * time.Second)
			res, err = r.client.Credentials.PostWithResponse(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("failed preparing 2nd credential creation request", err.Error())
				return
			}
		}
		if res.HasError != nil {
			resp.Diagnostics.AddError("failed credential creation", err.Error())
			return
		}
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("failed parsing credential creation response", "JSON200 == nil")
	}

	if err := r.applyClientResponse(ctx, &cred, res.JSON200); err != nil {
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
	res, err := r.client.Credentials.GetWithResponse(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString(), cred.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed preparing get credential request", err.Error())
	}
	if res.HasError != nil {
		if res.StatusCode() == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read credential", res.HasError.Error())
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("failed parsing get credential response", "JSON200 == nil")
	}

	if err := r.applyClientResponse(ctx, &cred, res.JSON200); err != nil {
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

// Update - lifecycle function - not used for this resource
func (r Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var cred Credential
	resp.Diagnostics.Append(req.State.Get(ctx, &cred)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Credentials.DeleteWithResponse(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString(), cred.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed preparing delete credential request", err.Error())
	}
	if res.HasError != nil {
		if !strings.Contains(err.Error(), "EOF") {
			resp.Diagnostics.AddError("failed to delete credential", res.HasError.Error())
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
			fmt.Sprintf("Expected import identifier with format: `project_id,instance_id,credential_id`.\nInstead got: %q", req.ID),
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

	if resp.Diagnostics.HasError() {
		return
	}
}
