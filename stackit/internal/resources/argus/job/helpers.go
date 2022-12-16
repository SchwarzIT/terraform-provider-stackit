package job

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/jobs"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	default_metrics_path                = "/metrics"
	default_scheme                      = "https"
	default_scrape_interval             = "5m"
	default_scrape_timeout              = "2m"
	default_saml2_enable_url_parameters = true
)

func (j *Job) setDefaults(job *jobs.Job) {
	if j.MetricsPath.Null || j.MetricsPath.Unknown {
		job.MetricsPath = default_metrics_path
	}
	if j.Scheme.Null || j.Scheme.Unknown {
		job.Scheme = default_scheme
	}
	if j.ScrapeInterval.Null || j.ScrapeInterval.Unknown {
		job.ScrapeInterval = default_scrape_interval
	}
	if j.ScrapeTimeout.Null || j.ScrapeTimeout.Unknown {
		job.ScrapeTimeout = default_scrape_timeout
	}
}

func (j *Job) ToClientJob() jobs.Job {
	job := jobs.Job{
		JobName:        j.Name.ValueString(),
		Scheme:         j.Scheme.ValueString(),
		MetricsPath:    j.MetricsPath.ValueString(),
		ScrapeInterval: j.ScrapeInterval.ValueString(),
		ScrapeTimeout:  j.ScrapeTimeout.ValueString(),
	}

	j.setDefaults(&job)

	if j.SAML2 != nil && !j.SAML2.EnableURLParameters.ValueBool() {
		if job.Params == nil {
			job.Params = map[string]interface{}{}
		}
		job.Params["saml2"] = []string{"disabled"}
	}

	if j.BasicAuth != nil {
		if job.BasicAuth == nil {
			job.BasicAuth = &jobs.BasicAuth{
				Username: j.BasicAuth.Username.ValueString(),
				Password: j.BasicAuth.Password.ValueString(),
			}
		}
	}

	t := make([]jobs.StaticConfig, len(j.Targets))
	for i, target := range j.Targets {
		ti := jobs.StaticConfig{}
		ti.Targets = make([]string, len(target.URLs))
		for k, v := range target.URLs {
			ti.Targets[k] = v.ValueString()
		}

		ti.Labels = make(map[string]string, len(target.Labels.Elems))
		for k, v := range target.Labels.Elems {
			ti.Labels[k], _ = common.ToString(context.TODO(), v)
		}
		t[i] = ti
	}
	job.StaticConfigs = t
	return job
}

func (j *Job) FromClientJob(cj jobs.Job) {
	j.ID = types.StringValue(cj.JobName)
	j.Name = types.StringValue(cj.JobName)
	j.MetricsPath = types.StringValue(cj.MetricsPath)
	j.Scheme = types.StringValue(cj.Scheme)
	j.ScrapeInterval = types.StringValue(cj.ScrapeInterval)
	j.ScrapeTimeout = types.StringValue(cj.ScrapeTimeout)
	j.handleSAML2(cj)
	j.handleBasicAuth(cj)
	j.handleTargets(cj)
}

func (j *Job) handleBasicAuth(cj jobs.Job) {
	if cj.BasicAuth == nil {
		j.BasicAuth = nil
		return
	}
	j.BasicAuth = &BasicAuth{
		Username: types.StringValue(cj.BasicAuth.Username),
		Password: types.StringValue(cj.BasicAuth.Password),
	}
}

func (j *Job) handleSAML2(cj jobs.Job) {
	if cj.Params == nil && j.SAML2 == nil {
		return
	}

	if j.SAML2 == nil {
		j.SAML2 = &SAML2{}
	}

	flag := true
	if v, ok := cj.Params["saml2"]; ok {
		if sl, ok := v.([]string); ok {
			if len(sl) == 1 && sl[0] == "disabled" {
				flag = false
			}
		}
	}

	j.SAML2 = &SAML2{
		EnableURLParameters: types.Bool{Value: flag},
	}
}

func (j *Job) handleTargets(cj jobs.Job) {
	newTargets := []Target{}
	for i, sc := range cj.StaticConfigs {
		nt := Target{
			URLs: []types.String{},
		}
		for _, v := range sc.Targets {
			nt.URLs = append(nt.URLs, types.StringValue(v))
		}

		if len(j.Targets) > i && j.Targets[i].Labels.IsNull() {
			nt.Labels = j.Targets[i].Labels
		} else {
			nt.Labels = types.Map{ElemType: types.StringType}
			if len(sc.Labels) > 0 {
				nt.Labels.Elems = make(map[string]attr.Value, len(sc.Labels))
			}
			for k, v := range sc.Labels {
				nt.Labels.Elems[k] = types.StringValue(v)
			}
		}
		newTargets = append(newTargets, nt)
	}
	j.Targets = newTargets
}
