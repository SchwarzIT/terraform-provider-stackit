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

	cl, err := c.Services.Kubernetes.Cluster.GetClusterWithResponse(ctx, config.ProjectID.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to prepare read request for cluster", err.Error())
		return
	}
	if cl.HasError != nil {
		resp.Diagnostics.AddError("failed to read cluster", cl.HasError.Error())
		return
	}
	transform(&config, cl.JSON200)

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
	res, err := c.Services.Kubernetes.Credentials.GetClusterCredentialsWithResponse(ctx, cl.ProjectID.ValueString(), cl.Name.ValueString())
	if err != nil {
		diags.AddError("failed to prepare get request for cluster credentials", err.Error())
		return
	}
	if res.HasError != nil {
		diags.AddError("failed to get cluster credentials", res.HasError.Error())
		return
	}
	if res.JSON200.Kubeconfig != nil {
		cl.KubeConfig = types.StringValue(*res.JSON200.Kubeconfig)
	} else {
		cl.KubeConfig = types.StringNull()
	}
}
