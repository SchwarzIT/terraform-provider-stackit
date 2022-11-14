package stackit

import (
	"context"

	dataArgusInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/argus/instance"
	dataArgusJob "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/argus/job"
	dataElasticSearchInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/elasticsearch/instance"
	dataKubernetes "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/kubernetes"
	dataObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/bucket"
	dataObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credential"
	dataObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credentials-group"
	dataPostgresFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/postgres-flex/instance"
	dataProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/project"
	resourceArgusInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/instance"
	resourceArgusJob "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/job"
	resourceElasticsearchInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/elasticsearch/instance"
	resourceKubernetes "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes"
	resourceObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/bucket"
	resourceObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credential"
	resourceObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credentials-group"
	resourcePostgresFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/postgres-flex/instance"
	resourceProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/project"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// New returns a new STACKIT provider function
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &StackitProvider{
			version: version,
		}
	}
}

type StackitProvider struct {
	version string
}

var _ = provider.Provider(&StackitProvider{})

// Provider schema struct
type providerSchema struct {
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
	ServiceAccountToken types.String `tfsdk:"service_account_token"`
}

// GetSchema returns the provider's schema
func (p *StackitProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `
This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

~> **Note:** The provider is built using Terraform's plugin framework, therefore we recommend using Terraform CLI v1.x which supports Protocol v6
		`,
		Attributes: map[string]tfsdk.Attribute{
			"service_account_email": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Service Account Email.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_EMAIL` environment variable instead.",
			},
			"service_account_token": {
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Service Account Token.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_TOKEN` environment variable instead.",
			},
		},
	}, nil
}

func (p *StackitProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "stackit"
	resp.Version = p.version
}

// GetResources - Defines provider resources
func (p *StackitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resourceArgusInstance.New,
		resourceArgusJob.New,
		resourceElasticsearchInstance.New,
		resourceKubernetes.New,
		resourceObjectStorageBucket.New,
		resourceObjectStorageCredential.New,
		resourceObjectStorageCredentialsGroup.New,
		resourceProject.New,
		resourcePostgresFlexInstance.New,
	}
}

// GetDataSources - Defines provider data sources
func (p *StackitProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		dataArgusInstance.New,
		dataArgusJob.New,
		dataElasticSearchInstance.New,
		dataKubernetes.New,
		dataObjectStorageBucket.New,
		dataObjectStorageCredential.New,
		dataObjectStorageCredentialsGroup.New,
		dataPostgresFlexInstance.New,
		dataProject.New,
	}
}
