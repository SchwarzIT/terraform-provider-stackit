package credential

import (
	"context"
	"fmt"

	client "github.com/SchwarzIT/community-stackit-go-client"
	dataservices "github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0/generated"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ResourceService string

const (
	ElasticSearch ResourceService = "elasticsearch"
	LogMe         ResourceService = "logme"
	MariaDB       ResourceService = "mariadb"
	Postgres      ResourceService = "postgres"
	Redis         ResourceService = "redis"
	RabbitMQ      ResourceService = "rabbitmq"
)

func (s ResourceService) Display() string {
	switch s {
	case ElasticSearch:
		return "ElasticSearch"
	case LogMe:
		return "LogMe"
	case MariaDB:
		return "MariaDB"
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
	service ResourceService
}

var _ = resource.Resource(&Resource{})

// Metadata returns data resource metadata
func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = fmt.Sprintf("stackit_%s_credential", r.service)
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
		r.client = c.ElasticSearch
	case LogMe:
		r.client = c.LogMe
	case MariaDB:
		r.client = c.MariaDB
	case Postgres:
		r.client = c.PostgresDB
	case Redis:
		r.client = c.Reddis
	case RabbitMQ:
		r.client = c.RabbitMQ
	}
}
