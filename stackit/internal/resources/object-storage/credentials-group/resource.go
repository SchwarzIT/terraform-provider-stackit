package credentialsgroup

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// New returns a new configured resource
func New(p common.Provider) func() resource.Resource {
	return func() resource.Resource {
		return &Resource{
			Provider: p,
		}
	}
}

// Resource is the exported resource
type Resource struct {
	Provider common.Provider
}

var _ = resource.Resource(&Resource{})

// Metadata returns data resource metadata
func (r Resource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = "stackit_object_storage_credentials_group"
}
