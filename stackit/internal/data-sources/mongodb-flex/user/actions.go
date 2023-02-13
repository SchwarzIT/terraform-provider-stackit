package instance

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client.MongoDBFlex
	var config User
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.User.GetWithResponse(ctx, config.ProjectID.ValueString(), config.InstanceID.ValueString(), config.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200.Item"); agg != nil {
		diags.AddError("failed to list mongodb instances", agg.Error())
		return
	}

	item := res.JSON200.Item
	config.Username = nullOrValStr(item.Username)
	config.Host = nullOrValStr(item.Host)
	config.Port = nullOrValInt64(item.Port)
	config.Database = nullOrValStr(item.Database)
	if roles := item.Roles; roles != nil && len(*roles) > 0 {
		r := *roles
		config.Role = types.StringValue(r[0])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func nullOrValStr(v *string) basetypes.StringValue {
	a := types.StringNull()
	if v != nil {
		a = types.StringValue(*v)
	}
	return a
}

func nullOrValInt64(v *int) basetypes.Int64Value {
	a := types.Int64Null()
	if v != nil {
		a = types.Int64Value(int64(*v))
	}
	return a
}
