package bucket

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/object-storage/v1.0.1/bucket"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var bucket Bucket
	diags := req.Plan.Get(ctx, &bucket)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := bucket.Timeouts.Create(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	res := r.createBucket(ctx, resp, bucket, timeout)
	if resp.Diagnostics.HasError() {
		return
	}

	b := res.JSON200
	// update state
	diags = resp.State.Set(ctx, Bucket{
		ID:                     types.StringValue(b.Bucket.Name),
		Name:                   types.StringValue(b.Bucket.Name),
		ObjectStorageProjectID: types.StringValue(b.Project),
		Region:                 types.StringValue(b.Bucket.Region),
		HostStyleURL:           types.StringValue(b.Bucket.UrlVirtualHostedStyle),
		PathStyleURL:           types.StringValue(b.Bucket.UrlPathStyle),
		Timeouts:               bucket.Timeouts,
	})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createBucket(ctx context.Context, resp *resource.CreateResponse, plan Bucket, timeout time.Duration) *bucket.GetResponse {
	c := r.client
	b := &bucket.GetResponse{}

	// Create bucket
	res, err := c.ObjectStorage.Bucket.Create(ctx, plan.ObjectStorageProjectID.ValueString(), plan.Name.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON201"); agg != nil {
		resp.Diagnostics.AddError("failed to create bucket", agg.Error())
		return b
	}
	process := res.WaitHandler(ctx, c.ObjectStorage.Bucket, plan.ObjectStorageProjectID.ValueString(), plan.Name.ValueString())
	process.SetTimeout(timeout)
	tmp, err := process.WaitWithContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to verify bucket creation", err.Error())
		return b
	}
	nb, ok := tmp.(*bucket.GetResponse)
	if !ok {
		resp.Diagnostics.AddError("failed to parse wait response", "not *bucket.GetResponse")
		return b
	}
	return nb
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var state Bucket

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.ObjectStorage.Bucket.Get(ctx, state.ObjectStorageProjectID.ValueString(), state.Name.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200.Bucket"); agg != nil {
		resp.Diagnostics.AddError("failed to read bucket", agg.Error())
		return
	}

	if res.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Region = types.StringValue(res.JSON200.Bucket.Region)
	state.HostStyleURL = types.StringValue(res.JSON200.Bucket.UrlVirtualHostedStyle)
	state.PathStyleURL = types.StringValue(res.JSON200.Bucket.UrlPathStyle)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function - not used for this resource
func (r Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Bucket
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle timeout
	timeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	res, err := c.ObjectStorage.Bucket.Delete(ctx, state.ObjectStorageProjectID.ValueString(), state.Name.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		if !validate.StatusEquals(res, http.StatusNotFound) {
			resp.Diagnostics.AddError("failed to delete bucket", agg.Error())
			return
		}
	}

	process := res.WaitHandler(ctx, c.ObjectStorage.Bucket, state.ObjectStorageProjectID.ValueString(), state.Name.ValueString())
	process.SetTimeout(timeout)
	if _, err = process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed to verify bucket deletion", err.Error())
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
			fmt.Sprintf("Expected import identifier with format: `object_storage_project_id,name` where `name` is the cluster name.\nInstead got: %q", req.ID),
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_storage_project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

}
