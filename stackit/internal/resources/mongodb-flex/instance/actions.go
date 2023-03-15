package instance

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/mongodb-flex/v1.0/generated/instance"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
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
			Class: types.StringValue(DefaultStorageClass),
			Size:  types.Int64Value(DefaultStorageSize),
		}
	} else {
		resp.Diagnostics.Append(plan.Storage.As(ctx, &storage, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	cl := storage.Class.ValueString()
	sz := int(storage.Size.ValueInt64())

	// handle creation
	bus := plan.BackupSchedule.ValueString()
	fid := plan.MachineType.ValueString()
	name := plan.Name.ValueString()
	rpl := int(plan.Replicas.ValueInt64())
	ver := plan.Version.ValueString()
	optType := plan.Type.ValueString()
	body := instance.InstanceCreateInstanceRequest{
		ACL: &instance.InstanceACL{Items: &acl},
		Storage: &instance.InstanceStorage{
			Class: &cl,
			Size:  &sz,
		},
		BackupSchedule: &bus,
		FlavorID:       &fid,
		Labels:         &plan.Labels,
		Name:           &name,
		Options: &struct {
			Type *string `json:"type,omitempty"`
		}{
			Type: &optType,
		},
		Replicas: &rpl,
		Version:  &ver,
	}

	res, err := r.client.MongoDBFlex.Instance.CreateWithResponse(ctx, plan.ProjectID.ValueString(), body)
	if agg := validate.Response(res, err, "JSON202.ID"); agg != nil {
		resp.Diagnostics.AddError("failed MongoDB flex instance creation", agg.Error())
		return
	}

	// set state

	instanceID := *res.JSON202.ID
	plan.ID = types.StringValue(instanceID)
	if instanceID == "" {
		resp.Diagnostics.AddError("received an empty instance ID", fmt.Sprintf("invalid instance id: %+v", *res.JSON202))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), instanceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), plan.ProjectID.ValueString())...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API currently has a bug that causes the instance to initially get a FAILED status
	// To overcome the bug, we'll wait an initial 30 sec
	time.Sleep(60 * time.Second)

	process := res.WaitHandler(ctx, r.client.MongoDBFlex.Instance, plan.ProjectID.ValueString(), instanceID)
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed MongoDB instance creation validation", err.Error())
		return
	}

	// read cluster
	get, err := r.client.MongoDBFlex.Instance.GetWithResponse(ctx, plan.ProjectID.ValueString(), instanceID)
	if agg := validate.Response(get, err, "JSON200.Item"); agg != nil {
		resp.Diagnostics.AddError("failed to get instance after create", agg.Error())
		return
	}

	if err := applyClientResponse(&plan, get.JSON200.Item); err != nil {
		resp.Diagnostics.AddError("failed to process client response", err.Error())
		return
	}

	// update state with user
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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
	res, err := r.client.MongoDBFlex.Instance.GetWithResponse(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200.Item"); agg != nil {
		resp.Diagnostics.AddError("failed to read instance", agg.Error())
		return
	}
	if err := applyClientResponse(&state, res.JSON200.Item); err != nil {
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
			Class: types.StringValue(DefaultStorageClass),
			Size:  types.Int64Value(DefaultStorageSize),
		}
	} else {
		resp.Diagnostics.Append(plan.Storage.As(ctx, &storage, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	cl := storage.Class.ValueString()
	sz := int(storage.Size.ValueInt64())

	// handle creation
	bus := plan.BackupSchedule.ValueString()
	fid := plan.MachineType.ValueString()
	name := plan.Name.ValueString()
	rpl := int(plan.Replicas.ValueInt64())
	ver := plan.Version.ValueString()
	body := instance.InstanceUpdateInstanceRequest{
		ACL: &instance.InstanceACL{Items: &acl},
		Storage: &instance.InstanceStorage{
			Class: &cl,
			Size:  &sz,
		},
		BackupSchedule: &bus,
		FlavorID:       &fid,
		Labels:         &plan.Labels,
		Name:           &name,
		Options: &map[string]string{
			"type": plan.Type.ValueString(),
		},
		Replicas: &rpl,
		Version:  &ver,
	}

	// handle update
	res, err := r.client.MongoDBFlex.Instance.PatchWithResponse(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString(), body)
	if agg := validate.Response(res, err, "JSON202.Item"); agg != nil {
		resp.Diagnostics.AddError("failed updating mongodb flex instance", agg.Error())
		return
	}

	process := res.WaitHandler(ctx, r.client.MongoDBFlex.Instance, plan.ProjectID.ValueString(), plan.ID.ValueString())
	if _, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed MongoDB instance update validation", err.Error())
		return
	}

	// read cluster
	get, err := r.client.MongoDBFlex.Instance.GetWithResponse(ctx, plan.ProjectID.ValueString(), plan.ID.ValueString())
	if agg := validate.Response(get, err, "JSON200.Item"); agg != nil {
		resp.Diagnostics.AddError("failed to get instance after create", agg.Error())
		return
	}

	if err := applyClientResponse(&plan, get.JSON200.Item); err != nil {
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

	res, err := r.client.MongoDBFlex.Instance.DeleteWithResponse(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to delete instance", agg.Error())
		return
	}

	process := res.WaitHandler(ctx, r.client.MongoDBFlex.Instance, state.ProjectID.ValueString(), state.ID.ValueString())
	if _, err = process.WaitWithContext(ctx); err != nil {
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
