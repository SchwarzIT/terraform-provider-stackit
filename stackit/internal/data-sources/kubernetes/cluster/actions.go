package cluster

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Cluster

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cl, err := c.Kubernetes.Cluster.Get(ctx, config.KubernetesProjectID.ValueString(), config.Name.ValueString())
	if agg := validate.Response(cl, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to read cluster", agg.Error())
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

func (r DataSource) getCredential(ctx context.Context, diags *diag.Diagnostics, cl *Cluster) {
	c := r.client
	res, err := c.Kubernetes.Credentials.List(ctx, cl.KubernetesProjectID.ValueString(), cl.Name.ValueString())
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		diags.AddError("failed to get cluster credentials", agg.Error())
		return
	}
	if res.JSON200.Kubeconfig != nil {
		cl.KubeConfig = types.StringValue(*res.JSON200.Kubeconfig)
	} else {
		cl.KubeConfig = types.StringNull()
	}
}
