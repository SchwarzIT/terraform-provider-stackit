package instance

import (
	"context"
	"fmt"

	client "github.com/SchwarzIT/community-stackit-go-client"
	dataservices "github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0/generated"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type DataSourceService string

const (
	ElasticSearch DataSourceService = "elasticsearch"
	LogMe         DataSourceService = "logme"
	MariaDB       DataSourceService = "mariadb"
	Postgres      DataSourceService = "postgres"
	Redis         DataSourceService = "redis"
	RabbitMQ      DataSourceService = "rabbitmq"
)

// NewElasticSearch returns a new configured resource
func NewElasticSearch() resource.Resource {
	return &Resource{service: ElasticSearch}
}

// NewLogMe returns a new configured resource
func NewLogMe() resource.Resource {
	return &Resource{service: LogMe}
}

// NewMariaDB returns a new configured resource
func NewMariaDB() resource.Resource {
	return &Resource{service: MariaDB}
}

// NewPostgres returns a new configured resource
func NewPostgres() resource.Resource {
	return &Resource{service: Postgres}
}

// NewRedis returns a new configured resource
func NewRedis() resource.Resource {
	return &Resource{service: Redis}
}

// NewRabbitMQ returns a new configured resource
func NewRabbitMQ() resource.Resource {
	return &Resource{service: RabbitMQ}
}

// Resource is the exported resource
type Resource struct {
	client  *dataservices.ClientWithResponses
	service DataSourceService
}

var _ = resource.Resource(&Resource{})

// Metadata returns data resource metadata
func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = fmt.Sprintf("stackit_%s_instance", r.service)
}

// Configure the resource client
func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	switch r.service {
	case ElasticSearch:
		r.client = c.Services.ElasticSearch
	case LogMe:
		r.client = c.Services.LogMe
	case MariaDB:
		r.client = c.Services.MariaDB
	case Postgres:
		r.client = c.Services.PostgresDB
	case Redis:
		r.client = c.Services.Reddis
	case RabbitMQ:
		r.client = c.Services.RabbitMQ
	}
}
