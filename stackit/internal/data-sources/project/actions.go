package project

import (
	"context"

	rmv2 "github.com/SchwarzIT/community-stackit-go-client/pkg/services/resource-management/v2.0"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var p Project

	diags := req.Config.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.ResourceManagement.Get(ctx, p.ContainerID.ValueString(), &rmv2.GetParams{})
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed reading project", agg.Error())
		return
	}

	project := res.JSON200
	p.ID = types.StringValue(project.ProjectID.String())
	p.Name = types.StringValue(project.Name)
	p.ParentContainerID = types.StringValue(project.Parent.ContainerID)
	p.ContainerID = types.StringValue(project.ContainerID)
	p.BillingRef = types.StringNull()
	if project.Labels != nil {
		l := *project.Labels
		if v, ok := l["billingReference"]; ok {
			p.BillingRef = types.StringValue(v)
		}
		for k, v := range l {
			p.Labels[k] = v
		}
	}
	diags = resp.State.Set(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
