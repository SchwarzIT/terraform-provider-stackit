package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
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
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to list instances", agg.Error())
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
	if agg := validate.Response(ores, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to get offerings", agg.Error())
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
		resp.Diagnostics.AddError("received an empty instance ID", "InstanceID == nil")
		return
	}
	config.ID = types.StringValue(*instance.InstanceID)
	config.PlanID = types.StringValue(instance.PlanID)
	config.DashboardURL = types.StringNull()
	config.DashboardURL = types.StringValue(instance.DashboardUrl)
	config.CFGUID = types.StringValue(instance.CFGUID)
	config.CFSpaceGUID = types.StringValue(instance.CFSpaceGUID)
	config.CFOrganizationGUID = types.StringValue("")
	if instance.OrganizationGUID != nil {
		config.CFOrganizationGUID = types.StringValue(*instance.OrganizationGUID)
	}
	elems := []attr.Value{}
	if acl, ok := instance.Parameters["sgw_acl"]; ok {
		aclString, ok := acl.(string)
		if !ok {
			resp.Diagnostics.AddError("couldn't parse ACL as string", "ACL interface isn't a string")
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
