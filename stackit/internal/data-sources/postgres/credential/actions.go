package credential

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Credential
	diags := req.Config.Get(ctx, &config)
	es := d.client.DataServices.PostgresDB

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	i, err := es.Credentials.Get(ctx, config.ProjectID.ValueString(), config.InstanceID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get credential", err.Error())
		return
	}

	// set computed fields
	config.CACert = types.StringValue(i.Raw.Credential.Cacrt)
	config.Host = types.StringValue(i.Raw.Credential.Host)
	if len(i.Raw.Credential.Hosts) > 0 {
		var d diag.Diagnostics
		if len(i.Raw.Credential.Hosts) > 0 {
			config.Hosts, d = types.ListValueFrom(ctx, types.StringType, i.Raw.Credential.Hosts)
			resp.Diagnostics.Append(d...)
		}
	}
	config.Username = types.StringValue(i.Raw.Credential.Username)
	config.Password = types.StringValue(i.Raw.Credential.Password)
	config.Port = types.Int64Value(int64(i.Raw.Credential.Port))
	config.SyslogDrainURL = types.StringValue(i.Raw.SyslogDrainURL)
	config.RouteServiceURL = types.StringValue(i.Raw.RouteServiceURL)
	config.Schema = types.StringValue(i.Raw.Credential.Scheme)
	config.URI = types.StringValue(i.URI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
