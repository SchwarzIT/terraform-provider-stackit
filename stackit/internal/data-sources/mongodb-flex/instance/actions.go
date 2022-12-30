package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/mongodb-flex/v1.0/generated/instance"
	resource "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/mongodb-flex/instance"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client.Services.MongoDBFlex
	var config resource.Instance
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.Instance.ListWithResponse(ctx, config.ProjectID.ValueString(), &instance.ListParams{})
	if err != nil {
		resp.Diagnostics.AddError("failed making list instances request", err.Error())
		return
	}
	if res.HasError != nil {
		resp.Diagnostics.AddError("list instances response has an error", res.HasError.Error())
		return
	}
	if res.JSON200 == nil || res.JSON200.Items == nil {
		resp.Diagnostics.AddError("failed to parse response", "JSON200 == nil or .Items == nil")
		return
	}

	list := *res.JSON200.Items
	found := -1
	existing := ""
	for i, instance := range list {
		if instance.Name == nil || instance.ID == nil {
			continue
		}
		if strings.EqualFold(*instance.Name, config.Name.ValueString()) {
			found = i
			break
		}
		if existing == "" {
			existing = "\navailable instances in the project are:"
		}
		existing = fmt.Sprintf("%s\n- %s", existing, *instance.Name)
	}

	if found == -1 {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find instance", "instance could not be found."+existing)
		resp.Diagnostics.Append(diags...)
		return
	}

	// set found instance
	instance := list[found]
	ires, err := c.Instance.GetWithResponse(ctx, config.ProjectID.ValueString(), *instance.ID)
	if err != nil {
		resp.Diagnostics.AddError("failed making get instance request", err.Error())
		return
	}
	if ires.HasError != nil {
		resp.Diagnostics.AddError("list instances response has an error", ires.HasError.Error())
		return
	}
	if ires.JSON200 == nil || ires.JSON200.Item == nil {
		resp.Diagnostics.AddError("failed to parse response", "JSON200 == nil or .Items == nil")
		return
	}

	resource.ApplyClientResponse(&config, ires.JSON200.Item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
