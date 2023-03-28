package stackit

import (
	"context"
	"errors"
	"fmt"
	"os"

	client "github.com/SchwarzIT/community-stackit-go-client"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/clients"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/env"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var warned = false

func (p *StackitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerSchema
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Token flow
	if config.ServiceAccountEmail.IsUnknown() || config.ServiceAccountEmail.IsNull() {
		config.ServiceAccountEmail = types.StringValue(os.Getenv(ServiceAccountEmail))
	}
	if config.ServiceAccountToken.IsUnknown() || config.ServiceAccountToken.IsNull() {
		config.ServiceAccountToken = types.StringValue(os.Getenv(ServiceAccountToken))
	}

	// Key flow
	if config.ServiceAccountKey.IsUnknown() || config.ServiceAccountKey.IsNull() {
		config.ServiceAccountKey = types.StringValue(os.Getenv(ServiceAccountKey))
	}
	if config.ServiceAccountKeyPath.IsUnknown() || config.ServiceAccountKeyPath.IsNull() {
		config.ServiceAccountKeyPath = types.StringValue(os.Getenv(ServiceAccountKeyPath))
	}
	if config.PrivateKey.IsUnknown() || config.PrivateKey.IsNull() {
		config.PrivateKey = types.StringValue(os.Getenv(PrivateKey))
	}
	if config.PrivateKeyPath.IsUnknown() || config.PrivateKeyPath.IsNull() {
		config.PrivateKeyPath = types.StringValue(os.Getenv(PrivateKeyPath))
	}

	kfcl, err := keyFlow(ctx, config)
	if err == nil {
		resp.DataSourceData = kfcl
		resp.ResourceData = kfcl
		return
	}

	tfcl, err2 := tokenFlow(ctx, config)
	if err2 == nil {
		if !warned {
			resp.Diagnostics.AddWarning("Using token flow authentication", "It is recommended to use `Key flow` authentication, as it is more secure. Please consider modifying your setup.")
			warned = true
		}
		resp.DataSourceData = tfcl
		resp.ResourceData = tfcl
		return
	}

	resp.Diagnostics.AddError("couldn't initialize client with an authentication flow", fmt.Sprintf("key flow client auth:\n%s\n\ntoken flow client auth:\n%s", err.Error(), err2.Error()))
}

func keyFlow(ctx context.Context, config providerSchema) (*services.Services, error) {
	if config.ServiceAccountKey.ValueString() != "" &&
		config.PrivateKey.ValueString() != "" {
		return client.NewClientWithKeyAuth(ctx, clients.KeyFlowConfig{
			ServiceAccountKey: []byte(config.ServiceAccountKey.ValueString()),
			PrivateKey:        []byte(config.PrivateKey.ValueString()),
			Environment:       env.Environment(config.Environment.ValueString()),
		})
	}
	if config.ServiceAccountKeyPath.ValueString() != "" &&
		config.PrivateKeyPath.ValueString() != "" {
		return client.NewClientWithKeyAuth(ctx, clients.KeyFlowConfig{
			ServiceAccountKeyPath: config.ServiceAccountKeyPath.ValueString(),
			PrivateKeyPath:        config.PrivateKeyPath.ValueString(),
			Environment:           env.Environment(config.Environment.ValueString()),
		})
	}
	return nil, errors.New("no proper settings found for key flow")
}

func tokenFlow(ctx context.Context, config providerSchema) (*services.Services, error) {
	if config.ServiceAccountEmail.ValueString() != "" &&
		config.ServiceAccountToken.ValueString() != "" {
		return client.NewClientWithTokenAuth(ctx, clients.TokenFlowConfig{
			ServiceAccountEmail: config.ServiceAccountEmail.ValueString(),
			ServiceAccountToken: config.ServiceAccountToken.ValueString(),
			Environment:         env.Environment(config.Environment.ValueString()),
		})
	}
	return nil, errors.New("no proper settings found for token flow")
}
