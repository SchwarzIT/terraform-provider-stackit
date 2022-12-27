package credential

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/credentials"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r Resource) applyClientResponse(ctx context.Context, c *Credential, cgr credentials.GetResponse) error {
	c.ID = types.StringValue(cgr.ID)
	c.Host = types.StringValue(cgr.Raw.Credential.Host)

	c.Hosts = types.ListNull(types.StringType)
	if len(cgr.Raw.Credential.Hosts) > 0 {
		h := []attr.Value{}
		for _, v := range cgr.Raw.Credential.Hosts {
			h = append(h, types.StringValue(v))
		}
		c.Hosts = types.ListValueMust(types.StringType, h)
	}
	c.Username = types.StringValue(cgr.Raw.Credential.Username)
	c.Password = types.StringValue(cgr.Raw.Credential.Password)
	c.Port = types.Int64Value(int64(cgr.Raw.Credential.Port))
	c.SyslogDrainURL = types.StringValue(cgr.Raw.SyslogDrainURL)
	c.RouteServiceURL = types.StringValue(cgr.Raw.RouteServiceURL)
	c.URI = types.StringValue(cgr.URI)
	return nil
}
