package stackit

import (
	"context"
	"os"
	"time"

	client "github.com/SchwarzIT/community-stackit-go-client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (p *StackitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerSchema
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var email string
	if config.ServiceAccountEmail.IsUnknown() || config.ServiceAccountEmail.IsNull() {
		email = os.Getenv("STACKIT_SERVICE_ACCOUNT_EMAIL")
		config.ServiceAccountEmail = types.StringValue(email)
	} else {
		email = config.ServiceAccountEmail.ValueString()
	}

	if email == "" {
		resp.Diagnostics.AddError("missing mandatory field", "STACKIT service account email must be provided")
		return
	}

	var token string
	if config.ServiceAccountToken.IsUnknown() || config.ServiceAccountToken.IsNull() {
		token = os.Getenv("STACKIT_SERVICE_ACCOUNT_TOKEN")
		config.ServiceAccountToken = types.StringValue(token)
	} else {
		token = config.ServiceAccountToken.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError("missing mandatory field", "STACKIT service account token must be provided")
		return
	}

	c, err := client.New(context.Background(), client.Config{
		ServiceAccountEmail: email,
		ServiceAccountToken: token,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create STACKIT client:\n"+err.Error(),
		)
		return
	}

	httpClient := c.GetHTTPClient()
	httpClient.Timeout = time.Minute

	resp.DataSourceData = c
	resp.ResourceData = c
}
