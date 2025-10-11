package cluster

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.1/credentials"
	serviceenablement "github.com/SchwarzIT/community-stackit-go-client/pkg/services/service-enablement/v1"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.1/cluster"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
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

	timeout, d := plan.Timeouts.Create(ctx, 1*time.Hour)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	// pre process plan
	r.preProcessConfig(&resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle project init
	r.enableProject(ctx, &resp.Diagnostics, &plan, timeout)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	r.createOrUpdateCluster(ctx, &resp.Diagnostics, &plan, timeout)
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

func (r Resource) preProcessConfig(diags *diag.Diagnostics, cl *Cluster) {
	projectID := cl.ProjectID.ValueString()
	kubernetesProjectID := cl.KubernetesProjectID.ValueString()
	if projectID == "" && kubernetesProjectID == "" {
		diags.AddError("project_id or kubernetes_project_id must be set", "please note that kubernetes_project_id is deprecated and will be removed in a future release")
		return
	}
	if projectID == "" {
		cl.ProjectID = cl.KubernetesProjectID
	}
	if kubernetesProjectID == "" {
		cl.KubernetesProjectID = cl.ProjectID
	}
}

func (r Resource) enableProject(ctx context.Context, diags *diag.Diagnostics, cl *Cluster, timeout time.Duration) {
	projectID := cl.ProjectID.ValueString()
	c := r.client.ServiceEnablement

	serviceID := "cloud.stackit.ske"

	found := true
	status, err := c.GetService(ctx, projectID, serviceID)
	if agg := common.Validate(&diag.Diagnostics{}, status, err, "JSON200.State"); agg != nil {
		if status == nil || status.StatusCode() != http.StatusNotFound {
			diags.AddError("failed to fetch SKE project status", agg.Error())
			return
		}
		found = false
	}

	// check if project already enabled
	if found && *status.JSON200.State == serviceenablement.ENABLED {
		return
	}

	res, err := c.EnableService(ctx, projectID, serviceID)
	if agg := common.Validate(diags, res, err); agg != nil {
		diags.AddError("failed during SKE service enablement", agg.Error())
		return
	}

	process := res.WaitHandler(ctx, c, projectID, serviceID).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
		diags.AddError("failed to verify SKE service enablement", err.Error())
		return
	}
}

func (r Resource) createOrUpdateCluster(ctx context.Context, diags *diag.Diagnostics, cl *Cluster, timeout time.Duration) {
	c := r.client
	versions, err := r.loadAvaiableVersions(ctx, diags)
	if err != nil {
		diags.AddError("failed while loading version options", err.Error())
		return
	}

	// cluster vars
	projectID := cl.ProjectID.ValueString()
	clusterName := cl.Name.ValueString()
	clusterConfig, err := cl.clusterConfig(versions)
	networkID := cl.NetworkID.ValueString()
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

	if err := r.validate(ctx, diags, projectID, clusterName, clusterConfig, &nodePools, maintenance, hibernations, extensions); err != nil {
		diags.AddError(
			"Failed cluster validation",
			err.Error(),
		)
		return
	}

	clusterData := cluster.SkeServiceCreateOrUpdateClusterRequest{
		Extensions:  extensions,
		Hibernation: hibernations,
		Kubernetes:  clusterConfig,
		Maintenance: maintenance,
		Nodepools:   nodePools,
	}

	if networkID != "" {
		clusterData.Network = &cluster.Network{
			ID: &networkID,
		}
	}

	resp, err := c.Kubernetes.Cluster.CreateOrUpdate(ctx, projectID, clusterName, clusterData)
	if agg := common.Validate(diags, resp, err); agg != nil {
		diags.AddError("failed during SKE create/update", agg.Error())
		return
	}

	process := resp.WaitHandler(ctx, c.Kubernetes.Cluster, projectID, clusterName).SetTimeout(timeout)
	res, err := process.WaitWithContext(ctx)
	if agg := common.Validate(diags, res, err, "JSON200.Status.Aggregated"); agg != nil {
		diags.AddError("failed to validate SKE create/update", agg.Error())
		return
	}
	result, ok := res.(*cluster.GetResponse)
	if !ok {
		diags.AddError("failed to parse Wait() response", "response is not *cluster.GetClusterResponse")
		return
	}
	cl.Status = types.StringValue(string(*result.JSON200.Status.Aggregated))
	cl.Transform(*result.JSON200)
}

func (r Resource) getCredential(ctx context.Context, diags *diag.Diagnostics, cl *Cluster) {
	c := r.client

	// TODO: allow this to be changed later
	expirationSeconds := "2678400"

	res, err := c.Kubernetes.Credentials.CreateKubeconfig(ctx, cl.ProjectID.ValueString(), cl.Name.ValueString(), credentials.CreateKubeconfigJSONRequestBody{
		ExpirationSeconds: &expirationSeconds,
	})

	if agg := common.Validate(diags, res, err, "JSON200.Kubeconfig"); agg != nil {
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

	// pre process state
	r.preProcessConfig(&resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster
	res, err := c.Kubernetes.Cluster.Get(ctx, state.ProjectID.ValueString(), state.Name.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
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

	// pre process state
	r.preProcessConfig(&resp.Diagnostics, &plan)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Update(ctx, 1*time.Hour)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	r.createOrUpdateCluster(ctx, &resp.Diagnostics, &plan, timeout)
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

	// pre process plan
	r.preProcessConfig(&resp.Diagnostics, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	res, err := c.Kubernetes.Cluster.Delete(ctx, state.ProjectID.ValueString(), state.Name.ValueString())
	if agg := common.Validate(&resp.Diagnostics, res, err); agg != nil {
		resp.Diagnostics.AddError("failed deleting cluster", agg.Error())
		return
	}

	timeout, d := state.Timeouts.Delete(ctx, 1*time.Hour)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	process := res.WaitHandler(ctx, c.Kubernetes.Cluster, state.ProjectID.ValueString(), state.Name.ValueString()).SetTimeout(timeout)
	if _, err := process.WaitWithContext(ctx); err != nil {
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}

	// pre-read imports
	c := r.client
	res, err := c.Kubernetes.Cluster.Get(ctx, idParts[0], idParts[1])
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
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
