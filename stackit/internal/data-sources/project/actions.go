package project

import (
	"context"
	"time"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/projects"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const default_retry_duration = 10 * time.Minute

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.Provider.Client()
	var project projects.Project
	var p Project

	diags := req.Config.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.RetryContext(ctx, default_retry_duration, func() *helper.RetryError {
		var err error
		project, err = c.Projects.Get(ctx, p.ID.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read project", err.Error())
		return
	}

	p.ID = types.String{Value: project.ID}
	p.Name = types.String{Value: project.Name}
	p.BillingRef = types.String{Value: project.BillingReference}

	diags = resp.State.Set(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
