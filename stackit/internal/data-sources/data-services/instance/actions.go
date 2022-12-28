package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Instance
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.Instances.ListWithResponse(ctx, config.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to prepare list instances request", err.Error())
		return
	}
	if res.HasError != nil {
		resp.Diagnostics.AddError("failed to make list instances request", res.HasError.Error())
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("failed to parse list instances response", "JSON200 == nil")
		return
	}

	list := res.JSON200
	found := -1
	existing := ""
	for i, instance := range list.Instances {
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
	instance := list.Instances[found]

	ores, err := d.client.Offerings.GetWithResponse(ctx, config.ProjectID.ValueString())
	if err != nil {
		diags.AddError(err.Error(), "failed to prepare get offerings call")
		return
	}
	if ores.HasError != nil {
		diags.AddError(res.HasError.Error(), "failed making get offerings call")
		return
	}
	if ores.JSON200 == nil {
		diags.AddError("received an empty response for offerings", "JSON200 == nil")
		return
	}

	for _, offer := range ores.JSON200.Offerings {
		for _, p := range offer.Plans {
			if p.ID != instance.PlanID {
				continue
			}
			config.Plan = types.StringValue(p.Name)
			config.Version = types.StringValue(offer.Version)
		}
	}

	if instance.InstanceID == nil {
		diags.AddError("received an empty instance ID", "InstanceID == nil")
		return
	}
	config.ID = types.StringValue(*instance.InstanceID)
	config.PlanID = types.StringValue(instance.PlanID)
	config.DashboardURL = types.StringNull()
	config.DashboardURL = types.StringValue(instance.DashboardUrl)
	config.CFGUID = types.StringValue(instance.CFGUID)
	config.CFSpaceGUID = types.StringValue(instance.CFSpaceGUID)
	config.CFOrganizationGUID = types.StringValue(instance.CFSpaceGUID)
	elems := []attr.Value{}
	if acl, ok := instance.Parameters["sgw_acl"]; ok {
		aclString, ok := acl.(string)
		if !ok {
			diags.AddError("couldn't parse ACL as string", "ACL interface isn't a string")
			return
		}
		items := strings.Split(aclString, ",")
		for _, v := range items {
			elems = append(elems, types.StringValue(v))
		}
	}
	config.ACL = types.ListValueMust(types.StringType, elems)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
