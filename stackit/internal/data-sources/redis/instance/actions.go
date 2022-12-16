package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Instance
	diags := req.Config.Get(ctx, &config)
	rds := c.DataServices.Redis

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := rds.Instances.List(ctx, config.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to list instances", err.Error())
		return
	}

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
	intance := list.Instances[found]

	res, err := rds.Options.GetOfferings(ctx, config.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get offerings", err.Error())
		return
	}

	for _, offer := range res.Offerings {
		for _, p := range offer.Plans {
			if p.ID != intance.PlanID {
				continue
			}
			config.Plan = types.StringValue(p.Name)
			config.Version = types.StringValue(offer.Version)
		}
	}
	config.ID = types.StringValue(intance.InstanceID)
	config.PlanID = types.StringValue(intance.PlanID)
	config.DashboardURL = types.StringValue(intance.DashboardURL)
	config.CFGUID = types.StringValue(intance.CFGUID)
	config.CFSpaceGUID = types.StringValue(intance.CFSpaceGUID)
	config.CFOrganizationGUID = types.StringValue(intance.CFOrganizationGUID)
	config.ACL = types.List{ElemType: types.StringType}
	if aclString, ok := intance.Parameters["sgw_acl"]; ok {
		items := strings.Split(aclString, ",")
		for _, v := range items {
			config.ACL.Elems = append(config.ACL.Elems, types.StringValue(v))
		}
	} else {
		config.ACL.Null = true
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
