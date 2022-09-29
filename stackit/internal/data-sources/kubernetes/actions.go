package kubernetes

import (
	"context"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const default_retry_duration = 10 * time.Minute

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.Provider.Client()
	var config kubernetes.Cluster

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read cluster
	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		cl, err := c.Kubernetes.Clusters.Get(ctx, config.ProjectID.Value, config.Name.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		transform(&config, cl)
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read cluster", err.Error())
		return
	}

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
	c := r.Provider.Client()

	var cred clusters.Credentials
	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		var err error
		cred, err = c.Kubernetes.Clusters.GetCredential(ctx, cl.ProjectID.Value, cl.Name.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		diags.AddError("failed to get cluster credentials", err.Error())
		return
	}
	cl.KubeConfig = types.String{Value: cred.Kubeconfig}
}
