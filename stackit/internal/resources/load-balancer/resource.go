package loadbalancer

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/baseurl"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services"
	loadbalancer "github.com/SchwarzIT/community-stackit-go-client/pkg/services/load-balancer/1.3.0"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// New returns a new configured resource
func New() resource.Resource {
	return &Resource{
		urls: loadbalancer.BaseURLs,
	}
}

// Resource is the exported resource
type Resource struct {
	client *services.Services
	urls   baseurl.BaseURL
}

var _ = resource.Resource(&Resource{})

// Metadata returns data resource metadata
func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = "stackit_load_balancer"
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
