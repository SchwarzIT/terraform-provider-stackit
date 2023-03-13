package job

import (
	"context"

	scrapeconfig "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/generated/scrape-config"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
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
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		resp.Diagnostics.AddError("failed to read job", agg.Error())
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
