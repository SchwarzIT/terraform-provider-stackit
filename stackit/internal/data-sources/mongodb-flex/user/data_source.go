package user

import (
	"context"
	"fmt"

	client "github.com/SchwarzIT/community-stackit-go-client"
	mongodbflex "github.com/SchwarzIT/community-stackit-go-client/pkg/services/mongodb-flex/v1.0/generated"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/urls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// New returns a new configured data source
func New() datasource.DataSource {
	return &DataSource{
		urls: mongodbflex.BaseURLs,
	}
}

// DataSource is the exported data source
type DataSource struct {
	client *client.Client
	urls   urls.ByEnvs
}

var _ = datasource.DataSource(&DataSource{})

// Metadata returns data resource metadata
func (d *DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = "stackit_mongodb_flex_user"
}

// Configure configures the data source client
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
