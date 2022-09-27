package credentialsgroup

import (
	"context"
	"fmt"
	"strings"
	"time"

	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	default_retry_duration = 10 * time.Minute
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

	var data CredentialsGroup
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created := false
	if err := helper.RetryContext(ctx, default_retry_duration, r.createCredentialGroup(ctx, &data, &created)); err != nil {
		resp.Diagnostics.AddError("failed to create credential group", err.Error())
		return
	}

	// update state
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createCredentialGroup(ctx context.Context, data *CredentialsGroup, created *bool) func() *helper.RetryError {
	c := r.Provider.Client()
	return func() *helper.RetryError {
		if !*created {
			if err := c.ObjectStorage.CredentialsGroup.Create(ctx, data.ProjectID.Value, data.Name.Value); err != nil {
				if common.IsNonRetryable(err) {
					return helper.NonRetryableError(err)
				}
				return helper.RetryableError(err)
			}
			x := true
			created = &x
		}
		res, err := c.ObjectStorage.CredentialsGroup.List(ctx, data.ProjectID.Value)
		if err != nil {
			if common.IsNonRetryable(err) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}

		for _, group := range res.CredentialsGroups {
			if group.DisplayName == data.Name.Value {
				data.ID = types.String{Value: group.CredentialsGroupID}
				data.URN = types.String{Value: group.URN}
				return nil
			}
		}

		y := false
		created = &y
		return nil
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.Provider.Client()
	var state CredentialsGroup

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found := false
	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		var err error
		res, err := c.ObjectStorage.CredentialsGroup.List(ctx, state.ProjectID.Value)
		if err != nil {
			if common.IsNonRetryable(err) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}
		for _, group := range res.CredentialsGroups {
			if group.DisplayName == state.Name.Value {
				found = true
				state.ID = types.String{Value: group.CredentialsGroupID}
				state.URN = types.String{Value: group.URN}
				return nil
			}
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read credential group", err.Error())
		return
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

	c := r.Provider.Client()
	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		if err := c.ObjectStorage.CredentialsGroup.Delete(ctx, state.ProjectID.Value, state.ID.Value); err != nil {
			if common.IsNonRetryable(err) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
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
