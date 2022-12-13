package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
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
	default_version                        = "1.23"
)

func (r Resource) loadAvaiableVersions(ctx context.Context) ([]*semver.Version, error) {
	c := r.client
	var versionOptions []*semver.Version
	opts, err := c.Kubernetes.Options.List(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't fetch K8s options")
	}

	versionOptions = []*semver.Version{}
	for _, v := range opts.KubernetesVersions {
		if !strings.EqualFold(v.State, consts.SKE_VERSION_STATE_SUPPORTED) {
			continue
		}
		versionOption, err := semver.NewVersion(v.Version)
		if err != nil {
			return nil, err
		}
		versionOptions = append(versionOptions, versionOption)
	}
	return versionOptions, nil
}

func (c *Cluster) clusterConfig(versionOptions []*semver.Version) (clusters.Kubernetes, error) {
	if c.KubernetesVersion.IsNull() || c.KubernetesVersion.IsUnknown() {
		c.KubernetesVersion = types.StringValue(default_version)
	}

	clusterConfigVersion, err := semver.NewVersion(c.KubernetesVersion.ValueString())
	if err != nil {
		return clusters.Kubernetes{}, err
	}
	clusterVersionConstraint, err := toVersionConstraint(clusterConfigVersion)
	if err != nil {
		return clusters.Kubernetes{}, err
	}
	clusterConfigVersion = maxVersionOption(clusterVersionConstraint, versionOptions)
	if clusterConfigVersion == nil {
		return clusters.Kubernetes{}, fmt.Errorf("returned version is nil\nthe options were: %+v", versionOptions)
	}

	cfg := clusters.Kubernetes{
		Version:                   clusterConfigVersion.String(),
		AllowPrivilegedContainers: c.AllowPrivilegedContainers.ValueBool(),
	}

	if c.AllowPrivilegedContainers.IsNull() || c.AllowPrivilegedContainers.IsUnknown() {
		cfg.AllowPrivilegedContainers = default_allow_privileged
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

func (c *Cluster) nodePools() []clusters.NodePool {
	cnps := []clusters.NodePool{}
	for _, p := range c.NodePools {
		// taints
		ts := []clusters.Taint{}
		for _, v := range p.Taints {
			t := clusters.Taint{
				Effect: v.Effect.ValueString(),
				Key:    v.Key.ValueString(),
				Value:  v.Value.ValueString(),
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

		cnp := clusters.NodePool{
			Name:           p.Name.ValueString(),
			Minimum:        int(p.Minimum.ValueInt64()),
			Maximum:        int(p.Maximum.ValueInt64()),
			MaxSurge:       int(p.MaxSurge.ValueInt64()),
			MaxUnavailable: int(p.MaxUnavailable.ValueInt64()),
			Machine: clusters.Machine{
				Type: p.MachineType.ValueString(),
				Image: clusters.MachineImage{
					Name:    p.OSName.ValueString(),
					Version: p.OSVersion.ValueString(),
				},
			},
			Volume: clusters.Volume{
				Type: p.VolumeType.ValueString(),
				Size: int(p.VolumeSizeGB.ValueInt64()),
			},
			Taints: ts,
			CRI: clusters.CRI{
				Name: p.ContainerRuntime.ValueString(),
			},
			Labels:            ls,
			AvailabilityZones: zs,
		}
		cnps = append(cnps, cnp)
	}
	return cnps
}

func setNodepoolDefaults(nps []clusters.NodePool) []clusters.NodePool {
	for i, np := range nps {
		if np.Machine.Image.Name == "" {
			nps[i].Machine.Image.Name = default_os_name
		}
		if np.Minimum == 0 {
			nps[i].Minimum = int(default_nodepool_min)
		}
		if np.Maximum == 0 {
			nps[i].Maximum = int(default_nodepool_max)
		}
		if np.MaxSurge == 0 {
			nps[i].MaxSurge = int(default_nodepool_max_surge)
		}
		if np.MaxUnavailable == 0 {
			nps[i].MaxUnavailable = int(default_nodepool_max_unavailable)
		}
		if np.Volume.Type == "" {
			nps[i].Volume.Type = default_volume_type
		}
		if np.Volume.Size == 0 {
			nps[i].Volume.Size = int(default_volume_size_gb)
		}
		if np.CRI.Name == "" {
			nps[i].CRI.Name = default_cri
		}
		if len(np.AvailabilityZones) == 0 {
			nps[i].AvailabilityZones = []string{default_zone}
		}
	}
	return nps
}

func (c *Cluster) hibernations() *clusters.Hibernation {
	scs := []clusters.HibernationScedule{}
	for _, h := range c.Hibernations {
		sc := clusters.HibernationScedule{
			Start: h.Start.ValueString(),
			End:   h.End.ValueString(),
		}
		if !h.Timezone.IsNull() && !h.Timezone.IsUnknown() {
			sc.Timezone = h.Timezone.ValueString()
		}
		scs = append(scs, sc)
	}

	if len(scs) == 0 {
		return nil
	}

	return &clusters.Hibernation{
		Schedules: scs,
	}
}

func (c *Cluster) extensions() *clusters.Extensions {
	if c.Extensions == nil || c.Extensions.Argus == nil {
		return nil
	}

	return &clusters.Extensions{
		Argus: &clusters.ArgusExtension{
			Enabled:         c.Extensions.Argus.Enabled.ValueBool(),
			ArgusInstanceID: c.Extensions.Argus.ArgusInstanceID.ValueString(),
		},
	}
}

func (c *Cluster) maintenance() *clusters.Maintenance {
	if c.Maintenance == nil {
		return nil
	}

	return &clusters.Maintenance{
		AutoUpdate: clusters.MaintenanceAutoUpdate{
			KubernetesVersion:   c.Maintenance.EnableKubernetesVersionUpdates.ValueBool(),
			MachineImageVersion: c.Maintenance.EnableMachineImageVersionUpdates.ValueBool(),
		},
		TimeWindow: clusters.MaintenanceTimeWindow{
			Start: c.Maintenance.Start.ValueString(),
			End:   c.Maintenance.End.ValueString(),
		},
	}
}

// Transform transforms clusters.Cluster structure to Cluster
func (c *Cluster) Transform(cl clusters.Cluster) {
	c.ID = types.StringValue(cl.Name)
	if c.KubernetesVersion.IsNull() || c.KubernetesVersion.IsUnknown() {
		c.KubernetesVersion = types.StringValue(cl.Kubernetes.Version)
	}
	c.KubernetesVersionUsed = types.StringValue(cl.Kubernetes.Version)
	c.AllowPrivilegedContainers = types.Bool{Value: cl.Kubernetes.AllowPrivilegedContainers}
	c.Status = types.StringValue(cl.Status.Aggregated)
	c.NodePools = []NodePool{}
	for _, np := range cl.Nodepools {
		n := NodePool{
			Name:             types.StringValue(np.Name),
			MachineType:      types.StringValue(np.Machine.Type),
			OSName:           types.StringValue(np.Machine.Image.Name),
			OSVersion:        types.StringValue(np.Machine.Image.Version),
			Minimum:          types.Int64Value(int64(np.Minimum)),
			Maximum:          types.Int64Value(int64(np.Maximum)),
			MaxSurge:         types.Int64Value(int64(np.MaxSurge)),
			MaxUnavailable:   types.Int64Value(int64(np.MaxUnavailable)),
			VolumeType:       types.StringValue(np.Volume.Type),
			VolumeSizeGB:     types.Int64Value(int64(np.Volume.Size)),
			Labels:           types.Map{ElemType: types.StringType, Null: true},
			Taints:           nil,
			ContainerRuntime: types.StringValue(np.CRI.Name),
			Zones:            types.List{ElemType: types.StringType, Null: true},
		}
		for k, v := range np.Labels {
			if n.Labels.Null {
				n.Labels.Null = false
				n.Labels.Elems = make(map[string]attr.Value, len(np.Labels))
			}
			n.Labels.Elems[k] = types.StringValue(v)
		}
		for _, v := range np.Taints {
			if n.Taints == nil {
				n.Taints = []Taint{}
			}
			n.Taints = append(n.Taints, Taint{
				Effect: types.StringValue(v.Effect),
				Key:    types.StringValue(v.Key),
				Value:  types.StringValue(v.Value),
			})
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

func (c *Cluster) transformHibernations(cl clusters.Cluster) {
	if c.Hibernations == nil || cl.Hibernation == nil {
		return
	}

	c.Hibernations = []Hibernation{}
	for _, h := range cl.Hibernation.Schedules {
		c.Hibernations = append(c.Hibernations, Hibernation{
			Start:    types.StringValue(h.Start),
			End:      types.StringValue(h.End),
			Timezone: types.StringValue(h.Timezone),
		})
	}
}

func (c *Cluster) transformMaintenance(cl clusters.Cluster) {
	if c.Maintenance == nil || cl.Maintenance == nil {
		return
	}

	c.Maintenance = &Maintenance{
		EnableKubernetesVersionUpdates:   types.Bool{Value: cl.Maintenance.AutoUpdate.KubernetesVersion},
		EnableMachineImageVersionUpdates: types.Bool{Value: cl.Maintenance.AutoUpdate.MachineImageVersion},
		Start:                            types.StringValue(cl.Maintenance.TimeWindow.Start),
		End:                              types.StringValue(cl.Maintenance.TimeWindow.End),
	}
}

func (c *Cluster) transformExtensions(cl clusters.Cluster) {
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
