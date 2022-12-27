package instance

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/mongodb-flex/instances"
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
		resp.Diagnostics.AddError("failed mongodb validation", err.Error())
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
	res, wait, err := r.client.MongoDBFlex.Instances.Create(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), plan.MachineType.ValueString(), instances.Storage{
		Class: storage.Class.ValueString(),
		Size:  int(storage.Size.ValueInt64()),
	}, plan.Version.ValueString(), int(plan.Replicas.ValueInt64()), plan.BackupSchedule.ValueString(), plan.Labels, plan.Options, instances.ACL{Items: acl})

	if err != nil {
		resp.Diagnostics.AddError("failed MongoDB instance creation", err.Error())
		return
	}

	// set state
	plan.ID = types.StringValue(res.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), res.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), plan.ProjectID.ValueString())...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API currently has a bug that causes the instance to initially get a FAILED status
	// To overcome the bug, we'll wait an initial 30 sec
	time.Sleep(30 * time.Second)

	instance, err := wait.Wait()
	if err != nil {
		resp.Diagnostics.AddError("failed MongoDB instance creation validation", err.Error())
		return
	}

	i, ok := instance.(instances.Instance)
	if !ok {
		resp.Diagnostics.AddError("failed to parse client response", "response is not of instances.Instance")
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
	// these are the default user values
	// the current API doesn't read them yet, but in later releases
	// this will be the way to get the default user and database credentials
	// the default user credentials won't change
	username := "stackit"
	database := "stackit"
	roles := []string{}

	for maxTries := 10; maxTries > -1; maxTries-- {
		res, err := r.client.MongoDBFlex.Users.Create(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), username, database, roles)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) && maxTries > 0 {
				time.Sleep(time.Second * 5)
				continue
			}
			if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) && maxTries > 0 {
				time.Sleep(time.Second * 30)
				continue
			}
			d.AddError("failed to create user", err.Error())
			return
		}

		elems := []attr.Value{}
		for _, v := range res.Item.Roles {
			elems = append(elems, types.StringValue(v))
		}
		u, diags := types.ObjectValue(
			map[string]attr.Type{
				"id":       types.StringType,
				"username": types.StringType,
				"database": types.StringType,
				"password": types.StringType,
				"host":     types.StringType,
				"port":     types.Int64Type,
				"uri":      types.StringType,
				"roles":    types.ListType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"id":       types.StringValue(res.Item.ID),
				"username": types.StringValue(res.Item.Username),
				"database": types.StringValue(res.Item.Database),
				"password": types.StringValue(res.Item.Password),
				"host":     types.StringValue(res.Item.Host),
				"port":     types.Int64Value(int64(res.Item.Port)),
				"uri":      types.StringValue(res.Item.URI),
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
	instance, err := r.client.MongoDBFlex.Instances.Get(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read mongodb instance", err.Error())
		return
	}

	if err := applyClientResponse(&state, instance.Item); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
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

	// validate
	if err := r.validate(ctx, plan); err != nil {
		resp.Diagnostics.AddError("failed mongodb validation", err.Error())
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

	// handle update
	_, wait, err := r.client.MongoDBFlex.Instances.Update(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), plan.Name.ValueString(), plan.MachineType.ValueString(), instances.Storage{
		Class: storage.Class.ValueString(),
		Size:  int(storage.Size.ValueInt64()),
	}, plan.Version.ValueString(), int(plan.Replicas.ValueInt64()), plan.BackupSchedule.ValueString(), plan.Labels, plan.Options, instances.ACL{Items: acl})
	if err != nil {
		resp.Diagnostics.AddError("failed MongoDB instance update", err.Error())
		return
	}

	instance, err := wait.Wait()
	if err != nil {
		resp.Diagnostics.AddError("failed MongoDB instance update validation", err.Error())
		return
	}

	i, ok := instance.(instances.Instance)
	if !ok {
		resp.Diagnostics.AddError("failed to parse client response", "response is not of instances.Instance")
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

	process, err := r.client.MongoDBFlex.Instances.Delete(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete mongodb instance", err.Error())
		return
	}

	// allow long wait
	httpClient := r.client.GetHTTPClient()
	t := httpClient.Timeout
	httpClient.Timeout = time.Minute
	_, err = process.Wait()

	// revert
	httpClient.Timeout = t

	if err != nil {
		resp.Diagnostics.AddError("failed to verify mongodb instance deletion", err.Error())
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
			fmt.Sprintf("Expected import identifier with format: `project_id,mongodb_instance_id`.\nInstead got: %q", req.ID),
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
