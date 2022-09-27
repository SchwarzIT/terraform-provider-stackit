package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.Provider.IsConfigured() {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on another resource.",
		)
		return
	}

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
	c := r.Provider.Client()
	created := false

	// cluster vars
	projectID := cl.ProjectID.Value
	clusterName := cl.Name.Value
	clusterConfig := cl.clusterConfig()
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

	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		// Create cluster
		if !created {
			_, err := c.Kubernetes.Clusters.Create(ctx,
				projectID,
				clusterName,
				clusterConfig,
				nodePools,
				maintenance,
				hibernations,
				extensions,
			)

			if err != nil {
				if common.IsNonRetryable(err) {
					return helper.NonRetryableError(err)
				}
				return helper.RetryableError(err)
			}
			created = true
		}

		// Check cluster status
		s, err := c.Kubernetes.Clusters.Get(ctx, cl.ProjectID.Value, cl.Name.Value)
		if err != nil {
			if common.IsNonRetryable(err) {
				return helper.NonRetryableError(fmt.Errorf("error receiving cluster state: %s", err))
			}
			return helper.RetryableError(err)
		}

		cl.Status = types.String{Value: s.Status.Aggregated}
		if s.Status.Aggregated != consts.SKE_CLUSTER_STATUS_HEALTHY {
			return helper.RetryableError(fmt.Errorf("expected cluster to be active & healthy but it was in state %s", s.Status.Aggregated))
		}

		cl.Transform(s)
		return nil
	}); err != nil {
		diags.AddError("failed to verify cluster", err.Error())
		return
	}
}

func (r Resource) getCredential(ctx context.Context, diags *diag.Diagnostics, cl *Cluster) {
	c := r.Provider.Client()

	var cred clusters.Credentials
	if err := helper.RetryContext(ctx, common.DURATION_5M, func() *helper.RetryError {
		var err error
		cred, err = c.Kubernetes.Clusters.GetCredential(ctx, cl.ProjectID.Value, cl.Name.Value)
		if err != nil {
			if common.IsNonRetryable(err) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		diags.AddError("failed to get cluster credentials", err.Error())
		return
	}
	cl.KubeConfig = types.String{Value: cred.Kubeconfig}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.Provider.Client()
	var state, plan Cluster

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster
	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		cl, err := c.Kubernetes.Clusters.Get(ctx, state.ProjectID.Value, state.Name.Value)
		if err != nil {
			if common.IsNonRetryable(err) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}
		state.Transform(cl)
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read cluster", err.Error())
		return
	}

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

	c := r.Provider.Client()
	deleted := false
	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		if !deleted {
			if err := c.Kubernetes.Clusters.Delete(ctx, state.ProjectID.Value, state.Name.Value); err != nil {
				return helper.RetryableError(err)
			}
			deleted = true
		}

		// Verify cluster deletion
		s, err := c.Kubernetes.Clusters.Get(ctx, state.ProjectID.Value, state.Name.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return nil
			}
			return helper.RetryableError(err)
		}
		return helper.RetryableError(fmt.Errorf("expected cluster to be deleted, but was it was in state %s", s.Status.Aggregated))
	}); err != nil {
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
	c := r.Provider.Client()
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
