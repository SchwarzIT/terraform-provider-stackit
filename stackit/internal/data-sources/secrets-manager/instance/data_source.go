package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/baseurl"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services"
	secretsmanager "github.com/SchwarzIT/community-stackit-go-client/pkg/services/secrets-manager/v1.1.0"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// New returns a new configured data source
func New() datasource.DataSource {
	return &DataSource{
		urls: secretsmanager.BaseURLs,
	}
}

// DataSource is the exported data source
type DataSource struct {
	client *services.Services
	urls   baseurl.BaseURL
}

var _ = datasource.DataSource(&DataSource{})

// Metadata returns data resource metadata
func (d *DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = "stackit_secrets_manager_instance"
}

// Configure configures the data source client
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*services.Services)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
