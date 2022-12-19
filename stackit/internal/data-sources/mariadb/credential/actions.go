package credential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (d DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Credential
	diags := req.Config.Get(ctx, &config)
	service := d.client.DataServices.MariaDB

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	i, err := service.Credentials.Get(ctx, config.ProjectID.ValueString(), config.InstanceID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get credential", err.Error())
		return
	}

	// set computed fields
	config.CACert = types.StringValue(i.Raw.Credential.Cacrt)
	config.Host = types.StringValue(i.Raw.Credential.Host)
	config.Hosts = types.List{ElemType: types.StringType}
	if len(i.Raw.Credential.Hosts) > 0 {
		config.Hosts.Elems = make([]attr.Value, len(i.Raw.Credential.Hosts))
		for k, v := range i.Raw.Credential.Hosts {
			config.Hosts.Elems[k] = types.StringValue(v)
		}
	}
	config.Username = types.StringValue(i.Raw.Credential.Username)
	config.Password = types.StringValue(i.Raw.Credential.Password)
	config.Port = types.Int64Value(int64(i.Raw.Credential.Port))
	config.SyslogDrainURL = types.StringValue(i.Raw.SyslogDrainURL)
	config.RouteServiceURL = types.StringValue(i.Raw.RouteServiceURL)
	config.URI = types.StringValue(i.URI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
