package bucket

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/buckets"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
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

	// handle creation
	b := r.createBucket(ctx, resp, bucket)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, Bucket{
		ID:           types.String{Value: b.Bucket.Name},
		Name:         types.String{Value: b.Bucket.Name},
		ProjectID:    types.String{Value: b.Project},
		Region:       types.String{Value: b.Bucket.Region},
		HostStyleURL: types.String{Value: b.Bucket.URLVirtualHostedStyle},
		PathStyleURL: types.String{Value: b.Bucket.URLPathStyle},
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createBucket(ctx context.Context, resp *resource.CreateResponse, plan Bucket) buckets.BucketResponse {
	c := r.client
	var b buckets.BucketResponse

	// Create bucket
	process, err := c.ObjectStorage.Buckets.Create(ctx, plan.ProjectID.Value, plan.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to verify bucket creation", err.Error())
		return b
	}
	process.SetTimeout(10 * time.Minute)
	tmp, err := process.Wait()
	if err != nil {
		resp.Diagnostics.AddError("failed to verify bucket creation", err.Error())
		return b
	}
	b = tmp.(buckets.BucketResponse)
	if err != nil {
		resp.Diagnostics.AddError("failed to verify bucket creation", err.Error())
		return b
	}

	return b
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

	b, err := c.ObjectStorage.Buckets.Get(ctx, state.ProjectID.Value, state.Name.Value)
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read bucket", err.Error())
		return
	}

	state.Region = types.String{Value: b.Bucket.Region}
	state.HostStyleURL = types.String{Value: b.Bucket.URLVirtualHostedStyle}
	state.PathStyleURL = types.String{Value: b.Bucket.URLPathStyle}

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

	c := r.client

	httpClient := c.GetHTTPClient()
	t := httpClient.Timeout

	httpClient.Timeout = time.Minute
	process, err := c.ObjectStorage.Buckets.Delete(ctx, state.ProjectID.Value, state.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to verify bucket deletion", err.Error())
		return
	}
	process.SetTimeout(10 * time.Minute)
	_, err = process.Wait()
	if err != nil {
		resp.Diagnostics.AddError("failed to verify bucket deletion", err.Error())
		return
	}
	httpClient.Timeout = t

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,name` where `name` is the cluster name.\nInstead got: %q", req.ID),
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

	// validate bucket name
	if err := buckets.ValidateBucketName(idParts[1]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate bucket name.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

}
