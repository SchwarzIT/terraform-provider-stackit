package credential

import (
	"context"
	"net/http"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Credential
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.Credentials.GetWithResponse(ctx, config.ProjectID.ValueString(), config.InstanceID.ValueString(), config.ID.ValueString())
	if agg := validate.Response(res, err, "JSON200.Raw"); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read credential", agg.Error())
		return
	}

	i := res.JSON200
	// set computed fields
	config.Host = types.StringValue(i.Raw.Credentials.Host)
	config.Hosts = types.ListNull(types.StringType)
	if len(i.Raw.Credentials.Hosts) > 0 {
		h := []attr.Value{}
		for _, v := range i.Raw.Credentials.Hosts {
			h = append(h, types.StringValue(v))
		}
		config.Hosts = types.ListValueMust(types.StringType, h)
	}

	config.DatabaseName = types.StringValue(i.Raw.Credentials.Name)
	config.Username = types.StringValue(i.Raw.Credentials.Username)
	config.Password = types.StringValue(i.Raw.Credentials.Password)
	config.Port = types.Int64Value(int64(i.Raw.Credentials.Port))
	config.SyslogDrainURL = types.StringValue(i.Raw.SyslogDrainUrl)
	config.RouteServiceURL = types.StringValue(i.Raw.RouteServiceUrl)
	config.URI = types.StringValue(i.Uri)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
