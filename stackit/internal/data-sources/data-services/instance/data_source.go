package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/env"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services"
	dataservices "github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type DataSourceService string

const (
	ElasticSearch DataSourceService = "elasticsearch"
	LogMe         DataSourceService = "logme"
	MariaDB       DataSourceService = "mariadb"
	Opensearch    DataSourceService = "opensearch"
	Postgres      DataSourceService = "postgres"
	Redis         DataSourceService = "redis"
	RabbitMQ      DataSourceService = "rabbitmq"
)

func (s DataSourceService) Display() string {
	switch s {
	case ElasticSearch:
		return "ElasticSearch"
	case LogMe:
		return "LogMe"
	case MariaDB:
		return "MariaDB"
	case Opensearch:
		return "Opensearch"
	case Postgres:
		return "Postgres"
	case Redis:
		return "Redis"
	case RabbitMQ:
		return "RabbitMQ"
	}
	return ""
}

// NewElasticSearch returns a new configured resource
func NewElasticSearch() datasource.DataSource {
	return &DataSource{
		service: ElasticSearch,
		urls:    dataservices.GetBaseURLs(dataservices.ElasticSearch),
	}
}

// NewLogMe returns a new configured resource
func NewLogMe() datasource.DataSource {
	return &DataSource{
		service: LogMe,
		urls:    dataservices.GetBaseURLs(dataservices.LogMe),
	}
}

// NewMariaDB returns a new configured resource
func NewMariaDB() datasource.DataSource {
	return &DataSource{
		service: MariaDB,
		urls:    dataservices.GetBaseURLs(dataservices.MariaDB),
	}
}

// NewOpensearch returns a new configured resource
func NewOpensearch() datasource.DataSource {
	return &DataSource{
		service: Opensearch,
		urls:    dataservices.GetBaseURLs(dataservices.Opensearch),
	}
}

// NewPostgres returns a new configured resource
func NewPostgres() datasource.DataSource {
	return &DataSource{
		service: Postgres,
		urls:    dataservices.GetBaseURLs(dataservices.PostgresDB),
	}
}

// NewRedis returns a new configured resource
func NewRedis() datasource.DataSource {
	return &DataSource{
		service: Redis,
		urls:    dataservices.GetBaseURLs(dataservices.Redis),
	}
}

// NewRabbitMQ returns a new configured resource
func NewRabbitMQ() datasource.DataSource {
	return &DataSource{
		service: RabbitMQ,
		urls:    dataservices.GetBaseURLs(dataservices.RabbitMQ),
	}
}

// DataSource is the exported data source
type DataSource struct {
	client  *dataservices.ClientWithResponses
	service DataSourceService
	urls    env.EnvironmentURLs
}

var _ = datasource.DataSource(&DataSource{})

// Metadata returns data resource metadata
func (d *DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = fmt.Sprintf("stackit_%s_instance", d.service)
}

// Configure configures the data source client
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*services.Services)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	switch d.service {
	case ElasticSearch:
		d.client = c.ElasticSearch
	case LogMe:
		d.client = c.LogMe
	case MariaDB:
		d.client = c.MariaDB
	case Opensearch:
		d.client = c.Opensearch
	case Postgres:
		d.client = c.PostgresDB
	case Redis:
		d.client = c.Redis
	case RabbitMQ:
		d.client = c.RabbitMQ
	}
}
