package credential

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/credentials"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Create - lifecycle function
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var cred Credential
	diags := req.Plan.Get(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	es := r.client.DataServices.ElasticSearch

	// handle creation
	res, err := es.Credentials.Create(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "service bind failed") {
			time.Sleep(30 * time.Second)
			res, err = es.Credentials.Create(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString())
		}
		if err != nil {
			resp.Diagnostics.AddError("failed credential creation", err.Error())
			return
		}
	}

	if err := r.applyClientResponse(ctx, &cred, credentials.GetResponse(res)); err != nil {
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
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var cred Credential
	diags := req.State.Get(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	es := r.client.DataServices.ElasticSearch

	// read instance credential
	res, err := es.Credentials.Get(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString(), cred.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read credential", err.Error())
		return
	}

	if err := r.applyClientResponse(ctx, &cred, res); err != nil {
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
func (r *Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

// Delete - lifecycle function
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var cred Credential
	resp.Diagnostics.Append(req.State.Get(ctx, &cred)...)
	if resp.Diagnostics.HasError() {
		return
	}

	es := r.client.DataServices.ElasticSearch

	res, err := es.Credentials.Delete(ctx, cred.ProjectID.ValueString(), cred.InstanceID.ValueString(), cred.ID.ValueString())
	if err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			resp.Diagnostics.AddError("failed to delete credential", err.Error())
			return
		}
	}
	if res.Error != "" {
		resp.Diagnostics.AddError("failed to delete credential", res.Error)
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
