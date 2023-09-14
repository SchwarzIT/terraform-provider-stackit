package loadbalancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg Instance
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.LoadBalancer.Instances.Get(ctx, cfg.ProjectID.ValueString(), cfg.Name.ValueString())
	if agg := validate.Response(res, err, "JSON200.Name"); agg != nil {
		if res.StatusCode() == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Couldn't get instance information", agg.Error())
		return
	}
	if res.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Couldn't get instance information", fmt.Sprintf("Received status code %d", res.StatusCode()))
		common.Dump(&resp.Diagnostics, res.Body)
		return
	}

	cfg.parse(ctx, *res.JSON200, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
