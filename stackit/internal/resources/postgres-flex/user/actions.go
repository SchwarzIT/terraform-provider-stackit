package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/generated/users"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pkg/errors"
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
	if username == "" {
		username = "psqluser"
	}

	var roles []string
	resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roles, true)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(roles) == 0 {
		roles = []string{"login"}
	}

	body := users.CreateUserJSONRequestBody{
		Roles:    &roles,
		Username: &username,
	}

	res, err := r.client.PostgresFlex.Users.CreateUserWithResponse(ctx, plan.ProjectID.ValueString(), plan.InstanceID.ValueString(), body)
	if agg := validate.Response(res, err, "JSON201.Item"); agg != nil {
		if res.StatusCode() == http.StatusBadRequest {
			j := ""
			if res.JSON400 != nil {
				b, _ := json.Marshal(res.JSON400)
				j = string(b)
			}
			resp.Diagnostics.AddError("failed creating postgres flex db user", errors.Wrapf(agg, j).Error())
			return
		}
		resp.Diagnostics.AddError("failed creating postgres flex db user", agg.Error())
		return
	}

	item := *res.JSON201.Item

	elems := []attr.Value{}
	if *item.Roles != nil {
		for _, v := range *item.Roles {
			elems = append(elems, types.StringValue(v))
		}
	}
	if item.Password == nil {
		resp.Diagnostics.AddError("received an empty password", fmt.Sprintf("full response: %+v", item))
		return
	}

	if item.ID == nil {
		resp.Diagnostics.AddError("received an empty ID", fmt.Sprintf("full response: %+v", item))
		return
	}
	plan.ID = nullOrValStr(item.ID)
	plan.Password = nullOrValStr(item.Password)
	plan.Host = nullOrValStr(item.Host)
	plan.Port = nullOrValInt64(item.Port)
	plan.URI = nullOrValStr(item.URI)
	plan.Roles = types.ListValueMust(types.StringType, elems)

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
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster
	res, err := r.client.PostgresFlex.Users.GetUserWithResponse(ctx, state.ProjectID.ValueString(), state.InstanceID.ValueString(), state.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200.Item"); agg != nil {
		if res.JSON400 != nil {
			// verify the instance exists
			res, err := r.client.PostgresFlex.Instance.ListWithResponse(ctx, state.ProjectID.ValueString())
			if agg2 := validate.Response(res, err, "JSON200.Items"); agg2 != nil {
				resp.Diagnostics.AddError("failed making read user request", agg.Error())
				resp.Diagnostics.AddError("failed verifying instance status", agg2.Error())
				return
			}
			for _, item := range *res.JSON200.Items {
				if item.ID != nil && *item.ID == state.InstanceID.ValueString() {
					resp.Diagnostics.AddError("failed making read user request", agg.Error())
					return
				}
			}
			// instance doesn't exists:
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed making read user request", agg.Error())
		return
	}

	item := res.JSON200.Item
	state.Username = nullOrValStr(item.Username)
	state.Host = nullOrValStr(item.Host)
	state.Port = nullOrValInt64(item.Port)
	roles := []attr.Value{}
	if r := item.Roles; r != nil {
		for _, v := range *r {
			roles = append(roles, types.StringValue(v))
		}
	}
	state.Roles = types.ListValueMust(types.StringType, roles)
	if state.URI.IsUnknown() {
		state.URI = types.StringNull()
	}
	if state.Password.IsUnknown() {
		state.Password = types.StringNull()
	}

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	res, err := r.client.PostgresFlex.Users.DeleteUserWithResponse(ctx, state.ProjectID.ValueString(), state.InstanceID.ValueString(), state.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to delete user", agg.Error())
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
			fmt.Sprintf("Expected import identifier with format: `project_id,postgres_flex_instance_id,user_id`.\nInstead got: %q", req.ID),
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
