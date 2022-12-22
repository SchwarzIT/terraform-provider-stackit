package stackit

import (
	"context"

	dataArgusInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/argus/instance"
	dataArgusJob "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/argus/job"
	dataElasticSearchCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/elasticsearch/credential"
	dataElasticSearchInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/elasticsearch/instance"
	dataKubernetesCluster "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/kubernetes/cluster"
	dataKubernetesProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/kubernetes/project"
	dataLogMeCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/logme/credential"
	dataLogMeInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/logme/instance"
	dataMariaDBCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/mariadb/credential"
	dataMariaDBInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/mariadb/instance"
	dataMongoDBFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/mongodb-flex/instance"
	dataObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/bucket"
	dataObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credential"
	dataObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credentials-group"
	dataPostgresFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/postgres-flex/instance"
	dataPostgresCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/postgres/credential"
	dataPostgresDBInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/postgres/instance"
	dataProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/project"
	dataRabbitMQInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/rabbitmq/instance"
	dataRedisCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/redis/credential"
	dataRedisInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/redis/instance"

	resourceArgusInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/instance"
	resourceArgusJob "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/job"
	resourceElasticSearchCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/elasticsearch/credential"
	resourceElasticsearchInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/elasticsearch/instance"
	resourceKubernetesCluster "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes/cluster"
	resourceKubernetesProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes/project"
	resourceLogMeCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/logme/credential"
	resourceLogMeInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/logme/instance"
	resourceMariaDBCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/mariadb/credential"
	resourceMariaDBInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/mariadb/instance"
	resourceMongoDBFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/mongodb-flex/instance"
	resourceObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/bucket"
	resourceObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credential"
	resourceObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credentials-group"
	resourceObjectStorageProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/project"
	resourcePostgresFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/postgres-flex/instance"
	resourcePostgresCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/postgres/credential"
	resourcePostgresDBInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/postgres/instance"
	resourceProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/project"
	resourceRabbitMQInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/rabbitmq/instance"
	resourceRedisCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/redis/credential"
	resourceRedisInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/redis/instance"

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

func (p *StackitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "stackit"
	resp.Version = p.version
}

// GetResources - Defines provider resources
func (p *StackitProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resourceArgusInstance.New,
		resourceArgusJob.New,
		resourceElasticSearchCredential.New,
		resourceElasticsearchInstance.New,
		resourceKubernetesCluster.New,
		resourceKubernetesProject.New,
		resourceLogMeCredential.New,
		resourceLogMeInstance.New,
		resourceMariaDBCredential.New,
		resourceMariaDBInstance.New,
		resourceMongoDBFlexInstance.New,
		resourceObjectStorageBucket.New,
		resourceObjectStorageCredential.New,
		resourceObjectStorageCredentialsGroup.New,
		resourceObjectStorageProject.New,
		resourcePostgresCredential.New,
		resourcePostgresDBInstance.New,
		resourceProject.New,
		resourcePostgresFlexInstance.New,
		resourceRabbitMQInstance.New,
		resourceRedisCredential.New,
		resourceRedisInstance.New,
	}
}

// GetDataSources - Defines provider data sources
func (p *StackitProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		dataArgusInstance.New,
		dataArgusJob.New,
		dataElasticSearchCredential.New,
		dataElasticSearchInstance.New,
		dataKubernetesCluster.New,
		dataKubernetesProject.New,
		dataLogMeCredential.New,
		dataLogMeInstance.New,
		dataMariaDBCredential.New,
		dataMariaDBInstance.New,
		dataMongoDBFlexInstance.New,
		dataObjectStorageBucket.New,
		dataObjectStorageCredential.New,
		dataObjectStorageCredentialsGroup.New,
		dataPostgresCredential.New,
		dataPostgresDBInstance.New,
		dataPostgresFlexInstance.New,
		dataProject.New,
		dataRabbitMQInstance.New,
		dataRedisCredential.New,
		dataRedisInstance.New,
	}
}
