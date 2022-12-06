package credential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config Credential
	diags := req.Config.Get(ctx, &config)
	es := r.client.DataServices.ElasticSearch

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	i, err := es.Credentials.Get(ctx, config.ProjectID.Value, config.InstanceID.Value, config.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to get credential", err.Error())
		return
	}

	// set computed fields
	config.CACert = types.String{Value: i.Raw.Credential.Cacrt}
	config.Host = types.String{Value: i.Raw.Credential.Host}

	config.Hosts = types.List{ElemType: types.StringType}
	if len(i.Raw.Credential.Hosts) > 0 {
		config.Hosts.Elems = make([]attr.Value, len(i.Raw.Credential.Hosts))
		for k, v := range i.Raw.Credential.Hosts {
			config.Hosts.Elems[k] = types.String{Value: v}
		}
	}

	config.Username = types.String{Value: i.Raw.Credential.Username}
	config.Password = types.String{Value: i.Raw.Credential.Password}
	config.Port = types.Int64{Value: int64(i.Raw.Credential.Port)}
	config.SyslogDrainURL = types.String{Value: i.Raw.SyslogDrainURL}
	config.RouteServiceURL = types.String{Value: i.Raw.RouteServiceURL}
	config.Schema = types.String{Value: i.Raw.Credential.Scheme}
	config.URI = types.String{Value: i.URI}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
