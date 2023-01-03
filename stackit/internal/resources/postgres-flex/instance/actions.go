package postgresinstance

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/generated/instance"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/generated/users"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.setDefaults()

	// validate
	if err := r.validate(ctx, plan); err != nil {
		resp.Diagnostics.AddError("failed postgres validation", err.Error())
		return
	}

	acl := []string{}
	for _, v := range plan.ACL.Elements() {
		nv, err := common.ToString(context.Background(), v)
		if err != nil {
			continue
		}
		acl = append(acl, nv)
	}

	storage := Storage{}
	if plan.Storage.IsUnknown() {
		storage = Storage{
			Class: types.StringValue(default_storage_class),
			Size:  types.Int64Value(default_storage_size),
		}
	} else {
		resp.Diagnostics.Append(plan.Storage.As(ctx, &storage, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// handle creation

	c := r.client.PostgresFlex

	// prepare values
	name := plan.Name.ValueString()
	bu := plan.BackupSchedule.ValueString()
	flavorID := plan.MachineType.ValueString()
	repl := int(plan.Replicas.ValueInt64())
	sc := storage.Class.ValueString()
	ss := int(storage.Size.ValueInt64())
	v := plan.Version.ValueString()

	body := instance.InstanceCreateInstanceRequest{
		Name: &name,
		ACL: &instance.InstanceACL{
			Items: &acl,
		},
		BackupSchedule: &bu,
		FlavorID:       &flavorID,
		Labels:         &plan.Labels,
		Options:        &plan.Options,
		Replicas:       &repl,
		Storage: &instance.InstanceStorage{
			Class: &sc,
			Size:  &ss,
		},
		Version: &v,
	}
	res, err := c.Instance.CreateWithResponse(ctx, plan.ProjectID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("failed preparing Postgres instance creation request", err.Error())
		return
	}
	if res.HasError != nil {
		resp.Diagnostics.AddError("failed Postgres instance creation", res.HasError.Error())
		return
	}
	if res.JSON200.ID == nil {
		resp.Diagnostics.AddError("received nil ID", "postgres flex API returned a nil ID")
	}

	// set state
	plan.ID = types.StringValue(*res.JSON200.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), *res.JSON200.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), plan.ProjectID.ValueString())...)
	if resp.Diagnostics.HasError() {
		return
	}

	process := res.WaitHandler(ctx, c.Instance, plan.ProjectID.ValueString(), *res.JSON200.ID)
	ins, err := process.WaitWithContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed Postgres instance creation validation", err.Error())
		return
	}

	i, ok := ins.(*instance.InstanceSingleInstance)
	if !ok {
		resp.Diagnostics.AddError("failed to parse client response", "response is not of *instance.InstanceSingleInstance")
		return
	}

	if err := applyClientResponse(&plan, i); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	r.createUser(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state with user
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createUser(ctx context.Context, plan *Instance, d *diag.Diagnostics) {
	// TODO: remove this code hack when the API returns
	// a proper instance status
	time.Sleep(time.Minute * 2)

	// these are the default user values
	// the current API doesn't read them yet, but in later releases
	// this will be the way to get the default user and database credentials
	// the default user credentials won't change
	username := "stackit"
	database := "stackit"
	roles := []string{}

	for maxTries := 10; maxTries > -1; maxTries-- {
		c := r.client.PostgresFlex
		body := users.CreateUserJSONRequestBody{
			Database: &database,
			Roles:    &roles,
			Username: &username,
		}
		res, err := c.Users.CreateUserWithResponse(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), body)
		if err != nil {
			d.AddError("failed prepare create user request", err.Error())
			return
		}
		if (res.StatusCode() == http.StatusNotFound ||
			res.StatusCode() == http.StatusBadRequest) &&
			maxTries > 0 {
			time.Sleep(time.Second * 5)
			continue
		}
		if res.HasError != nil {
			d.AddError("failed create user request", res.HasError.Error())
		}

		if res.JSON200 == nil {
			d.AddError("response is nil", "api returned nil response")
			return
		}

		item := users.InstanceUser{}
		if res.JSON200.Item != nil {
			item = *res.JSON200.Item
		}
		roles := []string{}
		if item.Roles != nil {
			roles = *item.Roles
		}
		elems := []attr.Value{}
		for _, v := range roles {
			elems = append(elems, types.StringValue(v))
		}

		id := ""
		if item.ID != nil {
			id = *item.ID
		}
		un := ""
		if item.Username != nil {
			un = *item.Username
		}
		db := ""
		if item.Database != nil {
			db = *item.Database
		}
		pw := ""
		if item.Password != nil {
			pw = *item.Password
		}
		ho := ""
		if item.Host != nil {
			ho = *item.Host
		}
		var po int
		if item.Port != nil {
			po = *item.Port
		}
		uri := ""
		if item.Uri != nil {
			uri = *item.Uri
		}
		u, diags := types.ObjectValue(
			map[string]attr.Type{
				"id":       types.StringType,
				"username": types.StringType,
				"database": types.StringType,
				"password": types.StringType,
				"hostname": types.StringType,
				"port":     types.Int64Type,
				"uri":      types.StringType,
				"roles":    types.ListType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"id":       types.StringValue(id),
				"username": types.StringValue(un),
				"database": types.StringValue(db),
				"password": types.StringValue(pw),
				"hostname": types.StringValue(ho),
				"port":     types.Int64Value(int64(po)),
				"uri":      types.StringValue(uri),
				"roles":    types.ListValueMust(types.StringType, elems),
			},
		)
		plan.User = u
		d.Append(diags...)
		break
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Instance

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster

	c := r.client.PostgresFlex
	res, err := c.Instance.GetWithResponse(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to prepare read postgres instance request", err.Error())
		return
	}
	if res.HasError != nil {
		if res.StatusCode() == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read postgres instance", err.Error())
		return
	}

	if res.JSON200 == nil {
		resp.Diagnostics.AddError("instance response is nil", "JSON200 is nil")
		return
	}

	if err := applyClientResponse(&state, res.JSON200.Item); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	if state.User.IsNull() || state.User.IsUnknown() {
		r.createUser(ctx, &state, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
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
	var plan, state Instance
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID

	// validate
	if err := r.validate(ctx, plan); err != nil {
		resp.Diagnostics.AddError("failed postgres validation", err.Error())
		return
	}

	acl := []string{}
	for _, v := range plan.ACL.Elements() {
		nv, err := common.ToString(context.Background(), v)
		if err != nil {
			continue
		}
		acl = append(acl, nv)
	}

	storage := Storage{}
	if plan.Storage.IsUnknown() {
		storage = Storage{
			Class: types.StringValue(default_storage_class),
			Size:  types.Int64Value(default_storage_size),
		}
	} else {
		resp.Diagnostics.Append(plan.Storage.As(ctx, &storage, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// prepare values
	name := plan.Name.ValueString()
	bu := plan.BackupSchedule.ValueString()
	flavorID := plan.MachineType.ValueString()
	repl := int(plan.Replicas.ValueInt64())
	sc := storage.Class.ValueString()
	ss := int(storage.Size.ValueInt64())
	v := plan.Version.ValueString()

	body := instance.InstanceUpdateInstanceRequest{
		Name: &name,
		ACL: &instance.InstanceACL{
			Items: &acl,
		},
		BackupSchedule: &bu,
		FlavorID:       &flavorID,
		Labels:         &plan.Labels,
		Options:        &plan.Options,
		Replicas:       &repl,
		Storage: &instance.InstanceStorage{
			Class: &sc,
			Size:  &ss,
		},
		Version: &v,
	}

	// handle update
	c := r.client.PostgresFlex.Instance
	res, err := c.UpdateWithResponse(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("failed prepare Postgres instance update request", err.Error())
		return
	}
	if res.HasError != nil {
		if res.StatusCode() == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to update postgres instance", err.Error())
		return
	}

	process := res.WaitHandler(ctx, c, plan.ProjectID.ValueString(), plan.ID.ValueString())
	isi, err := process.WaitWithContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed Postgres instance update validation", err.Error())
		return
	}

	i, ok := isi.(*instance.InstanceSingleInstance)
	if !ok {
		resp.Diagnostics.AddError("failed to parse client response", "response is not of *instance.InstanceSingleInstance")
		return
	}

	if err := applyClientResponse(&plan, i); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	// update state
	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Instance
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client.PostgresFlex.Instance
	res, err := c.DeleteWithResponse(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to prepare delete postgres instance request", err.Error())
		return
	}
	if res.HasError != nil {
		if res.StatusCode() == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to make delete postgres instance request", res.HasError.Error())
		return
	}

	process := res.WaitHandler(ctx, c, state.ProjectID.ValueString(), state.ID.ValueString())
	if _, err := process.WaitWithContext(ctx); err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to verify postgres instance deletion", err.Error())
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
			fmt.Sprintf("Expected import identifier with format: `project_id,postgres_instance_id`.\nInstead got: %q", req.ID),
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}
