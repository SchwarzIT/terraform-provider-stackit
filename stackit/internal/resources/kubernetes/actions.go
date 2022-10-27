package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
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
	plan.Status = types.String{Value: consts.SKE_CLUSTER_STATUS_HEALTHY}
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
	projectID := cl.ProjectID.Value
	clusterName := cl.Name.Value
	clusterConfig, err := cl.clusterConfig(versions)
	if err != nil {
		diags.AddError("Failed to create cluster config", err.Error())
		return
	}
	nodePools := setNodepoolDefaults(cl.nodePools())
	maintenance := cl.maintenance()
	hibernations := cl.hibernations()
	extensions := cl.extensions()

	if err := r.validate(ctx, projectID, clusterName, clusterConfig, nodePools, maintenance, hibernations, extensions); err != nil {
		diags.AddError(
			"Failed cluster validation",
			err.Error(),
		)
		return
	}

	_, process, err := c.Kubernetes.Clusters.CreateOrUpdate(ctx,
		projectID,
		clusterName,
		clusterConfig,
		nodePools,
		maintenance,
		hibernations,
		extensions,
	)
	if err != nil {
		diags.AddError("failed during SKE Create/Update", err.Error())
		return
	}

	res, err := process.Wait()
	if err != nil {
		diags.AddError("failed to validate SKE Create/Update", err.Error())
		return
	}

	result, ok := res.(clusters.Cluster)
	if !ok {
		diags.AddError("failed to parse Wait() response", "response is not clusters.Cluster")
		return
	}

	cl.Status = types.String{Value: result.Status.Aggregated}
	cl.Transform(result)
}

func (r Resource) getCredential(ctx context.Context, diags *diag.Diagnostics, cl *Cluster) {
	c := r.client
	cred, err := c.Kubernetes.Clusters.GetCredential(ctx, cl.ProjectID.Value, cl.Name.Value)
	if err != nil {
		diags.AddError("failed to get cluster credentials", err.Error())
		return
	}
	cl.KubeConfig = types.String{Value: cred.Kubeconfig}
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
	cl, err := c.Kubernetes.Clusters.Get(ctx, state.ProjectID.Value, state.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to read cluster", err.Error())
		return
	}
	state.Transform(cl)

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
	process, err := c.Kubernetes.Clusters.Delete(ctx, state.ProjectID.Value, state.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete cluster", err.Error())
		return
	}

	if _, err := process.Wait(); err != nil {
		resp.Diagnostics.AddError("failed to verify cluster deletion", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,name` where `name` is the cluster name.\nInstead got: %q", req.ID),
		)
		return
	}

	// validate cluster name
	if err := clusters.ValidateClusterName(idParts[1]); err != nil {
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
			fmt.Sprintf("Couldn't validate project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}

	// pre-read imports
	c := r.client
	res, err := c.Kubernetes.Clusters.Get(ctx, idParts[0], idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error during import", err.Error())
		return
	}

	if res.Extensions != nil {
		extensions := &Extensions{}

		if res.Extensions.Argus != nil {
			extensions.Argus = &ArgusExtension{
				Enabled:         types.Bool{Value: res.Extensions.Argus.Enabled},
				ArgusInstanceID: types.String{Value: res.Extensions.Argus.ArgusInstanceID},
			}
		}

		diags := resp.State.SetAttribute(ctx, path.Root("extensions"), extensions)
		resp.Diagnostics.Append(diags...)
	}

	if res.Hibernation != nil {
		hibernations := []Hibernation{}
		for _, h := range res.Hibernation.Schedules {
			hibernations = append(hibernations, Hibernation{
				Start:    types.String{Value: h.Start},
				End:      types.String{Value: h.End},
				Timezone: types.String{Value: h.Timezone},
			})
		}
		diags := resp.State.SetAttribute(ctx, path.Root("hibernations"), hibernations)
		resp.Diagnostics.Append(diags...)
	}

	if res.Maintenance != nil {
		digas := resp.State.SetAttribute(ctx, path.Root("maintenance"), &Maintenance{
			EnableKubernetesVersionUpdates:   types.Bool{Value: res.Maintenance.AutoUpdate.KubernetesVersion},
			EnableMachineImageVersionUpdates: types.Bool{Value: res.Maintenance.AutoUpdate.MachineImageVersion},
			Start:                            types.String{Value: res.Maintenance.TimeWindow.Start},
			End:                              types.String{Value: res.Maintenance.TimeWindow.End},
		})
		resp.Diagnostics.Append(digas...)
	}
}
