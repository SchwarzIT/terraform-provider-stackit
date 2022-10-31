package postgresinstance

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/postgres/instances"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PostgresInstance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	res, wait, err := r.client.Incubator.Postgres.Instances.Create(ctx, plan.ProjectID.Value, plan.Name.Value, plan.FlavorID.Value, instances.Storage{
		Class: plan.Storage.Class.Value,
		Size:  int(plan.Storage.Size.Value),
	}, plan.Version.Value, int(plan.Replicas.Value), plan.BackupSchedule.Value, plan.Labels, plan.Options, instances.ACL{Items: plan.ACL})

	if err != nil {
		resp.Diagnostics.AddError("failed Postgres instance creation", err.Error())
		return
	}

	// set state
	plan.ID = types.String{Value: res.ID}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := wait.Wait(); err != nil {
		resp.Diagnostics.AddError("failed Postgres instance creation validation", err.Error())
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PostgresInstance

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster
	instance, err := r.client.Incubator.Postgres.Instances.Get(ctx, state.ProjectID.Value, state.ID.Value)
	if err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read postgres instance", err.Error())
		return
	}

	_ = instance
	state.ACL = instance.Item.ACL.Items
	state.BackupSchedule = types.String{Value: instance.Item.BackupSchedule}
	state.FlavorID = types.String{Value: instance.Item.Flavor.ID}
	state.Name = types.String{Value: instance.Item.Name}
	state.Replicas = types.Int64{Value: int64(instance.Item.Replicas)}
	state.Storage = Storage{
		Class: types.String{Value: instance.Item.Storage.Class},
		Size:  types.Int64{Value: int64(instance.Item.Storage.Size)},
	}
	state.Version = types.String{Value: instance.Item.Version}
	state.Users = []User{}

	for _, user := range instance.Item.Users {
		state.Users = append(state.Users, User{
			ID:       types.String{Value: user.ID},
			Username: types.String{Value: user.Username},
			Password: types.String{Value: user.Password},
			Hostname: types.String{Value: user.Hostname},
			Database: types.String{Value: user.Database},
			Port:     types.Int64{Value: int64(user.Port)},
			URI:      types.String{Value: user.URI},
			Roles:    user.Roles,
		})
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
	var plan PostgresInstance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	_, wait, err := r.client.Incubator.Postgres.Instances.Update(ctx, plan.ProjectID.Value, plan.ID.Value, plan.FlavorID.Value, plan.BackupSchedule.Value, plan.Labels, plan.Options, instances.ACL{Items: plan.ACL})
	if err != nil {
		resp.Diagnostics.AddError("failed Postgres instance update", err.Error())
		return
	}

	if _, err := wait.Wait(); err != nil {
		resp.Diagnostics.AddError("failed Postgres instance update validation", err.Error())
		return
	}

	// update state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PostgresInstance
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	process, err := r.client.Incubator.Postgres.Instances.Delete(ctx, state.ProjectID.Value, state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete postgres instance", err.Error())
		return
	}

	if _, err := process.Wait(); err != nil {
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
