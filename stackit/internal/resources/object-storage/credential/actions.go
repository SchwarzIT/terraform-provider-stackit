package credential

import (
	"context"

	keys "github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/access-keys"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.Provider.IsConfigured() {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on another resource.",
		)
		return
	}

	var data Credential
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	k := r.createAccessKey(ctx, resp, data)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, Credential{
		ID:                 types.String{Value: k.KeyID},
		ProjectID:          types.String{Value: k.Project},
		CredentialsGroupID: types.String{Value: data.CredentialsGroupID.Value},
		Expiry:             types.String{Value: k.Expires},
		DisplayName:        types.String{Value: k.DisplayName},
		AccessKey:          types.String{Value: k.AccessKey},
		SecretAccessKey:    types.String{Value: k.SecretAccessKey},
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createAccessKey(ctx context.Context, resp *resource.CreateResponse, key Credential) keys.AccessKeyCreateResponse {
	c := r.Provider.Client()
	res, err := c.ObjectStorage.AccessKeys.Create(ctx, key.ProjectID.Value, key.Expiry.Value, key.CredentialsGroupID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to create credential", err.Error())
		return res
	}
	return res
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.Provider.Client()
	var state Credential

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := c.ObjectStorage.AccessKeys.List(ctx, state.ProjectID.Value, state.CredentialsGroupID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to read credential list", err.Error())
		return
	}

	found := false
	for _, k := range list.AccessKeys {
		if k.KeyID != state.ID.Value {
			continue
		}
		found = true
		state.DisplayName = types.String{Value: k.DisplayName}
		state.Expiry = types.String{Value: k.Expires}
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		break
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
}

// Update - lifecycle function - not used for this resource
func (r Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Credential
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.Provider.Client()
	err := c.ObjectStorage.AccessKeys.Delete(ctx, state.ProjectID.Value, state.ID.Value, state.CredentialsGroupID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete credential", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
