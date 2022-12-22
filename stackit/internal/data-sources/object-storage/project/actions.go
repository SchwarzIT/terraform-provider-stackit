package project

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/project"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config project.ObjectStorageProject

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	c := r.client.ObjectStorage.Projects
	_, err := c.Get(ctx, config.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read request for Object Storage project", err.Error())
		return
	}

	config.ID = types.StringValue(config.ProjectID.ValueString())

	// update state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
