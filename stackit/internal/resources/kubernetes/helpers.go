package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.4/generated/cluster"
	"github.com/pkg/errors"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	timeFormat                             = "2006-01-02T15:04:05.999Z"
	default_allow_privileged               = true
	default_os_name                        = "flatcar"
	default_nodepool_min             int64 = 1
	default_nodepool_max             int64 = 2
	default_nodepool_max_surge       int64 = 1
	default_nodepool_max_unavailable int64 = 1
	default_volume_type                    = "storage_premium_perf1"
	default_volume_size_gb           int64 = 20
	default_cri                            = "containerd"
	default_zone                           = "eu01-m"
	default_version                        = "1.24"
)

func (r Resource) loadAvaiableVersions(ctx context.Context) ([]*semver.Version, error) {
	c := r.client
	var versionOptions []*semver.Version
	resp, err := c.Services.Kubernetes.ProviderOptions.GetProviderOptionsWithResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed making options request")
	}
	if resp.HasError != nil {
		return nil, errors.Wrap(resp.HasError, "Couldn't fetch options")
	}

	opts := resp.JSON200
	versionOptions = []*semver.Version{}
	for _, v := range *opts.KubernetesVersions {
		if v.State == nil || v.Version == nil {
			continue
		}
		if !strings.EqualFold(*v.State, consts.SKE_VERSION_STATE_SUPPORTED) {
			continue
		}
		versionOption, err := semver.NewVersion(*v.Version)
		if err != nil {
			return nil, err
		}
		versionOptions = append(versionOptions, versionOption)
	}
	return versionOptions, nil
}

func (c *Cluster) clusterConfig(versionOptions []*semver.Version) (cluster.Kubernetes, error) {
	if c.KubernetesVersion.IsNull() || c.KubernetesVersion.IsUnknown() {
		c.KubernetesVersion = types.StringValue(default_version)
	}

	clusterConfigVersion, err := semver.NewVersion(c.KubernetesVersion.ValueString())
	if err != nil {
		return cluster.Kubernetes{}, err
	}
	clusterVersionConstraint, err := toVersionConstraint(clusterConfigVersion)
	if err != nil {
		return cluster.Kubernetes{}, err
	}
	clusterConfigVersion = maxVersionOption(clusterVersionConstraint, versionOptions)
	if clusterConfigVersion == nil {
		return cluster.Kubernetes{}, fmt.Errorf("returned version is nil\nthe options were: %+v", versionOptions)
	}

	pvlg := c.AllowPrivilegedContainers.ValueBool()
	cfg := cluster.Kubernetes{
		Version:                   clusterConfigVersion.String(),
		AllowPrivilegedContainers: &pvlg,
	}

	if c.AllowPrivilegedContainers.IsNull() || c.AllowPrivilegedContainers.IsUnknown() {
		pvlg := default_allow_privileged
		cfg.AllowPrivilegedContainers = &pvlg
	}
	return cfg, nil
}

// toVersionConstraint matches the patch version if given, or else any version with same major and minor version.
func toVersionConstraint(version *semver.Version) (*semver.Constraints, error) {
	if version.String() == version.Original() { // patch version given
		return semver.NewConstraint(fmt.Sprintf("= %s", version.String()))
	}
	nextVersion := version.IncMinor()
	return semver.NewConstraint(fmt.Sprintf(">= %s, < %s", version.String(), nextVersion.String()))
}

// maxVersionOption returns the maximal version that matches the given version. A matching option is required.
// If the given version only contains major and minor version, the latest patch version is returned.
func maxVersionOption(versionConstraint *semver.Constraints, versionOptions []*semver.Version) *semver.Version {
	if len(versionOptions) == 0 || versionOptions[0] == nil {
		return nil
	}
	ret := versionOptions[0]
	for _, v := range versionOptions[1:] {
		if versionConstraint.Check(v) && v.GreaterThan(ret) {
			ret = v
		}
	}
	return ret
}

func (c *Cluster) nodePools() []cluster.Nodepool {
	cnps := []cluster.Nodepool{}
	for _, p := range c.NodePools {
		// taints
		ts := []cluster.Taint{}
		for _, v := range p.Taints {
			val := v.Value.ValueString()
			t := cluster.Taint{
				Effect: cluster.TaintEffect(v.Effect.ValueString()),
				Key:    v.Key.ValueString(),
				Value:  &val,
			}
			ts = append(ts, t)
		}

		// labels
		ls := map[string]string{}
		for k, v := range p.Labels.Elems {
			nv, err := common.ToString(context.Background(), v)
			if err != nil {
				ls[k] = ""
				continue
			}
			ls[k] = nv
		}

		// zones
		zs := []string{}
		for _, v := range p.Zones.Elems {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			s, err := common.ToString(context.TODO(), v)
			if err != nil {
				continue
			}
			zs = append(zs, s)
		}

		ms := int(p.MaxSurge.ValueInt64())
		mu := int(p.MaxUnavailable.ValueInt64())
		in := p.OSName.ValueString()
		vt := p.VolumeType.ValueString()
		cn := cluster.CRIName(p.ContainerRuntime.ValueString())
		cnp := cluster.Nodepool{
			Name:           p.Name.ValueString(),
			Minimum:        int(p.Minimum.ValueInt64()),
			Maximum:        int(p.Maximum.ValueInt64()),
			MaxSurge:       &ms,
			MaxUnavailable: &mu,
			Machine: cluster.Machine{
				Type: p.MachineType.ValueString(),
				Image: cluster.Image{
					Name:    &in,
					Version: p.OSVersion.ValueString(),
				},
			},
			Volume: cluster.Volume{
				Type: &vt,
				Size: int(p.VolumeSizeGB.ValueInt64()),
			},
			Taints: &ts,
			CRI: &cluster.CRI{
				Name: &cn,
			},
			Labels:            &ls,
			AvailabilityZones: zs,
		}
		cnps = append(cnps, cnp)
	}
	return cnps
}

func setNodepoolDefaults(nps []cluster.Nodepool) []cluster.Nodepool {
	for i, np := range nps {
		if np.Machine.Image.Name == nil || *np.Machine.Image.Name == "" {
			d := default_os_name
			nps[i].Machine.Image.Name = &d
		}
		if np.Minimum == 0 {
			nps[i].Minimum = int(default_nodepool_min)
		}
		if np.Maximum == 0 {
			nps[i].Maximum = int(default_nodepool_max)
		}
		if np.MaxSurge == nil || *np.MaxSurge == 0 {
			d := int(default_nodepool_max_surge)
			nps[i].MaxSurge = &d
		}
		if np.MaxUnavailable == nil || *np.MaxUnavailable == 0 {
			d := int(default_nodepool_max_unavailable)
			nps[i].MaxUnavailable = &d
		}
		if np.Volume.Type == nil || *np.Volume.Type == "" {
			s := default_volume_type
			nps[i].Volume.Type = &s
		}
		if np.Volume.Size == 0 {
			nps[i].Volume.Size = int(default_volume_size_gb)
		}
		if np.CRI != nil && (np.CRI.Name == nil || *np.CRI.Name == "") {
			s := cluster.CRIName(default_cri)
			nps[i].CRI.Name = &s
		}
		if len(np.AvailabilityZones) == 0 {
			nps[i].AvailabilityZones = []string{default_zone}
		}
	}
	return nps
}

func (c *Cluster) hibernations() *cluster.Hibernation {
	scs := []cluster.HibernationSchedule{}
	for _, h := range c.Hibernations {
		sc := cluster.HibernationSchedule{
			Start: h.Start.ValueString(),
			End:   h.End.ValueString(),
		}
		if !h.Timezone.IsNull() && !h.Timezone.IsUnknown() {
			tz := h.Timezone.ValueString()
			sc.Timezone = &tz
		}
		scs = append(scs, sc)
	}

	if len(scs) == 0 {
		return nil
	}

	return &cluster.Hibernation{
		Schedules: scs,
	}
}

func (c *Cluster) extensions() *cluster.Extension {
	if c.Extensions == nil || c.Extensions.Argus == nil {
		return nil
	}

	return &cluster.Extension{
		Argus: &cluster.Argus{
			Enabled:         c.Extensions.Argus.Enabled.ValueBool(),
			ArgusInstanceID: c.Extensions.Argus.ArgusInstanceID.ValueString(),
		},
	}
}

func (c *Cluster) maintenance() *cluster.Maintenance {
	if c.Maintenance == nil {
		return nil
	}

	kv := c.Maintenance.EnableKubernetesVersionUpdates.ValueBool()
	miv := c.Maintenance.EnableMachineImageVersionUpdates.ValueBool()
	return &cluster.Maintenance{
		AutoUpdate: cluster.MaintenanceAutoUpdate{
			KubernetesVersion:   &kv,
			MachineImageVersion: &miv,
		},
		TimeWindow: cluster.TimeWindow{
			Start: c.Maintenance.Start.ValueString(),
			End:   c.Maintenance.End.ValueString(),
		},
	}
}

// Transform transforms clusters.Cluster structure to Cluster
func (c *Cluster) Transform(cl cluster.Cluster) {
	c.ID = types.StringValue(*cl.Name)
	if c.KubernetesVersion.IsNull() || c.KubernetesVersion.IsUnknown() {
		c.KubernetesVersion = types.StringValue(cl.Kubernetes.Version)
	}
	c.KubernetesVersionUsed = types.StringValue(cl.Kubernetes.Version)
	c.AllowPrivilegedContainers = types.Bool{Value: *cl.Kubernetes.AllowPrivilegedContainers}
	c.Status = types.StringValue(string(*cl.Status.Aggregated))
	c.NodePools = []NodePool{}
	for _, np := range cl.Nodepools {
		n := NodePool{
			Name:             types.StringValue(np.Name),
			MachineType:      types.StringValue(np.Machine.Type),
			OSName:           types.StringValue(*np.Machine.Image.Name),
			OSVersion:        types.StringValue(np.Machine.Image.Version),
			Minimum:          types.Int64Value(int64(np.Minimum)),
			Maximum:          types.Int64Value(int64(np.Maximum)),
			MaxSurge:         types.Int64Value(int64(*np.MaxSurge)),
			MaxUnavailable:   types.Int64Value(int64(*np.MaxUnavailable)),
			VolumeType:       types.StringValue(*np.Volume.Type),
			VolumeSizeGB:     types.Int64Value(int64(np.Volume.Size)),
			Labels:           types.Map{ElemType: types.StringType, Null: true},
			Taints:           nil,
			ContainerRuntime: types.StringValue(string(*np.CRI.Name)),
			Zones:            types.List{ElemType: types.StringType, Null: true},
		}
		if np.Labels != nil {
			for k, v := range *np.Labels {
				if n.Labels.Null {
					n.Labels.Null = false
					n.Labels.Elems = make(map[string]attr.Value, len(*np.Labels))
				}
				n.Labels.Elems[k] = types.StringValue(v)
			}
		}
		if np.Taints != nil {
			for _, v := range *np.Taints {
				if n.Taints == nil {
					n.Taints = []Taint{}
				}
				n.Taints = append(n.Taints, Taint{
					Effect: types.StringValue(string(v.Effect)),
					Key:    types.StringValue(v.Key),
					Value:  types.StringValue(*v.Value),
				})
			}
		}
		for _, v := range np.AvailabilityZones {
			if n.Zones.Null {
				n.Zones.Null = false
			}
			n.Zones.Elems = append(n.Zones.Elems, types.StringValue(v))
		}
		c.NodePools = append(c.NodePools, n)
	}

	c.transformMaintenance(cl)
	c.transformHibernations(cl)
	c.transformExtensions(cl)
}

func (c *Cluster) transformHibernations(cl cluster.Cluster) {
	if c.Hibernations == nil || cl.Hibernation == nil {
		return
	}

	c.Hibernations = []Hibernation{}
	for _, h := range cl.Hibernation.Schedules {
		c.Hibernations = append(c.Hibernations, Hibernation{
			Start:    types.StringValue(h.Start),
			End:      types.StringValue(h.End),
			Timezone: types.StringValue(*h.Timezone),
		})
	}
}

func (c *Cluster) transformMaintenance(cl cluster.Cluster) {
	if c.Maintenance == nil || cl.Maintenance == nil {
		return
	}

	c.Maintenance = &Maintenance{
		EnableKubernetesVersionUpdates:   types.Bool{Value: *cl.Maintenance.AutoUpdate.KubernetesVersion},
		EnableMachineImageVersionUpdates: types.Bool{Value: *cl.Maintenance.AutoUpdate.MachineImageVersion},
		Start:                            types.StringValue(cl.Maintenance.TimeWindow.Start),
		End:                              types.StringValue(cl.Maintenance.TimeWindow.End),
	}
}

func (c *Cluster) transformExtensions(cl cluster.Cluster) {
	if c.Extensions == nil || cl.Extensions == nil {
		return
	}

	if cl.Extensions.Argus != nil && c.Extensions.Argus != nil {
		c.Extensions.Argus = &ArgusExtension{
			Enabled:         types.Bool{Value: cl.Extensions.Argus.Enabled},
			ArgusInstanceID: types.StringValue(cl.Extensions.Argus.ArgusInstanceID),
		}
	}
}
