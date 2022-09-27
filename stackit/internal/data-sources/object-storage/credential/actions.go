package credential

import (
	"context"
	"time"

	keys "github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/object-storage/access-keys"
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
	var config Credential
	var list keys.AccessKeyListResponse

	diags := req.Config.Get(ctx, &config)

	if config.DisplayName.Value == "" && config.ID.Value == "" {
		diags.AddError("missing configuration", "either display_name or id must be provided")
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		var err error
		list, err = c.ObjectStorage.AccessKeys.List(ctx, config.ProjectID.Value, "")
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read credential", err.Error())
		return
	}

	found := false
	for _, k := range list.AccessKeys {
		if k.KeyID != config.ID.Value && k.DisplayName != config.DisplayName.Value {
			continue
		}

		found = true

		config.ID = types.String{Value: k.KeyID}
		config.DisplayName = types.String{Value: k.DisplayName}
		config.Expiry = types.String{Value: k.Expires}
		diags = resp.State.Set(ctx, &config)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		break
	}

	if !found {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find credential", "credential could not be found")
		resp.Diagnostics.Append(diags...)
		return
	}
}
