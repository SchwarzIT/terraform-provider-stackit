package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/secrets-manager/v1.1.0/users"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	res, err := c.SecretsManager.Users.Create(ctx, uuid.MustParse(plan.ProjectID.ValueString()), uuid.MustParse(plan.InstanceID.ValueString()), users.UserCreate{
		Description: plan.Description.ValueString(),
		Write:       plan.Write.ValueBool(),
	})
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to create user", agg.Error())
		return
	}

	plan.ID = types.StringValue(res.JSON200.ID)
	plan.Username = types.StringValue(res.JSON200.Username)
	plan.Password = types.StringValue(res.JSON200.Password)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var state User

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.SecretsManager.Users.Get(ctx, uuid.MustParse(state.ProjectID.ValueString()), uuid.MustParse(state.InstanceID.ValueString()), uuid.MustParse(state.ID.ValueString()))
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		if res != nil && res.StatusCode() == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to get user", agg.Error())
		return
	}

	state.Description = types.StringValue(res.JSON200.Description)
	state.Write = types.BoolValue(res.JSON200.Write)
	state.Username = types.StringValue(res.JSON200.Username)

	// update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	write := plan.Write.ValueBool()
	_, err := c.SecretsManager.Users.Update(ctx, uuid.MustParse(plan.ProjectID.ValueString()), uuid.MustParse(plan.InstanceID.ValueString()), uuid.MustParse(plan.ID.ValueString()), users.UserUpdate{
		Write: &write,
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to update user", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("write_enabled"), write)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state User
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	_, err := c.SecretsManager.Users.Delete(ctx, uuid.MustParse(state.ProjectID.ValueString()), uuid.MustParse(state.InstanceID.ValueString()), uuid.MustParse(state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("failed to delete user", err.Error())
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
			fmt.Sprintf("Expected import identifier with format: `project_id,instance_id,user_id`\nInstead got: %q", req.ID),
		)
		return
	}

	// validate project id
	if err := clientValidate.ProjectID(idParts[0]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate kubernetes_project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("password"), "")...)

	if resp.Diagnostics.HasError() {
		return
	}

}
