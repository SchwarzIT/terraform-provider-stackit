package cluster

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/generated/cluster"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Cluster
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	r.createOrUpdateCluster(ctx, &resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle credential
	r.getCredential(ctx, &resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	plan.Status = types.StringValue(string(cluster.STATE_HEALTHY))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createOrUpdateCluster(ctx context.Context, diags *diag.Diagnostics, cl *Cluster) {
	c := r.client
	versions, err := r.loadAvaiableVersions(ctx)
	if err != nil {
		diags.AddError("failed while loading version options", err.Error())
		return
	}

	// cluster vars
	projectID := cl.KubernetesProjectID.ValueString()
	clusterName := cl.Name.ValueString()
	clusterConfig, err := cl.clusterConfig(versions)
	if err != nil {
		diags.AddError("Failed to create cluster config", err.Error())
		return
	}
	nodePools := setNodepoolDefaults(cl.nodePools())
	maintenance := cl.maintenance()
	hibernations := cl.hibernations()
	extensions, diag := cl.extensions(ctx)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	if err := r.validate(ctx, projectID, clusterName, clusterConfig, &nodePools, maintenance, hibernations, extensions); err != nil {
		diags.AddError(
			"Failed cluster validation",
			err.Error(),
		)
		return
	}

	resp, err := c.Kubernetes.Cluster.CreateOrUpdateClusterWithResponse(ctx,
		projectID,
		clusterName, cluster.SkeServiceCreateOrUpdateClusterRequest{
			Extensions:  extensions,
			Hibernation: hibernations,
			Kubernetes:  clusterConfig,
			Maintenance: maintenance,
			Nodepools:   nodePools,
		},
	)
	if agg := validate.Response(resp, err); agg != nil {
		diags.AddError("failed during SKE create/update", agg.Error())
	}

	process := resp.WaitHandler(ctx, c.Kubernetes.Cluster, projectID, clusterName)
	res, err := process.WaitWithContext(ctx)
	if agg := validate.Response(res, err, "JSON200.Status.Aggregated"); agg != nil {
		diags.AddError("failed to validate SKE create/update", agg.Error())
		return
	}
	result, ok := res.(*cluster.GetClusterResponse)
	if !ok {
		diags.AddError("failed to parse Wait() response", "response is not *cluster.GetClusterResponse")
		return
	}
	cl.Status = types.StringValue(string(*result.JSON200.Status.Aggregated))
	cl.Transform(*result.JSON200)
}

func (r Resource) getCredential(ctx context.Context, diags *diag.Diagnostics, cl *Cluster) {
	c := r.client
	res, err := c.Kubernetes.Credentials.GetClusterCredentialsWithResponse(ctx, cl.KubernetesProjectID.ValueString(), cl.Name.ValueString())
	if agg := validate.Response(res, err, "JSON200.Kubeconfig"); agg != nil {
		diags.AddError("failed fetching cluster credentials", agg.Error())
		return
	}

	cl.KubeConfig = types.StringValue(*res.JSON200.Kubeconfig)
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var state Cluster

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster
	res, err := c.Kubernetes.Cluster.GetClusterWithResponse(ctx, state.KubernetesProjectID.ValueString(), state.Name.ValueString())
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed fetching cluster", agg.Error())
		return
	}

	state.Transform(*res.JSON200)

	// read credential
	r.getCredential(ctx, &resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Cluster
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	r.createOrUpdateCluster(ctx, &resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle credential
	r.getCredential(ctx, &resp.Diagnostics, &plan)
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

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	res, err := c.Kubernetes.Cluster.DeleteClusterWithResponse(ctx, state.KubernetesProjectID.ValueString(), state.Name.ValueString())
	if agg := validate.Response(res, err); agg != nil {
		resp.Diagnostics.AddError("failed deleting cluster", agg.Error())
		return
	}

	if _, err := res.WaitHandler(ctx, c.Kubernetes.Cluster, state.KubernetesProjectID.ValueString(), state.Name.ValueString()).WaitWithContext(ctx); err != nil {
		if !strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			resp.Diagnostics.AddError("failed to verify cluster deletion", err.Error())
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `kubernetes_project_id,name` where `name` is the cluster name.\nInstead got: %q", req.ID),
		)
		return
	}

	// validate cluster name
	if err := cluster.ValidateClusterName(idParts[1]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate cluster name.\n%s", err.Error()),
		)
		return
	}

	// validate project id
	if err := clientValidate.ProjectID(idParts[0]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate kubernetes_project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("kubernetes_project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}

	// pre-read imports
	c := r.client
	res, err := c.Kubernetes.Cluster.GetClusterWithResponse(ctx, idParts[0], idParts[1])
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed import pre-read", agg.Error())
		return
	}

	if res.JSON200.Extensions != nil {
		extensions := &Extensions{}

		if res.JSON200.Extensions.Argus != nil {
			extensions.Argus = &ArgusExtension{
				Enabled:         types.BoolValue(res.JSON200.Extensions.Argus.Enabled),
				ArgusInstanceID: types.StringValue(res.JSON200.Extensions.Argus.ArgusInstanceID),
			}
		}

		diags := resp.State.SetAttribute(ctx, path.Root("extensions"), extensions)
		resp.Diagnostics.Append(diags...)
	}

	if res.JSON200.Hibernation != nil {
		hibernations := []Hibernation{}
		for _, h := range res.JSON200.Hibernation.Schedules {
			hibernations = append(hibernations, Hibernation{
				Start:    types.StringValue(h.Start),
				End:      types.StringValue(h.End),
				Timezone: types.StringValue(*h.Timezone),
			})
		}
		diags := resp.State.SetAttribute(ctx, path.Root("hibernations"), hibernations)
		resp.Diagnostics.Append(diags...)
	}

	if res.JSON200.Maintenance != nil {
		digas := resp.State.SetAttribute(ctx, path.Root("maintenance"), &Maintenance{
			EnableKubernetesVersionUpdates:   types.BoolValue(*res.JSON200.Maintenance.AutoUpdate.KubernetesVersion),
			EnableMachineImageVersionUpdates: types.BoolValue(*res.JSON200.Maintenance.AutoUpdate.MachineImageVersion),
			Start:                            types.StringValue(res.JSON200.Maintenance.TimeWindow.Start),
			End:                              types.StringValue(res.JSON200.Maintenance.TimeWindow.End),
		})
		resp.Diagnostics.Append(digas...)
	}
}
