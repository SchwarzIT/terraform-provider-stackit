package job

import (
	"context"

	scrapeconfig "github.com/SchwarzIT/community-stackit-go-client/pkg/services/argus/v1.0/scrape-config"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
)

const (
	DefaultMetricsPath              = "/metrics"
	DefaultScheme                   = "https"
	DefaultScrapeInterval           = "5m"
	DefaultScrapeTimeout            = "2m"
	DefaultSampleLimit              = 5000
	DefaultSAML2EnableURLParameters = true
)

func (j *Job) setDefaults(job *scrapeconfig.CreateJSONBody) {
	if job == nil {
		return
	}
	if j.MetricsPath.IsNull() || j.MetricsPath.IsUnknown() {
		s := DefaultMetricsPath
		job.MetricsPath = &s
	}
	if j.Scheme.IsNull() || j.Scheme.IsUnknown() {
		job.Scheme = DefaultScheme
	}
	if j.ScrapeInterval.IsNull() || j.ScrapeInterval.IsUnknown() {
		job.ScrapeInterval = DefaultScrapeInterval
	}
	if j.ScrapeTimeout.IsNull() || j.ScrapeTimeout.IsUnknown() {
		job.ScrapeTimeout = DefaultScrapeTimeout
	}
	if j.SampleLimit.IsNull() || j.SampleLimit.IsUnknown() {
		job.SampleLimit = toFloat32Ptr(DefaultSampleLimit)
	}
}

func (j *Job) setDefaultsUpdate(partialUpdate scrapeconfig.PartialUpdateJSONBody) {
	if len(partialUpdate) == 0 {
		return
	}
	job := &partialUpdate[0]
	if j.MetricsPath.IsNull() || j.MetricsPath.IsUnknown() {
		s := DefaultMetricsPath
		job.MetricsPath = &s
	}
	if j.Scheme.IsNull() || j.Scheme.IsUnknown() {
		job.Scheme = DefaultScheme
	}
	if j.ScrapeInterval.IsNull() || j.ScrapeInterval.IsUnknown() {
		job.ScrapeInterval = DefaultScrapeInterval
	}
	if j.ScrapeTimeout.IsNull() || j.ScrapeTimeout.IsUnknown() {
		job.ScrapeTimeout = DefaultScrapeTimeout
	}
}

func (j *Job) ToClientJob() scrapeconfig.CreateJSONBody {
	job := scrapeconfig.CreateJSONBody{
		JobName:        j.Name.ValueString(),
		Scheme:         scrapeconfig.CreateJSONBodyScheme(j.Scheme.ValueString()),
		MetricsPath:    j.MetricsPath.ValueStringPointer(),
		ScrapeInterval: j.ScrapeInterval.ValueString(),
		ScrapeTimeout:  j.ScrapeTimeout.ValueString(),
		// This conversion might be lossy if the value is greater than 16777215.
		SampleLimit: toFloat32Ptr(float32(j.SampleLimit.ValueInt64())),
	}

	j.setDefaults(&job)

	if j.SAML2 != nil && !j.SAML2.EnableURLParameters.ValueBool() {
		if job.Params == nil {
			job.Params = &map[string]interface{}{}
		}
		p := *job.Params
		p["saml2"] = []string{"disabled"}
		job.Params = &p
	}

	if j.BasicAuth != nil {
		if job.BasicAuth == nil {
			u := j.BasicAuth.Username.ValueString()
			p := j.BasicAuth.Password.ValueString()
			job.BasicAuth = &struct {
				Password *string `json:"password,omitempty"`
				Username *string `json:"username,omitempty"`
			}{
				Username: &u,
				Password: &p,
			}
		}
	}

	t := make([]struct {
		Labels  *map[string]interface{} `json:"labels,omitempty"`
		Targets []string                `json:"targets"`
	}, len(j.Targets))
	for i, target := range j.Targets {
		ti := struct {
			Labels  *map[string]interface{} `json:"labels,omitempty"`
			Targets []string                `json:"targets"`
		}{}
		ti.Targets = make([]string, len(target.URLs))
		for k, v := range target.URLs {
			ti.Targets[k] = v.ValueString()
		}

		ls := map[string]interface{}{}
		for k, v := range target.Labels.Elements() {
			ls[k], _ = common.ToString(context.TODO(), v)
		}
		ti.Labels = &ls
		t[i] = ti
	}
	job.StaticConfigs = t
	return job
}

func (j *Job) ToClientPartialUpdateJobs() scrapeconfig.PartialUpdateJSONBody {
	jobs := scrapeconfig.PartialUpdateJSONBody{{
		JobName:        j.Name.ValueString(),
		Scheme:         scrapeconfig.PartialUpdateJSONBodyScheme(j.Scheme.ValueString()),
		MetricsPath:    j.MetricsPath.ValueStringPointer(),
		ScrapeInterval: j.ScrapeInterval.ValueString(),
		ScrapeTimeout:  j.ScrapeTimeout.ValueString(),
		SampleLimit:    toFloat32Ptr(float32(j.SampleLimit.ValueInt64())),
	}}
	j.setDefaultsUpdate(jobs)

	job := &jobs[0]

	if j.SAML2 != nil && !j.SAML2.EnableURLParameters.ValueBool() {
		if job.Params == nil {
			job.Params = &map[string]interface{}{}
		}
		p := *job.Params
		p["saml2"] = []string{"disabled"}
		job.Params = &p
	}

	if j.BasicAuth != nil {
		if job.BasicAuth == nil {
			u := j.BasicAuth.Username.ValueString()
			p := j.BasicAuth.Password.ValueString()
			job.BasicAuth = &struct {
				Password *string `json:"password,omitempty"`
				Username *string `json:"username,omitempty"`
			}{
				Username: &u,
				Password: &p,
			}
		}
	}

	t := make([]struct {
		Labels  *map[string]interface{} `json:"labels,omitempty"`
		Targets []string                `json:"targets"`
	}, len(j.Targets))
	for i, target := range j.Targets {
		ti := struct {
			Labels  *map[string]interface{} `json:"labels,omitempty"`
			Targets []string                `json:"targets"`
		}{}
		ti.Targets = make([]string, len(target.URLs))
		for k, v := range target.URLs {
			ti.Targets[k] = v.ValueString()
		}

		ls := map[string]interface{}{}
		for k, v := range target.Labels.Elements() {
			ls[k], _ = common.ToString(context.TODO(), v)
		}
		ti.Labels = &ls
		t[i] = ti
	}
	job.StaticConfigs = t
	return jobs
}

func (j *Job) FromClientJob(cj scrapeconfig.Job) {
	j.ID = types.StringValue(cj.JobName)
	j.Name = types.StringValue(cj.JobName)
	if cj.MetricsPath != nil {
		j.MetricsPath = types.StringValue(*cj.MetricsPath)
	}
	if cj.Scheme != nil {
		j.Scheme = types.StringValue(string(*cj.Scheme))
	}
	j.ScrapeInterval = types.StringValue(cj.ScrapeInterval)
	j.ScrapeTimeout = types.StringValue(cj.ScrapeTimeout)
	j.SampleLimit = types.Int64Null()
	if cj.SampleLimit != nil {
		j.SampleLimit = types.Int64Value(int64(*cj.SampleLimit))
	}
	j.handleSAML2(cj)
	j.handleBasicAuth(cj)
	j.handleTargets(cj)
}

func (j *Job) handleBasicAuth(cj scrapeconfig.Job) {
	if cj.BasicAuth == nil {
		j.BasicAuth = nil
		return
	}
	j.BasicAuth = &BasicAuth{
		Username: types.StringValue(cj.BasicAuth.Username),
		Password: types.StringValue(cj.BasicAuth.Password),
	}
}

func (j *Job) handleSAML2(cj scrapeconfig.Job) {
	if cj.Params == nil && j.SAML2 == nil {
		return
	}

	if j.SAML2 == nil {
		j.SAML2 = &SAML2{}
	}

	flag := true
	if cj.Params == nil || *cj.Params == nil {
		return
	}
	p := *cj.Params
	if v, ok := p["saml2"]; ok {
		if len(v) == 1 && v[0] == "disabled" {
			flag = false
		}
	}

	j.SAML2 = &SAML2{
		EnableURLParameters: types.BoolValue(flag),
	}
}

func (j *Job) handleTargets(cj scrapeconfig.Job) {
	newTargets := []Target{}
	for i, sc := range cj.StaticConfigs {
		nt := Target{
			URLs: []types.String{},
		}
		for _, v := range sc.Targets {
			nt.URLs = append(nt.URLs, types.StringValue(v))
		}

		if len(j.Targets) > i && j.Targets[i].Labels.IsNull() || sc.Labels == nil {
			nt.Labels = types.MapNull(types.StringType)
		} else {
			newl := map[string]attr.Value{}
			for k, v := range *sc.Labels {
				newl[k] = types.StringValue(v)
			}
			nt.Labels = types.MapValueMust(types.StringType, newl)
		}
		newTargets = append(newTargets, nt)
	}
	j.Targets = newTargets
}

func toFloat32Ptr(v float32) *float32 {
	return &v
}
