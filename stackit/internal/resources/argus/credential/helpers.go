package credential

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/generated/instances"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r Resource) applyClientResponse(ctx context.Context, c *Credential, cgr instances.Credentials) error {
	c.ID = types.StringValue(cgr.Username)
	c.Username = types.StringValue(cgr.Username)
	c.Password = types.StringValue(cgr.Password)
	return nil
}
