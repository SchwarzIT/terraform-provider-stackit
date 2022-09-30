package credentialsgroup

import (
	"context"
	"time"

	clientCredentialsGroup "github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/credentials-group"
	credentialsGroup "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/object-storage/credentials-group"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	default_retry_duration = 10 * time.Minute
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.Provider.Client()
	var data credentialsGroup.CredentialsGroup
	var list clientCredentialsGroup.CredentialsGroupResponse

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found := false
	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		var err error
		list, err = c.ObjectStorage.CredentialsGroup.List(ctx, data.ProjectID.Value)
		if err != nil {
			return helper.RetryableError(err)
		}

		for _, group := range list.CredentialsGroups {
			if group.CredentialsGroupID == data.ID.Value {
				found = true
				data.Name = types.String{Value: group.DisplayName}
				data.URN = types.String{Value: group.URN}
				return nil
			}
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read credential group", err.Error())
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("couldn't find credential group", "credential group could not be found")
		return
	}

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
