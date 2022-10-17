package bucket

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/buckets"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/bucket"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.Provider.Client()
	var config bucket.Bucket
	var b buckets.BucketResponse

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	b, err := c.ObjectStorage.Buckets.Get(ctx, config.ProjectID.Value, config.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to read bucket", err.Error())
		return
	}

	config.ID = types.String{Value: b.Bucket.Name}
	config.Region = types.String{Value: b.Bucket.Region}
	config.HostStyleURL = types.String{Value: b.Bucket.URLVirtualHostedStyle}
	config.PathStyleURL = types.String{Value: b.Bucket.URLPathStyle}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
