package credential

import (
	"context"
	"net/http"
	"time"

	accesskey "github.com/SchwarzIT/community-stackit-go-client/pkg/services/object-storage/v1.0.1/access-key"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data Credential
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// pre process plan
	r.preProcessConfig(&resp.Diagnostics, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle project init
	r.enableProject(ctx, &resp.Diagnostics, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	res := r.createAccessKey(ctx, resp, data)
	if resp.Diagnostics.HasError() {
		return
	}

	k := res.JSON201
	// update state
	diags = resp.State.Set(ctx, Credential{
		ID:                     types.StringValue(k.KeyID),
		ObjectStorageProjectID: types.StringValue(k.Project),
		CredentialsGroupID:     types.StringValue(data.CredentialsGroupID.ValueString()),
		Expiry:                 types.StringValue(k.Expires),
		DisplayName:            types.StringValue(k.DisplayName),
		AccessKey:              types.StringValue(k.AccessKey),
		SecretAccessKey:        types.StringValue(k.SecretAccessKey),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) preProcessConfig(diags *diag.Diagnostics, b *Credential) {
	projectID := b.ProjectID.ValueString()
	osProjectID := b.ObjectStorageProjectID.ValueString()
	if projectID == "" && osProjectID == "" {
		diags.AddError("project_id or object_storage_project_id must be set", "please note that object_storage_project_id is deprecated and will be removed in a future release")
		return
	}
	if projectID == "" {
		b.ProjectID = b.ObjectStorageProjectID
	}
	if osProjectID == "" {
		b.ObjectStorageProjectID = b.ProjectID
	}
}

func (r Resource) enableProject(ctx context.Context, diags *diag.Diagnostics, b *Credential) {
	projectID := b.ProjectID.ValueString()
	c := r.client.ObjectStorage.Project

	status, err := c.Get(ctx, projectID)
	if agg := common.Validate(&diag.Diagnostics{}, status, err, "JSON200"); agg != nil {
		if status == nil || status.StatusCode() != http.StatusNotFound {
			diags.AddError("failed to fetch Object Storage project status", agg.Error())
			return
		}
	}

	res, err := c.Create(ctx, projectID)
	if agg := common.Validate(diags, res, err); agg != nil {
		diags.AddError("failed during Object Storage project init", agg.Error())
		return
	}
}

func (r Resource) createAccessKey(ctx context.Context, resp *resource.CreateResponse, key Credential) *accesskey.CreateResponse {
	c := r.client
	body := accesskey.CreateJSONRequestBody{}
	if !key.Expiry.IsNull() && !key.Expiry.IsUnknown() {
		t, err := time.Parse("2006-01-02T15:04:05.999Z", key.Expiry.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("couldn't parse expiry", err.Error())
			return nil
		}
		body.Expires = &t
	}
	cg := key.CredentialsGroupID.ValueString()
	params := &accesskey.CreateParams{
		CredentialsGroup: &cg,
	}
	if cg == "" {
		params.CredentialsGroup = nil
	}
	res, err := c.ObjectStorage.AccessKey.Create(ctx, key.ObjectStorageProjectID.ValueString(), params, body)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON201"); agg != nil {
		resp.Diagnostics.AddError("failed to create credential", agg.Error())
		return res
	}

	return res
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var state Credential

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// pre process plan
	r.preProcessConfig(&resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	cg := state.CredentialsGroupID.ValueString()
	params := &accesskey.GetParams{
		CredentialsGroup: &cg,
	}
	if cg == "" {
		params.CredentialsGroup = nil
	}
	res, err := c.ObjectStorage.AccessKey.Get(ctx, state.ObjectStorageProjectID.ValueString(), params)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200.AccessKeys"); agg != nil {
		resp.Diagnostics.AddError("failed to list credentials", agg.Error())
		return
	}

	found := false
	for _, k := range res.JSON200.AccessKeys {
		if k.KeyID != state.ID.ValueString() {
			continue
		}
		found = true
		state.DisplayName = types.StringValue(k.DisplayName)
		state.Expiry = types.StringValue(k.Expires)
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
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Credential
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// pre-process config
	r.preProcessConfig(&resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Credential
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// pre process plan
	r.preProcessConfig(&resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	cg := state.CredentialsGroupID.ValueString()
	params := &accesskey.DeleteParams{
		CredentialsGroup: &cg,
	}
	if cg == "" {
		params.CredentialsGroup = nil
	}

	c := r.client
	res, err := c.ObjectStorage.AccessKey.Delete(ctx, state.ObjectStorageProjectID.ValueString(), state.ID.ValueString(), params)
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		resp.Diagnostics.AddError("failed to delete credential", agg.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
