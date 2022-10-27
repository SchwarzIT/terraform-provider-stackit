package project

import (
	"context"

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

	project, err := c.ResourceManagement.Projects.Get(ctx, p.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to read project", err.Error())
		return
	}

	p.ID = types.String{Value: project.ProjectID}
	p.Name = types.String{Value: project.Name}

	if billing, ok := project.Labels["billingReference"]; ok {
		p.BillingRef = types.String{Value: billing}
	}

	diags = resp.State.Set(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
