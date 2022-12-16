package credential

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/credentials"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r Resource) applyClientResponse(ctx context.Context, c *Credential, cgr credentials.GetResponse) error {
	c.ID = types.StringValue(cgr.ID)
	c.CACert = types.StringValue(cgr.Raw.Credential.Cacrt)
	c.Host = types.StringValue(cgr.Raw.Credential.Host)
	c.Hosts = types.List{ElemType: types.StringType}
	if len(cgr.Raw.Credential.Hosts) > 0 {
		c.Hosts.Elems = make([]attr.Value, len(cgr.Raw.Credential.Hosts))
		for k, v := range cgr.Raw.Credential.Hosts {
			c.Hosts.Elems[k] = types.StringValue(v)
		}
	}
	c.Username = types.StringValue(cgr.Raw.Credential.Username)
	c.Password = types.StringValue(cgr.Raw.Credential.Password)
	c.Port = types.Int64Value(int64(cgr.Raw.Credential.Port))
	c.SyslogDrainURL = types.StringValue(cgr.Raw.SyslogDrainURL)
	c.RouteServiceURL = types.StringValue(cgr.Raw.RouteServiceURL)
	c.Schema = types.StringValue(cgr.Raw.Credential.Scheme)
	c.URI = types.StringValue(cgr.URI)
	return nil
}
