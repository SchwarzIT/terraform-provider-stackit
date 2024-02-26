package network

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"strconv"
	"strings"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (d DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := d.client
	var config Network

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID, _ := uuid.Parse(config.ProjectID.ValueString())

	resNetworkList, err := c.IAAS.V1ListNetworksInProject(ctx, projectID)
	if agg := common.Validate(&resp.Diagnostics, resNetworkList, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed instance read", agg.Error())
		return
	}

	networkList := resNetworkList.JSON200.Items
	var networkID uuid.UUID
	found := false

	for _, n := range networkList {
		if strings.EqualFold(n.Name, config.Name.ValueString()) {
			networkID = n.NetworkID
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find network", "network could not be found.")
		resp.Diagnostics.Append(diags...)
		return
	}

	resNetwork, err := c.IAAS.V1GetNetwork(ctx, projectID, networkID)
	if agg := common.Validate(&resp.Diagnostics, resNetwork, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed instance read", agg.Error())
		return
	}

	network := resNetwork.JSON200

	prefixes := make([]attr.Value, 0)
	if len(network.Prefixes) > 0 {
		for _, pr := range network.Prefixes {
			prefixes = append(prefixes, types.StringValue(pr))
		}
	}

	nameservers := make([]attr.Value, 0)
	if network.Nameservers != nil && len(*network.Nameservers) > 0 {
		for _, ns := range *network.Nameservers {
			nameservers = append(nameservers, types.StringValue(ns))
		}
	}

	config.ID = types.StringValue(network.NetworkID.String())
	config.ProjectID = types.StringValue(projectID.String())
	config.Name = types.StringValue(network.Name)
	config.NetworkID = types.StringValue(network.NetworkID.String())
	config.PublicIp = types.StringPointerValue(network.PublicIp)
	config.Prefixes = types.ListValueMust(types.StringType, prefixes)
	config.NameServers = types.ListValueMust(types.StringType, nameservers)

	// get the Prefix Length in a hacky way, otherwise fall back to default
	if len(network.Prefixes) > 0 {
		cidrSplit := strings.Split(network.Prefixes[0], "/")
		if len(cidrSplit) != 2 {
			resp.Diagnostics.AddError("Processing CIDR Prefix Length",
				"Processing CIDR Prefix Length")
			return
		}

		prefixLength, err := strconv.ParseInt(cidrSplit[1], 10, 64)
		if err != nil {
			resp.Diagnostics.AddError("Processing CIDR Prefix Length", err.Error())
			return
		}

		config.PrefixLengthV4 = types.Int64Value(prefixLength)
	} else {
		config.PrefixLengthV4 = types.Int64Value(25)
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
