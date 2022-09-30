package job

import (
	"context"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/jobs"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/argus/job"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"net/http"
	"strings"
)

func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.Provider.Client()
	var config job.Job
	var b jobs.GetJobResponse

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		var err error
		b, err = c.Argus.Jobs.Get(ctx, config.ProjectID.Value, config.ArgusInstanceID.Value, config.Name.Value)
		if err != nil {
			if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
				return helper.NonRetryableError(err)
			}
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read instance", err.Error())
		return
	}
	config.Name = types.String{Value: b.Data.JobName}
	config.MetricsPath = types.String{Value: b.Data.MetricsPath}
	config.Scheme = types.String{Value: b.Data.Scheme}
	config.ScrapeInterval = types.String{Value: b.Data.ScrapeInterval}
	config.ScrapeTimeout = types.String{Value: b.Data.ScrapeTimeout}
	config.BasicAuth.Username = types.String{Value: b.Data.BasicAuth.Username}
	config.BasicAuth.Password = types.String{Value: b.Data.BasicAuth.Password}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
