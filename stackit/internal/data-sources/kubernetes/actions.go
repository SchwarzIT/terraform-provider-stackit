package kubernetes

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config kubernetes.Cluster

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cl, err := c.Kubernetes.Clusters.Get(ctx, config.ProjectID.Value, config.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to read cluster", err.Error())
		return
	}
	transform(&config, cl)

	// read credential
	r.getCredential(ctx, &resp.Diagnostics, &config)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r DataSource) getCredential(ctx context.Context, diags *diag.Diagnostics, cl *kubernetes.Cluster) {
	c := r.client
	cred, err := c.Kubernetes.Clusters.GetCredential(ctx, cl.ProjectID.Value, cl.Name.Value)
	if err != nil {
		diags.AddError("failed to get cluster credentials", err.Error())
		return
	}
	cl.KubeConfig = types.StringValue(cred.Kubeconfig)
}
