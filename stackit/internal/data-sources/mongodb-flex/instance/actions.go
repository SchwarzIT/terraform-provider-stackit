package instance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/mongodb-flex/v1.0/instance"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

	res, err := c.Instance.List(ctx, config.ProjectID.ValueString(), &instance.ListParams{})
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200.Items"); agg != nil {
		diags.AddError("failed to list mongodb instances", agg.Error())
		return
	}

	list := *res.JSON200.Items
	found := -1
	existing := ""
	for i, instance := range list {
		if instance.Name == nil || instance.ID == nil {
			continue
		}
		if strings.EqualFold(*instance.Name, config.Name.ValueString()) {
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
	ires, err := c.Instance.Get(ctx, config.ProjectID.ValueString(), *instance.ID)
	if agg := common.Validate(&resp.Diagnostics, ires, err, "JSON200.Item"); agg != nil {
		resp.Diagnostics.AddError("failed to get mongodb instance", agg.Error())
		return
	}

	if err := ApplyClientResponse(&config, ires.JSON200.Item); err != nil {
		resp.Diagnostics.AddError("error during apply", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ApplyClientResponse(pi *Instance, i *instance.InstancesSingleInstance) error {
	elems := []attr.Value{}
	if i == nil {
		return errors.New("instance response is empty")
	}
	if i.ACL != nil && i.ACL.Items != nil {
		for _, v := range *i.ACL.Items {
			// only include correctly formatted CIDR range
			// this is to overcome a current bug in the API
			if strings.Contains(v, "/") {
				elems = append(elems, types.StringValue(v))
			}
		}
	}
	if i.ID == nil {
		return errors.New("received a nil ID")
	}
	pi.ID = types.StringValue(*i.ID)
	pi.ACL = types.ListValueMust(types.StringType, elems)
	pi.BackupSchedule = nullOrValStr(i.BackupSchedule)
	pi.MachineType = nullOrValStr(i.Flavor.ID)
	pi.Name = nullOrValStr(i.Name)
	pi.Replicas = nullOrValInt64(i.Replicas)
	storage, diags := types.ObjectValue(
		map[string]attr.Type{
			"class": types.StringType,
			"size":  types.Int64Type,
		},
		map[string]attr.Value{
			"class": nullOrValStr(i.Storage.Class),
			"size":  nullOrValInt64(i.Storage.Size),
		})
	if diags.HasError() {
		return errors.New("failed setting storage object")
	}
	pi.Storage = storage
	pi.Version = nullOrValStr(i.Version)
	if !pi.Version.IsNull() && len(pi.Version.ValueString()) > 3 {
		v := pi.Version.ValueString()
		pi.Version = types.StringValue(v[0:3])
	}
	return nil
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
