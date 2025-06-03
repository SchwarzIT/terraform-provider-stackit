package instance

import (
	"context"
	"fmt"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/baseurl"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services"
	dataservices "github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ResourceService string

const (
	ElasticSearch ResourceService = "elasticsearch"
	LogMe         ResourceService = "logme"
	MariaDB       ResourceService = "mariadb"
	Opensearch    ResourceService = "opensearch"
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
func NewElasticSearch() resource.Resource {
	return &Resource{
		service: ElasticSearch,
		urls:    dataservices.GetBaseURLs(dataservices.ElasticSearch),
	}
}

// NewLogMe returns a new configured resource
func NewLogMe() resource.Resource {
	return &Resource{
		service: LogMe,
		urls:    dataservices.GetBaseURLs(dataservices.LogMe),
	}
}

// NewMariaDB returns a new configured resource
func NewMariaDB() resource.Resource {
	return &Resource{
		service: MariaDB,
		urls:    dataservices.GetBaseURLs(dataservices.MariaDB),
	}
}

// NewPostgres returns a new configured resource
func NewPostgres() resource.Resource {
	return &Resource{
		service: Postgres,
		urls:    dataservices.GetBaseURLs(dataservices.PostgresDB),
	}
}

// NewRedis returns a new configured resource
func NewRedis() resource.Resource {
	return &Resource{
		service: Redis,
		urls:    dataservices.GetBaseURLs(dataservices.Redis),
	}
}

// NewRabbitMQ returns a new configured resource
func NewRabbitMQ() resource.Resource {
	return &Resource{
		service: RabbitMQ,
		urls:    dataservices.GetBaseURLs(dataservices.RabbitMQ),
	}
}

// NewOpensearch returns a new configured resource
func NewOpensearch() resource.Resource {
	return &Resource{
		service: Opensearch,
		urls:    dataservices.GetBaseURLs(dataservices.Opensearch),
	}
}

// Resource is the exported resource
type Resource struct {
	client  *dataservices.ClientWithResponses
	service ResourceService
	urls    baseurl.BaseURL
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

	c, ok := req.ProviderData.(*services.Services)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *services.Services, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.setClient(c)
}

// ModifyPlan only for RABBITMQ
func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// we just do things for RabbitMQ!
	if r.service != RabbitMQ {
		return
	}

	// Return early if we are deleting (plan is null) or creating (state is null)
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var planData, stateData *Instance
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateVersion, err := version.NewVersion(stateData.Version.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error parsing state version", "Error parsing state version: "+err.Error())
		return
	}

	planVersion, err := version.NewVersion(planData.Version.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error parsing plan version", "Error parsing plan version: "+err.Error())
		return
	}

	if stateVersion == nil || planVersion == nil {
		resp.Diagnostics.AddError("StateVersion or PlanVersion is nil", "StateVersion or PlanVersion is nil.")
		return
	}

	if stateVersion.Segments()[0] != planVersion.Segments()[0] {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("version"))
		resp.Diagnostics.AddAttributeWarning(path.Root("version"), "Changing Version on RabbitMQ require replacement", "Changing Version on RabbitMQ require replacement")
	}

}

func (r *Resource) setClient(c *services.Services) {
	switch r.service {
	case ElasticSearch:
		r.client = c.ElasticSearch
	case LogMe:
		r.client = c.LogMe
	case MariaDB:
		r.client = c.MariaDB
	case Opensearch:
		r.client = c.Opensearch
	case Postgres:
		r.client = c.PostgresDB
	case Redis:
		r.client = c.Redis
	case RabbitMQ:
		r.client = c.RabbitMQ
	}
}
