package job

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/argus/jobs"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
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
		JobName:        j.Name.Value,
		Scheme:         j.Scheme.Value,
		MetricsPath:    j.MetricsPath.Value,
		ScrapeInterval: j.ScrapeInterval.Value,
		ScrapeTimeout:  j.ScrapeTimeout.Value,
	}

	j.setDefaults(&job)

	if j.SAML2 != nil && !j.SAML2.EnableURLParameters.Value {
		if job.Params == nil {
			job.Params = make(map[string]interface{}, 1)
		}
		job.Params["saml2"] = []string{"disabled"}
	}

	if j.BasicAuth != nil {
		if job.BasicAuth == nil {
			job.BasicAuth = &jobs.BasicAuth{
				Username: j.BasicAuth.Username.Value,
				Password: j.BasicAuth.Password.Value,
			}
		}
	}

	t := make([]jobs.StaticConfig, len(j.Targets))
	for i, target := range j.Targets {
		ti := jobs.StaticConfig{}
		ti.Targets = make([]string, len(target.URLs))
		for k, v := range target.URLs {
			ti.Targets[k] = v.Value
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
	j.ID = types.String{Value: cj.JobName}
	j.Name = types.String{Value: cj.JobName}
	j.MetricsPath = types.String{Value: cj.MetricsPath}
	j.Scheme = types.String{Value: cj.Scheme}
	j.ScrapeInterval = types.String{Value: cj.ScrapeInterval}
	j.ScrapeTimeout = types.String{Value: cj.ScrapeTimeout}
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
		Username: types.String{Value: cj.BasicAuth.Username},
		Password: types.String{Value: cj.BasicAuth.Password},
	}
}

func (j *Job) handleSAML2(cj jobs.Job) {
	if cj.Params == nil {
		j.SAML2 = nil
		return
	}

	if j.SAML2 == nil {
		j.SAML2 = &SAML2{}
	}

	flag := true
	v, ok1 := cj.Params["saml2"]
	if sl, ok2 := v.([]string); ok1 && ok2 {
		if len(sl) == 1 && sl[0] == "disabled" {
			flag = false
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
			nt.URLs = append(nt.URLs, types.String{Value: v})
		}

		if len(j.Targets) > i && j.Targets[i].Labels.IsNull() {
			nt.Labels = j.Targets[i].Labels
		} else {
			nt.Labels = types.Map{ElemType: types.StringType}
			for k, v := range sc.Labels {
				nt.Labels.Elems[k] = types.String{Value: v}
			}
		}
		newTargets = append(newTargets, nt)
	}
	j.Targets = newTargets
}
