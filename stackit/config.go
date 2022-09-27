package stackit

import (
	"context"
	"os"

	client "github.com/SchwarzIT/community-stackit-go-client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Provider schema struct
type providerSchema struct {
	ServiceAccountID    types.String `tfsdk:"service_account_id"`
	ServiceAccountToken types.String `tfsdk:"service_account_token"`
	CustomerAccountID   types.String `tfsdk:"customer_account_id"`
}

func (p *StackitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerSchema
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id string
	if config.ServiceAccountID.IsUnknown() || config.ServiceAccountID.IsNull() {
		id = os.Getenv("STACKIT_SERVICE_ACCOUNT_ID")
		config.ServiceAccountID = types.String{Value: id}
	} else {
		id = config.ServiceAccountID.Value
	}

	if id == "" {
		resp.Diagnostics.AddError("missing mandatory field", "STACKIT service account ID must be provided")
		return
	}

	var token string
	if config.ServiceAccountToken.IsUnknown() || config.ServiceAccountToken.IsNull() {
		token = os.Getenv("STACKIT_SERVICE_ACCOUNT_TOKEN")
		config.ServiceAccountToken = types.String{Value: token}
	} else {
		token = config.ServiceAccountToken.Value
	}

	if token == "" {
		resp.Diagnostics.AddError("missing mandatory field", "STACKIT service account token must be provided")
		return
	}

	var ca string
	if config.CustomerAccountID.IsUnknown() || config.CustomerAccountID.IsNull() {
		ca = os.Getenv("STACKIT_CUSTOMER_ACCOUNT_ID")
		config.CustomerAccountID = types.String{Value: ca}
	} else {
		ca = config.CustomerAccountID.Value
	}

	if ca == "" {
		resp.Diagnostics.AddError("missing mandatory field", "STACKIT customer account ID must be provided")
		return
	}

	c, err := client.New(context.Background(), &client.Config{
		ServiceAccountID: id,
		Token:            token,
		OrganizationID:   ca,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create STACKIT client:\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
	p.serviceAccountID = id
}
