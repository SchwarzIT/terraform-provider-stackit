package credential

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0/credentials"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r Resource) applyClientResponse(ctx context.Context, c *Credential, cgr *credentials.CredentialsResponse) error {
	c.ID = types.StringValue(cgr.ID)
	c.Host = types.StringValue(cgr.Raw.Credentials.Host)

	c.Hosts = types.ListNull(types.StringType)
	if cgr.Raw.Credentials.Hosts != nil && len(*cgr.Raw.Credentials.Hosts) > 0 {
		h := []attr.Value{}
		for _, v := range *cgr.Raw.Credentials.Hosts {
			h = append(h, types.StringValue(v))
		}
		c.Hosts = types.ListValueMust(types.StringType, h)
	}
	c.DatabaseName = types.StringValue("")
	if cgr.Raw.Credentials.Name != nil {
		c.DatabaseName = types.StringValue(*cgr.Raw.Credentials.Name)
	}
	c.Port = types.Int64Value(0)
	if cgr.Raw.Credentials.Port != nil {
		c.Port = types.Int64Value(int64(*cgr.Raw.Credentials.Port))
	}
	c.Username = types.StringValue(cgr.Raw.Credentials.Username)
	c.Password = types.StringValue(cgr.Raw.Credentials.Password)
	c.SyslogDrainURL = types.StringValue(cgr.Raw.SyslogDrainUrl)
	c.RouteServiceURL = types.StringValue(cgr.Raw.RouteServiceUrl)
	c.URI = types.StringValue(cgr.Uri)
	return nil
}
