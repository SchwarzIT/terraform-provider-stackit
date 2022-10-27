package job

import (
	"context"

	client "github.com/SchwarzIT/community-stackit-go-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// New returns a new configured resource
func New() resource.Resource {
	return &Resource{}
}

// Resource is the exported resource
type Resource struct {
	client *client.Client
}

var _ = resource.Resource(&Resource{})

// Metadata returns data resource metadata
func (r Resource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = "stackit_argus_job"
}
