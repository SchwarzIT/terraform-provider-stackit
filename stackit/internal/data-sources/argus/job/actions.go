package job

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/jobs"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/job"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config job.Job

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	b, err := c.Argus.Jobs.Get(ctx, config.ProjectID.ValueString(), config.ArgusInstanceID.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to read job", err.Error())
		return
	}

	config.FromClientJob(b.Data)
	handleSAML2(&config, b.Data)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func handleSAML2(j *job.Job, cj jobs.Job) {
	flag := true
	if cj.Params != nil {
		if v, ok := cj.Params["saml2"]; ok {
			if sl, ok := v.([]string); ok {
				if len(sl) == 1 && sl[0] == "disabled" {
					flag = false
				}
			}
		}
	}

	j.SAML2 = &job.SAML2{
		EnableURLParameters: types.Bool{Value: flag},
	}
}
