package user

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/env"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services"
	postgresflex "github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// New returns a new configured resource
func New() resource.Resource {
	return &Resource{
		urls: postgresflex.BaseURLs,
	}
}

// Resource is the exported resource
type Resource struct {
	client *services.Services
	urls   env.EnvironmentURLs
}

var _ = resource.Resource(&Resource{})

// Metadata returns data resource metadata
func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = "stackit_postgres_flex_user"
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

	r.client = c
}
