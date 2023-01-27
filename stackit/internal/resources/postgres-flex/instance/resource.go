package postgresinstance

import (
	"context"
	"fmt"

	client "github.com/SchwarzIT/community-stackit-go-client"
	postgresflex "github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/generated"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/urls"
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
	client *client.Client
	urls   urls.ByEnvs
}

var _ = resource.Resource(&Resource{})

// Metadata returns data resource metadata
func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = "stackit_postgres_flex_instance"
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

	r.client = c
}
