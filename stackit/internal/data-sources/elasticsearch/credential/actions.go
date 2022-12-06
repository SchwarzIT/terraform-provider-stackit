package credential

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Credential
	diags := req.Config.Get(ctx, &config)
	es := c.DataServices.ElasticSearch

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := es.Credentials.List(ctx, config.ProjectID.Value, config.InstanceID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to list credentials", err.Error())
		return
	}

	found := -1
	existing := ""
	for i, cred := range list.CredentialsList {
		if cred.ID == config.ID.Value {
			found = i
			break
		}
		if existing == "" {
			existing = "\navailable credentials in the project are:"
		}
		existing = fmt.Sprintf("%s\n- %s", existing, cred.ID)
	}

	if found == -1 {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find credentials", "credentials could not be found."+existing)
		resp.Diagnostics.Append(diags...)
		return
	}

	// set found instance
	cred := list.CredentialsList[found]
	config.ID = types.String{Value: cred.ID}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
