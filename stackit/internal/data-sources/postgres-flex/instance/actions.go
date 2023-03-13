package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := d.client.PostgresFlex.Instance
	var config Instance
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.ListWithResponse(ctx, config.ProjectID.ValueString())
	if agg := validate.Response(res, err, "JSON200.Items"); agg != nil {
		resp.Diagnostics.AddError("failed to list postgres flex instances", agg.Error())
		return
	}

	list := *res.JSON200.Items
	found := -1
	existing := ""
	for i, instance := range list {
		if instance.Name == nil {
			continue
		}
		if *instance.Name == config.Name.ValueString() {
			found = i
			break
		}
		if existing == "" {
			existing = "\navailable instances in the project are:"
		}
		existing = fmt.Sprintf("%s\n- %s", existing, *instance.Name)
	}

	if found == -1 {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find instance", "instance could not be found."+existing)
		resp.Diagnostics.Append(diags...)
		return
	}

	// set found instance
	instance := list[found]
	if instance.ID == nil {
		resp.Diagnostics.AddError("received a nil instance ID", "instance.ID == nil")
		return
	}
	ires, err := c.GetWithResponse(ctx, config.ProjectID.ValueString(), *instance.ID)
	if agg := validate.Response(ires, err, "JSON200.Item"); agg != nil {
		resp.Diagnostics.AddError("failed to get postgres flex instance", agg.Error())
		return
	}

	i := *ires.JSON200.Item
	config.ID = types.StringValue(*instance.ID)

	elems := []attr.Value{}
	if i.ACL != nil && i.ACL.Items != nil {
		for _, v := range *i.ACL.Items {
			elems = append(elems, types.StringValue(v))
		}
	}
	config.ACL = types.ListValueMust(types.StringType, elems)
	config.BackupSchedule = types.StringNull()
	if i.BackupSchedule != nil {
		config.BackupSchedule = types.StringValue(*i.BackupSchedule)
	}
	config.MachineType = types.StringNull()
	if i.Flavor != nil && i.Flavor.ID != nil {
		config.MachineType = types.StringValue(*i.Flavor.ID)
	}
	config.Name = types.StringNull()
	if i.Name != nil {
		config.Name = types.StringValue(*i.Name)
	}
	config.Replicas = types.Int64Null()
	if i.Replicas != nil {
		config.Replicas = types.Int64Value(int64(*i.Replicas))
	}
	if i.Storage != nil {
		class := types.StringNull()
		if i.Storage.Class != nil {
			class = types.StringValue(*i.Storage.Class)
		}
		size := types.Int64Null()
		if i.Storage.Class != nil {
			size = types.Int64Value(int64(*i.Storage.Size))
		}
		storage, diags := types.ObjectValue(
			map[string]attr.Type{
				"class": types.StringType,
				"size":  types.Int64Type,
			},
			map[string]attr.Value{
				"class": class,
				"size":  size,
			})
		if diags.HasError() {
			return
		}
		config.Storage = storage
	}
	config.Version = types.StringNull()
	if i.Version != nil {
		config.Version = types.StringValue(*i.Version)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
