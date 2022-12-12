package credential

import (
	"context"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/credentials"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) applyClientResponse(ctx context.Context, c *Credential, cgr credentials.GetResponse) error {
	c.ID = types.StringValue(cgr.ID)
	c.CACert = types.StringValue(cgr.Raw.Credential.Cacrt)
	c.Host = types.StringValue(cgr.Raw.Credential.Host)
	if len(cgr.Raw.Credential.Hosts) > 0 {
		c.Hosts, _ = types.ListValueFrom(ctx, types.StringType, cgr.Raw.Credential.Hosts)
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
