package instance

import (
	"context"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// New returns a new configured data source
func New(p common.Provider) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &DataSource{
			Provider: p,
		}
	}
}

// DataSource is the exported data source
type DataSource struct {
	Provider common.Provider
}

var _ = datasource.DataSource(&DataSource{})

// Metadata returns data resource metadata
func (r DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = "stackit_argus_instance"
}
