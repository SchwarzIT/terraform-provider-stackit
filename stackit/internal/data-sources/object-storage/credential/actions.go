package credential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Credential

	diags := req.Config.Get(ctx, &config)

	if config.DisplayName.Value == "" && config.ID.Value == "" {
		diags.AddError("missing configuration", "either display_name or id must be provided")
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := c.ObjectStorage.AccessKeys.List(ctx, config.ProjectID.Value, "")
	if err != nil {
		resp.Diagnostics.AddError("failed to read credential", err.Error())
		return
	}

	found := false
	for _, k := range list.AccessKeys {
		if k.KeyID != config.ID.Value && k.DisplayName != config.DisplayName.Value {
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
