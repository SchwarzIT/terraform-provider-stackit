package network

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (d DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := d.client
	var n Network

	diags := req.Config.Get(ctx, &n)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID, _ := uuid.Parse(n.ProjectID.String())
	networkID, _ := uuid.Parse(n.NetworkID.String())
	res, err := c.IAAS.V1GetNetwork(ctx, projectID, networkID)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed reading project", agg.Error())
		return
	}

	network := res.JSON200
	prefixes := make([]types.String, len(network.Prefixes))
	for i, pr := range network.Prefixes {
		prefixes[i] = types.StringValue(pr)
	}
	n.NetworkID = types.StringValue(network.NetworkID.String())
	n.PublicIp = types.StringValue(*network.PublicIp)
	n.Prefixes = prefixes
	n.Name = types.StringValue(network.Name)
	n.ProjectID = types.StringValue(projectID.String())

	diags = resp.State.Set(ctx, &n)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
