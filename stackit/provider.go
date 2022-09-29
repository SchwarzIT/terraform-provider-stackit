package stackit

import (
	"context"

	client "github.com/SchwarzIT/community-stackit-go-client"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	dataKubernetes "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/kubernetes"
	dataObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/bucket"
	dataObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credential"
	dataObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credentials-group"
	dataProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/project"
	resourceArgusInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/instance"
	resourceArgusJob "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/job"
	resourceKubernetes "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes"
	resourceObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/bucket"
	resourceObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credential"
	resourceObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credentials-group"
	resourceProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/project"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// New returns a new STACKIT provider
func New() provider.Provider {
	return &StackitProvider{}
}

type StackitProvider struct {
	configured       bool
	client           *client.Client
	serviceAccountID string
}

var _ = provider.Provider(&StackitProvider{})
var _ = common.Provider(&StackitProvider{})

// GetSchema returns the provider's schema
func (p *StackitProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `
This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

~> **Note:** The provider is built using Terraform's plugin framework, therefore we recommend using Terraform CLI v1.x which supports Protocol v6
		`,
		Attributes: map[string]tfsdk.Attribute{
			"service_account_id": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Service Account ID.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_ID` environment variable instead.",
			},
			"service_account_token": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Service Account Token.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_TOKEN` environment variable instead.",
			},
			"customer_account_id": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Customer Account ID (Organization ID).<br />This attribute can also be loaded from `STACKIT_CUSTOMER_ACCOUNT_ID` environment variable instead.",
			},
		},
	}, nil
}

// GetResources - Defines provider resources
func (p *StackitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resourceArgusJob.New(p),
		resourceArgusInstance.New(p),
		resourceKubernetes.New(p),
		resourceObjectStorageBucket.New(p),
		resourceObjectStorageCredential.New(p),
		resourceObjectStorageCredentialsGroup.New(p),
		resourceProject.New(p),
	}
}

// GetDataSources - Defines provider data sources
func (p *StackitProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		dataObjectStorageBucket.New(p),
		dataObjectStorageCredential.New(p),
		dataObjectStorageCredentialsGroup.New(p),
		dataKubernetes.New(p),
		dataProject.New(p),
	}
}

// IsConfigured - returns true when the provider has been configured
func (p *StackitProvider) IsConfigured() bool {
	return p.configured
}

// Client - returns the STACKIT client
func (p *StackitProvider) Client() *client.Client {
	return p.client
}

// ServiceAccountID - returns the service account id
func (p *StackitProvider) ServiceAccountID() string {
	return p.serviceAccountID
}
