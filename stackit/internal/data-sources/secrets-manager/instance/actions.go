package instance

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Instance

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.SecretsManager.Instances.Get(ctx, uuid.MustParse(config.ProjectID.ValueString()), uuid.MustParse(config.ID.ValueString()))
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to get instance", agg.Error())
		return
	}

	config.Name = types.StringValue(res.JSON200.Name)
	config.Frontend = types.StringValue(res.JSON200.ApiUrl + "/ui")
	config.API = types.StringValue(res.JSON200.ApiUrl)

	r.readACLs(ctx, &config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// update config
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r DataSource) readACLs(ctx context.Context, config *Instance, diags *diag.Diagnostics) {
	c := r.client
	res, err := c.SecretsManager.Acls.List(ctx, uuid.MustParse(config.ProjectID.ValueString()), uuid.MustParse(config.ID.ValueString()))
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		diags.AddError("failed to get instance ACLs", agg.Error())
		return
	}
	els := []attr.Value{}
	for _, el := range res.JSON200.Acls {
		els = append(els, types.StringValue(el.Cidr))
	}
	config.ACL = types.SetValueMust(types.StringType, els)
	return
}
