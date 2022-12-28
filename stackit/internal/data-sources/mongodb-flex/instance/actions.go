package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client.MongoDBFlex
	var config Instance
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := c.Instances.List(ctx, config.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to list instances", err.Error())
		return
	}

	found := -1
	existing := ""
	for i, instance := range list.Items {
		if instance.Name == config.Name.ValueString() {
			found = i
			break
		}
		if existing == "" {
			existing = "\navailable instances in the project are:"
		}
		existing = fmt.Sprintf("%s\n- %s", existing, instance.Name)
	}

	if found == -1 {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find instance", "instance could not be found."+existing)
		resp.Diagnostics.Append(diags...)
		return
	}

	// set found instance
	instance := list.Items[found]
	ires, err := c.Instances.Get(ctx, config.ProjectID.ValueString(), instance.ID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get instances", err.Error())
		return
	}
	i := ires.Item
	config.ID = types.StringValue(instance.ID)
	elems := []attr.Value{}
	for _, v := range i.ACL.Items {
		elems = append(elems, types.StringValue(v))
	}
	config.ACL = types.ListValueMust(types.StringType, elems)
	config.BackupSchedule = types.StringValue(i.BackupSchedule)
	config.MachineType = types.StringValue(i.Flavor.ID)
	config.Name = types.StringValue(i.Name)
	config.Replicas = types.Int64Value(int64(i.Replicas))
	storage, d := types.ObjectValue(
		map[string]attr.Type{
			"class": types.StringType,
			"size":  types.Int64Type,
		},
		map[string]attr.Value{
			"class": types.StringValue(i.Storage.Class),
			"size":  types.Int64Value(int64(i.Storage.Size)),
		})

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Storage = storage

	if len(i.Version) > 3 {
		i.Version = i.Version[0:3]
	}
	config.Version = types.StringValue(i.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
