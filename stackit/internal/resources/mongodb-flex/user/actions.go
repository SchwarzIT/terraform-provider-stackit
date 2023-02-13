package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/mongodb-flex/v1.0/generated/user"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := plan.Username.ValueString()
	database := plan.Database.ValueString()
	roles := []string{plan.Role.ValueString()}

	body := user.InstanceCreateUserRequest{
		Database: database,
		Roles:    roles,
		Username: &username,
	}
	res, err := r.client.MongoDBFlex.User.CreateWithResponse(ctx, plan.ProjectID.ValueString(), plan.InstanceID.ValueString(), body)
	if agg := validate.Response(res, err, "JSON202.Item"); agg != nil {
		resp.Diagnostics.AddError("failed creating mongodb flex db user", agg.Error())
		return
	}

	item := *res.JSON202.Item
	elems := []attr.Value{}
	if *res.JSON202.Item.Roles != nil {
		for _, v := range *res.JSON202.Item.Roles {
			elems = append(elems, types.StringValue(v))
		}
	}
	if res.JSON202.Item.Password == nil {
		resp.Diagnostics.AddError("received an empty password", fmt.Sprintf("full response: %+v", res.JSON202))
		return
	}

	if item.ID == nil {
		resp.Diagnostics.AddError("received an empty ID", fmt.Sprintf("full response: %+v", res.JSON202))
		return
	}
	plan.ID = nullOrValStr(item.ID)
	plan.Password = nullOrValStr(item.Password)
	plan.Host = nullOrValStr(item.Host)
	plan.Port = nullOrValInt64(item.Port)
	plan.URI = nullOrValStr(item.Uri)

	// update state with user
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state User

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster
	res, err := r.client.MongoDBFlex.User.GetWithResponse(ctx, state.ProjectID.ValueString(), state.InstanceID.ValueString(), state.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200.Item"); agg != nil {
		resp.Diagnostics.AddError("failed making read user request", err.Error())
		return
	}

	item := res.JSON200.Item
	state.Username = nullOrValStr(item.Username)
	state.Host = nullOrValStr(item.Host)
	state.Port = nullOrValInt64(item.Port)
	state.Database = nullOrValStr(item.Database)
	if roles := item.Roles; roles != nil && len(*roles) > 0 {
		r := *roles
		state.Role = types.StringValue(r[0])
	}
	if state.URI.IsUnknown() {
		state.URI = types.StringNull()
	}
	if state.Password.IsUnknown() {
		state.Password = types.StringNull()
	}

	// update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state User
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.MongoDBFlex.User.DeleteWithResponse(ctx, state.ProjectID.ValueString(), state.InstanceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed making MongoDB user deletion request", err.Error())
		return
	}
	if res.HasError != nil {
		if res.StatusCode() == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("user deletion response has an error", res.HasError.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,mongodb_instance_id,user_id`.\nInstead got: %q", req.ID),
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

func nullOrValStr(v *string) basetypes.StringValue {
	a := types.StringNull()
	if v != nil {
		a = types.StringValue(*v)
	}
	return a
}

func nullOrValInt64(v *int) basetypes.Int64Value {
	a := types.Int64Null()
	if v != nil {
		a = types.Int64Value(int64(*v))
	}
	return a
}
