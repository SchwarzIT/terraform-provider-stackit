package bucket

import (
	"context"
	"net/http"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Bucket

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.ObjectStorage.Bucket.Get(ctx, config.ObjectStorageProjectID.ValueString(), config.Name.ValueString())
	if agg := validate.Response(res, err, "JSON200.Bucket"); agg != nil {
		resp.Diagnostics.AddError("failed to read bucket", agg.Error())
		return
	}

	if res.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	b := res.JSON200
	config.ID = types.StringValue(b.Bucket.Name)
	config.Region = types.StringValue(b.Bucket.Region)
	config.HostStyleURL = types.StringValue(b.Bucket.UrlVirtualHostedStyle)
	config.PathStyleURL = types.StringValue(b.Bucket.UrlPathStyle)
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
