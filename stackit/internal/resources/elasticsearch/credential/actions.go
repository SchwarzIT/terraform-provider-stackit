package credential

import (
	"context"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/credentials"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"net/http"
	"strings"
)

// TODO: Questions:
// - store settings for testing
// - how to test. created my own project
// - validation
// - schema
// - Update function
// - ressource, data-source

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var cred Credential
	diags := req.Plan.Get(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// validate
	if err := r.validate(ctx, &cred); err != nil {
		resp.Diagnostics.AddError("failed credential validation", err.Error())
		return
	}

	es := r.client.DataServices.ElasticSearch

	// handle creation
	res, err := es.Credentials.Create(ctx, cred.ProjectID.Value, cred.InstanceID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed credential creation", err.Error())
		return
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

// List - lifecycle function
func (r Resource) List(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var cred Credential
	diags := req.State.Get(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	es := r.client.DataServices.ElasticSearch

	// read instance
	res, err := es.Credentials.List(ctx, cred.ProjectID.Value, cred.InstanceID.Value)
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to list credentials", err.Error())
		return
	}

	_ = res
	// for _,_  := range res.CredentialsList {
	//
	//
	// }

}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var cred Credential
	diags := req.State.Get(ctx, &cred)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	es := r.client.DataServices.ElasticSearch

	// read instance
	res, err := es.Credentials.Get(ctx, cred.ProjectID.Value, cred.InstanceID.Value, cred.ID.Value)
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
func (r Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var cred Credential
	resp.Diagnostics.Append(req.State.Get(ctx, &cred)...)
	if resp.Diagnostics.HasError() {
		return
	}

	es := r.client.DataServices.ElasticSearch

	res, err := es.Credentials.Delete(ctx, cred.ProjectID.Value, cred.InstanceID.Value, cred.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete credential", err.Error())
		return
	}
	if res.Error != "" {
		resp.Diagnostics.AddError("failed to delete credential", res.Error)
		return
	}

	resp.State.RemoveResource(ctx)
}
