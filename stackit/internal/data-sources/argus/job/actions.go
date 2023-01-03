package job

import (
	"context"

	scrapeconfig "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/generated/scrape-config"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/job"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config job.Job
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.Argus.ScrapeConfig.ReadWithResponse(ctx, config.ProjectID.ValueString(), config.ArgusInstanceID.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed preparing read job request", err.Error())
		return
	}
	if res.HasError != nil {
		resp.Diagnostics.AddError("failed making read job request", res.HasError.Error())
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("failed parsing read job request", "JSON200 == nil")
		return
	}

	config.FromClientJob(res.JSON200.Data)
	handleSAML2(&config, res.JSON200.Data)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func handleSAML2(j *job.Job, cj scrapeconfig.Job) {
	flag := true
	if cj.Params != nil {
		param := *cj.Params
		if v, ok := param["saml2"]; ok {
			if len(v) == 1 && v[0] == "disabled" {
				flag = false
			}
		}
	}

	j.SAML2 = &job.SAML2{
		EnableURLParameters: types.BoolValue(flag),
	}
}
