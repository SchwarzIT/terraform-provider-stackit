package network

import (
	"context"
	"fmt"
	iaas "github.com/SchwarzIT/community-stackit-go-client/pkg/services/iaas-api/v1alpha"
	"github.com/google/uuid"
	"reflect"

	openapiTypes "github.com/SchwarzIT/community-stackit-go-client/pkg/helpers/types"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Network
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	plan = r.createNetwork(ctx, resp, plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createNetwork(ctx context.Context, resp *resource.CreateResponse, plan Network) Network {
	//labels := rmv2.ResourceLabels{
	//	"billingReference": plan.BillingRef.ValueString(),
	//	"scope":            "PUBLIC",
	//}

	//for k, v := range plan.Labels {
	//	labels[k] = v
	//}

	//owner := rmv2.PROJECT_OWNER
	//subj1 := r.client.Client.GetServiceAccountEmail()
	//subj2 := plan.OwnerEmail.ValueString()
	//members := []rmv2.ProjectMember{
	//	{
	//		Subject: &subj1,
	//		Role:    &owner,
	//	},
	//	{
	//		Subject: &subj2,
	//		Role:    &owner,
	//	},
	//}

	var ns iaas.V1Nameserver
	for _, i := range plan.NameServers {
		ns = append(ns, i.String())
	}

	pl := int(plan.PrefixLengthV4.ValueInt64())
	body := iaas.V1CreateNetworkJSONRequestBody{
		Name:           plan.Name.String(),
		Nameservers:    &ns,
		PrefixLengthV4: &pl,
	}
	res, err := r.client.IAAS.V1CreateNetwork(ctx, openapiTypes.UUID{}, body)

	//body := rmv2.ProjectRequestBody{
	//	ContainerParentID: plan.ParentContainerID.ValueString(),
	//	Labels:            &labels,
	//	Members:           members,
	//	Name:              plan.Name.ValueString(),
	//}
	//res, err := r.client.ResourceManagement.Create(ctx, body)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON201"); agg != nil {
		resp.Diagnostics.AddError("failed creating network", agg.Error())
		return plan
	}

	process := res.WaitHandler(ctx, r.client.IAAS, openapiTypes.UUID{}, openapiTypes.UUID{})
	createdNetwork, err := process.WaitWithContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed validating network %s creation", body.Name), err.Error())
		return plan
	}

	networks := createdNetwork.(*iaas.V1ListNetworksInProjectResponse)
	for _, n := range networks.JSON200.Items {
		if plan.Name.String() != n.Name {
			continue
		}

		prefixes := make([]types.String, len(n.Prefixes))
		for i, pr := range n.Prefixes {
			prefixes[i] = types.StringValue(pr)
		}
		plan.NetworkID = types.StringValue(n.NetworkID.String())
		plan.PublicIp = types.StringValue(*n.PublicIp)
		plan.Prefixes = prefixes
	}

	return plan
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var p Network

	diags := req.State.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.IAAS.V1ListNetworksInProject(ctx, openapiTypes.UUID{})
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed reading project", agg.Error())
		return
	}

	for _, n := range res.JSON200.Items {
		if n.Name != p.Name.String() {
			continue
		}
		prefixes := make([]types.String, len(n.Prefixes))
		for i, pr := range n.Prefixes {
			prefixes[i] = types.StringValue(pr)
		}
		p.NetworkID = types.StringValue(n.NetworkID.String())
		p.PublicIp = types.StringValue(*n.PublicIp)
		p.Prefixes = prefixes
		p.Name = types.StringValue(n.Name)
		p.NetworkID = types.StringValue(n.NetworkID.String())
	}

	diags = resp.State.Set(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Network
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state Network
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.NetworkID.IsUnknown() {
		plan.NetworkID = state.NetworkID
	}

	r.updateNetwork(ctx, plan, state, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) updateNetwork(ctx context.Context, projectID openapiTypes.UUID, plan, state Network, resp *resource.UpdateResponse) {
	if plan.Name.Equal(state.Name) && reflect.DeepEqual(plan.NameServers, state.NameServers) {
		return
	}

	ns := make([]iaas.V1IP, len(plan.NameServers))
	for i, s := range plan.NameServers {
		ns[i] = s.String()
	}
	n := plan.Name.String()
	body := iaas.V1UpdateNetworkJSONBody{
		Name:        &n,
		Nameservers: &ns,
	}

	id, err := uuid.Parse(state.NetworkID.String())
	if err != nil {
		return
	}

	res, err := r.client.IAAS.V1UpdateNetwork(ctx, projectID, id, iaas.V1UpdateNetworkJSONRequestBody(body))
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed updating project", agg.Error())
		return
	}
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Network
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := uuid.Parse(state.NetworkID.String())
	if err != nil {
		return
	}

	c := r.client
	res, err := c.IAAS.V1DeleteNetwork(ctx, openapiTypes.UUID{}, id)
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		resp.Diagnostics.AddError("failed deleting network", agg.Error())
		return
	}

	process := res.WaitHandler(ctx, c.IAAS, openapiTypes.UUID{}, state.NetworkID).SetTimeout(timeout)
	res, err := process.WaitWithContext(ctx); err != nil {
		resp.Diagnostics.AddError("failed to verify network deletion", err.Error())
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// validate project id
	if err := clientValidate.NetworkID(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate network_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)
}
