package project

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config KubernetesProject

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	c := r.client.Kubernetes.Project
	p, err := c.Get(ctx, config.ProjectID.ValueString())
	if agg := common.Validate(&resp.Diagnostics, p, err); agg != nil {
		resp.Diagnostics.AddError("failed to read SKE project", agg.Error())
		return
	}

	config.ID = types.StringValue(config.ProjectID.ValueString())

	// update state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
