package credentialsgroup

import (
	"context"

	credentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credentials-group"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var data credentialsGroup.CredentialsGroup

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := c.ObjectStorage.CredentialsGroup.List(ctx, data.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to list credentials group", err.Error())
		return
	}

	found := false
	for _, group := range list.CredentialsGroups {
		if group.CredentialsGroupID == data.ID.ValueString() {
			found = true
			data.Name = types.StringValue(group.DisplayName)
			data.URN = types.StringValue(group.URN)
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("not found", "credential group could not be found")
		return
	}

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
