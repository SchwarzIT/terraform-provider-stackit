package stackit

import (
	"context"

	dataArgusInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/argus/instance"
	dataArgusJob "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/argus/job"
	dataDataServicesCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/data-services/credential"
	dataDataServicesInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/data-services/instance"
	dataKubernetesCluster "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/kubernetes/cluster"
	dataKubernetesProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/kubernetes/project"
	dataMongoDBFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/mongodb-flex/instance"
	dataMongoDBFlexUser "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/mongodb-flex/user"
	dataObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/bucket"
	dataObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credential"
	dataObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/credentials-group"
	dataObjectStorageProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/object-storage/project"
	dataPostgresFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/postgres-flex/instance"
	dataPostgresFlexUser "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/postgres-flex/user"
	dataProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/data-sources/project"

	resourceArgusCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/credential"
	resourceArgusInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/instance"
	resourceArgusJob "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/job"
	resourceDataServicesCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/data-services/credential"
	resourceDataServicesInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/data-services/instance"
	resourceKubernetesCluster "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes/cluster"
	resourceKubernetesProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes/project"
	resourceMongoDBFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/mongodb-flex/instance"
	resourceMongoDBFlexUser "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/mongodb-flex/user"
	resourceObjectStorageBucket "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/bucket"
	resourceObjectStorageCredential "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credential"
	resourceObjectStorageCredentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credentials-group"
	resourceObjectStorageProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/project"
	resourcePostgresFlexInstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/postgres-flex/instance"
	resourcePostgresFlexUser "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/postgres-flex/user"
	resourceProject "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/project"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// Token Flow
	ServiceAccountEmail = "STACKIT_SERVICE_ACCOUNT_EMAIL"
	ServiceAccountToken = "STACKIT_SERVICE_ACCOUNT_TOKEN"

	// Key Flow optional env variable (1)
	ServiceAccountKey = "STACKIT_SERVICE_ACCOUNT_KEY"
	PrivateKey        = "STACKIT_PRIVATE_KEY"

	// Key Flow optional env variable (2) using file paths
	ServiceAccountKeyPath = "STACKIT_SERVICE_ACCOUNT_KEY_PATH"
	PrivateKeyPath        = "STACKIT_PRIVATE_KEY_PATH"
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
	// Token Flow
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
	ServiceAccountToken types.String `tfsdk:"service_account_token"`

	// Key Flow optional env variable (1)
	ServiceAccountKey types.String `tfsdk:"service_account_key"`
	PrivateKey        types.String `tfsdk:"private_key"`

	// Key Flow optional env variable (2) using file paths
	ServiceAccountKeyPath types.String `tfsdk:"service_account_key_path"`
	PrivateKeyPath        types.String `tfsdk:"private_key_path"`

	// Key Flow

	Environment types.String `tfsdk:"environment"`
}

// Schema returns the provider's schema
func (p *StackitProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: md,
		Attributes: map[string]schema.Attribute{
			"service_account_email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Service Account Email.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_EMAIL` environment variable instead.",
			},
			"service_account_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Service Account Token.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_TOKEN` environment variable instead.",
			},
			"service_account_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Service Account Key.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_KEY` environment variable instead.",
			},
			"private_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Service Account Key.<br />This attribute can also be loaded from `STACKIT_PRIVATE_KEY` environment variable instead.",
			},
			"service_account_key_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Service Account Key.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_KEY_PATH` environment variable instead.",
			},
			"private_key_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Service Account Key.<br />This attribute can also be loaded from `STACKIT_PRIVATE_KEY_PATH` environment variable instead.",
			},
			"environment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The API environment that the provider interacts with. Options are `dev`, `qa`, `prod`.<br />This attribute can also be loaded from `STACKIT_ENV` environment variable instead.",
				Validators: []validator.String{
					stringvalidator.OneOf("dev", "qa", "prod"),
				},
			},
		},
	}
}

func (p *StackitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "stackit"
	resp.Version = p.version
}

// GetResources - Defines provider resources
func (p *StackitProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resourceArgusCredential.New,
		resourceArgusInstance.New,
		resourceArgusJob.New,
		resourceDataServicesCredential.NewElasticSearch,
		resourceDataServicesCredential.NewLogMe,
		resourceDataServicesCredential.NewMariaDB,
		resourceDataServicesCredential.NewPostgres,
		resourceDataServicesCredential.NewRabbitMQ,
		resourceDataServicesCredential.NewRedis,
		resourceDataServicesInstance.NewElasticSearch,
		resourceDataServicesInstance.NewLogMe,
		resourceDataServicesInstance.NewMariaDB,
		resourceDataServicesInstance.NewPostgres,
		resourceDataServicesInstance.NewRabbitMQ,
		resourceDataServicesInstance.NewRedis,
		resourceKubernetesCluster.New,
		resourceKubernetesProject.New,
		resourceMongoDBFlexInstance.New,
		resourceMongoDBFlexUser.New,
		resourceObjectStorageBucket.New,
		resourceObjectStorageCredential.New,
		resourceObjectStorageCredentialsGroup.New,
		resourceObjectStorageProject.New,
		resourcePostgresFlexInstance.New,
		resourcePostgresFlexUser.New,
		resourceProject.New,
	}
}

// GetDataSources - Defines provider data sources
func (p *StackitProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		dataArgusInstance.New,
		dataArgusJob.New,
		dataDataServicesCredential.NewElasticSearch,
		dataDataServicesCredential.NewLogMe,
		dataDataServicesCredential.NewMariaDB,
		dataDataServicesCredential.NewPostgres,
		dataDataServicesCredential.NewRabbitMQ,
		dataDataServicesCredential.NewRedis,
		dataDataServicesInstance.NewElasticSearch,
		dataDataServicesInstance.NewLogMe,
		dataDataServicesInstance.NewMariaDB,
		dataDataServicesInstance.NewPostgres,
		dataDataServicesInstance.NewRabbitMQ,
		dataDataServicesInstance.NewRedis,
		dataKubernetesCluster.New,
		dataKubernetesProject.New,
		dataMongoDBFlexInstance.New,
		dataMongoDBFlexUser.New,
		dataObjectStorageBucket.New,
		dataObjectStorageCredential.New,
		dataObjectStorageCredentialsGroup.New,
		dataObjectStorageProject.New,
		dataPostgresFlexInstance.New,
		dataPostgresFlexUser.New,
		dataProject.New,
	}
}

var md = `
The STACKIT provider is a project developed and maintained by the STACKIT community within Schwarz IT. Please note that it is not an official provider endorsed or maintained by STACKIT.

~> **Note:** The provider is built using Terraform's plugin framework, therefore we recommend using Terraform CLI v1.x which supports Protocol v6

## Authentication

Before you can start using the client, you will need to create a STACKIT Service Account in your project and assign it the appropriate permissions (i.e. ` + "`project.owner`)." + `
After the service account has been created, you can authenticate to the client using the ` + "`Key flow`" + `  (recommended) or with the static ` + "`Token flow`" + ` (less secure as the token is long-lived).

### Key flow

1. Follow instructions in the [community client](https://github.com/SchwarzIT/community-stackit-go-client#key-flow)

2. You can configure the Key flow by providing paths using environment variables or by configuring the paths in the provider block (` + "`Key flow (1)`" + ` example below)

   ` + "```bash" + `
   export STACKIT_SERVICE_ACCOUNT_KEY_PATH="sa_key.json"
   export STACKIT_PRIVATE_KEY_PATH="private_key.pem"
   ` + "```" + `

3. Another way is to provide the contents directly, either with environment variables or by configuring the provider directly (` + "`Key flow (2)`" + ` example below)

   ` + "```bash" + `
   export STACKIT_SERVICE_ACCOUNT_KEY_="..."
   export STACKIT_PRIVATE_KEY="..."
   ` + "```" + `

### Token flow

1. Set the following environment variables or configure the provider directly (example below)

   ` + "```bash" + `
   export STACKIT_SERVICE_ACCOUNT_EMAIL=email
   export STACKIT_SERVICE_ACCOUNT_TOKEN=token
   ` + "```" + `

&nbsp;
`
