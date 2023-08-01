package user

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config User

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.SecretsManager.Users.List(ctx, uuid.MustParse(config.ProjectID.ValueString()), uuid.MustParse(config.InstanceID.ValueString()))
	if agg := common.Validate(&resp.Diagnostics, res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to list users", agg.Error())
		return
	}

	for _, user := range res.JSON200.Users {
		if user.Username == config.Username.ValueString() {
			config.ID = types.StringValue(user.ID)
			config.Description = types.StringValue(user.Description)
			config.Write = types.BoolValue(user.Write)
			break
		}
	}

	if config.ID.IsNull() || config.ID.IsUnknown() {
		resp.Diagnostics.AddError("failed to find user", "user not found")
		return
	}

	// update config
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
