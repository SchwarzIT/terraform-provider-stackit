package project

import (
	"context"
	"fmt"

	rmv2 "github.com/SchwarzIT/community-stackit-go-client/pkg/services/resource-management/v2.0"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
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

	projectType := rmv2.PROJECT
	res, err := c.ResourceManagement.GetContainersOfAnOrganization(ctx, p.ContainerID.ValueString(), &rmv2.GetContainersOfAnOrganizationParams{Type: &projectType})
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed reading project with container ID: %s", p.ContainerID.ValueString()), agg.Error())
		return
	}

	containers := *res.JSON200
	id := -1
	for i, project := range containers.Items {
		if project.Item.ContainerID != p.ContainerID.ValueString() {
			continue
		}
		id = i
	}
	if id == -1 {
		resp.Diagnostics.AddError("not found", "project container ID not found")
	}
	project := containers.Items[id].Item
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
	}
	diags = resp.State.Set(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
