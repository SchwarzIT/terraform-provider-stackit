package bucket

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/buckets"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	c := r.Provider.Client()
	created := false
	var b buckets.BucketResponse

	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		var err error

		// Create bucket
		if !created {
			if err = c.ObjectStorage.Buckets.Create(ctx, plan.ProjectID.Value, plan.Name.Value); err != nil {
				if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
					return helper.NonRetryableError(err)
				}
				return helper.RetryableError(err)
			}
			created = true
		}

		// Get bucket
		b, err = c.ObjectStorage.Buckets.Get(ctx, plan.ProjectID.Value, plan.Name.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to verify bucket creation", err.Error())
		return b
	}

	return b
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.Provider.Client()
	var state Bucket
	var b buckets.BucketResponse

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	missing := false
	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		var err error
		b, err = c.ObjectStorage.Buckets.Get(ctx, state.ProjectID.Value, state.Name.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				resp.State.RemoveResource(ctx)
				missing = true
				return nil
			}
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read bucket", err.Error())
		return
	}

	if missing {
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

	c := r.Provider.Client()

	httpClient := c.GetHTTPClient()
	t := httpClient.Timeout

	deleted := false
	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		if !deleted {
			httpClient.Timeout = time.Minute
			if err := c.ObjectStorage.Buckets.Delete(ctx, state.ProjectID.Value, state.Name.Value); err != nil {
				if !strings.Contains(err.Error(), common.ERR_CLIENT_TIMEOUT) {
					return helper.RetryableError(err)
				}
			}
			httpClient.Timeout = t
			deleted = true
		}

		// Verify bucket deletion
		res, err := c.ObjectStorage.Buckets.List(ctx, state.ProjectID.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return nil
			}
			return helper.RetryableError(err)
		}
		for _, b := range res.Buckets {
			if b.Name == state.Name.Value {
				return helper.RetryableError(fmt.Errorf("bucket %s is still available", b.Name))
			}
		}
		return nil
	}); err != nil {
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
