package credential

import (
	"context"

	accesskey "github.com/SchwarzIT/community-stackit-go-client/pkg/services/object-storage/v1.0.1/access-key"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Credential

	diags := req.Config.Get(ctx, &config)

	if config.DisplayName.ValueString() == "" && config.ID.ValueString() == "" {
		diags.AddError("missing configuration", "either display_name or id must be provided")
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	credGroup := config.CredentialsGroupID.ValueString()
	params := &accesskey.GetParams{
		CredentialsGroup: &credGroup,
	}
	res, err := c.ObjectStorage.AccessKey.Get(ctx, config.ProjectID.ValueString(), params)
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200.AccessKeys"); agg != nil {
		resp.Diagnostics.AddError("failed to list credentials", agg.Error())
		return
	}

	list := res.JSON200
	found := false
	for _, k := range list.AccessKeys {
		if k.KeyID != config.ID.ValueString() && k.DisplayName != config.DisplayName.ValueString() {
			continue
		}

		found = true
		config.ID = types.StringValue(k.KeyID)
		config.DisplayName = types.StringValue(k.DisplayName)
		config.Expiry = types.StringValue(k.Expires)
		diags = resp.State.Set(ctx, &config)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		break
	}

	if !found {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find credential", "credential could not be found")
		resp.Diagnostics.Append(diags...)
		return
	}
}
