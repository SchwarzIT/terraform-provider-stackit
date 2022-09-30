package bucket

import (
	"context"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/buckets"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/bucket"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		var err error
		b, err = c.ObjectStorage.Buckets.Get(ctx, config.ProjectID.Value, config.Name.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
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
