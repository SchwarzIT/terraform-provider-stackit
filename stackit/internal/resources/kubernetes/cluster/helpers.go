package cluster

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.1/cluster"
	"github.com/pkg/errors"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	timeFormat                          = "2006-01-02T15:04:05.999Z"
	DefaultAllowPrivileged              = true
	DefaultOSName                       = "flatcar"
	DefaultNodepoolMin            int64 = 1
	DefaultNodepoolMax            int64 = 2
	DefaultNodepoolMaxSurge       int64 = 1
	DefaultNodepoolMaxUnavailable int64 = 0
	DefaultVolumeType                   = "storage_premium_perf1"
	DefaultVolumeSizeGB           int64 = 20
	DefaultCRI                          = "containerd"
	DefaultZone                         = "eu01-m"
	DefaultVersion                      = "1.31"
)

func (r Resource) loadAvaiableVersions(ctx context.Context, diags *diag.Diagnostics) ([]*semver.Version, error) {
	c := r.client
	var versionOptions []*semver.Version
	res, err := c.Kubernetes.ProviderOptions.List(ctx)
	if agg := common.Validate(diags, res, err, "JSON200.KubernetesVersions"); agg != nil {
		return nil, errors.Wrap(agg, "failed fetching cluster versions")
	}

	opts := res.JSON200
	versionOptions = []*semver.Version{}
	for _, v := range *opts.KubernetesVersions {
		if v.State == nil || v.Version == nil {
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
		c.KubernetesVersion = types.StringValue(DefaultVersion)
	}

	clusterConfigVersion, err := semver.NewVersion(c.KubernetesVersion.ValueString())
	if err != nil {
		return cluster.Kubernetes{}, err
	}

	if c.KubernetesVersionUsed.ValueString() != "" {
		clusterCurrentVersionUsed, err := semver.NewVersion(c.KubernetesVersionUsed.ValueString())
		if err != nil {
			return cluster.Kubernetes{}, err
		}

		if clusterCurrentVersionUsed.GreaterThan(clusterConfigVersion) {
			clusterConfigVersion = clusterCurrentVersionUsed
		}
	} else {
		clusterVersionConstraint, err := toVersionConstraint(clusterConfigVersion)
		if err != nil {
			return cluster.Kubernetes{}, err
		}

		clusterConfigVersion = maxVersionOption(clusterVersionConstraint, versionOptions)
		if clusterConfigVersion == nil {
			return cluster.Kubernetes{}, fmt.Errorf("returned version is nil\nthe options were: %+v", versionOptions)
		}
	}

	pvlg := c.AllowPrivilegedContainers.ValueBool()
	cfg := cluster.Kubernetes{
		Version: clusterConfigVersion.String(),
	}

	if clusterConfigVersion.Compare(semver.MustParse("1.25.0")) == -1 {
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
		for k, v := range p.Labels.Elements() {
			nv, err := common.ToString(context.Background(), v)
			if err != nil {
				ls[k] = ""
				continue
			}
			ls[k] = nv
		}

		// zones
		zs := []string{}
		for _, v := range p.Zones.Elements() {
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
			d := DefaultOSName
			nps[i].Machine.Image.Name = &d
		}
		if np.Minimum == 0 {
			nps[i].Minimum = int(DefaultNodepoolMin)
		}
		if np.Maximum == 0 {
			nps[i].Maximum = int(DefaultNodepoolMax)
		}
		if np.MaxSurge == nil || *np.MaxSurge == 0 {
			d := int(DefaultNodepoolMaxSurge)
			nps[i].MaxSurge = &d
		}
		if np.MaxUnavailable == nil || *np.MaxUnavailable == 0 {
			d := int(DefaultNodepoolMaxUnavailable)
			nps[i].MaxUnavailable = &d
		}
		if np.Volume.Type == nil || *np.Volume.Type == "" {
			s := DefaultVolumeType
			nps[i].Volume.Type = &s
		}
		if np.Volume.Size == 0 {
			nps[i].Volume.Size = int(DefaultVolumeSizeGB)
		}
		if np.CRI != nil && (np.CRI.Name == nil || *np.CRI.Name == "") {
			s := cluster.CRIName(DefaultCRI)
			nps[i].CRI.Name = &s
		}
		if len(np.AvailabilityZones) == 0 {
			nps[i].AvailabilityZones = []string{DefaultZone}
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

func (c *Cluster) extensions(ctx context.Context) (*cluster.Extension, diag.Diagnostics) {
	if c.Extensions == nil {
		return nil, nil
	}
	ex := &cluster.Extension{}
	if c.Extensions.Argus != nil {
		ex.Argus = &cluster.Argus{
			Enabled:         c.Extensions.Argus.Enabled.ValueBool(),
			ArgusInstanceID: c.Extensions.Argus.ArgusInstanceID.ValueString(),
		}
	}
	if c.Extensions.ACL != nil {
		var cidrs []string
		diags := c.Extensions.ACL.AllowedCIDRs.ElementsAs(ctx, &cidrs, true)
		if diags.HasError() {
			return nil, diags
		}
		ex.Acl = &cluster.ACL{
			Enabled:      c.Extensions.ACL.Enabled.ValueBool(),
			AllowedCidrs: cidrs,
		}
	}
	return ex, nil
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

	if cl.Kubernetes.AllowPrivilegedContainers != nil {
		c.AllowPrivilegedContainers = types.BoolValue(*cl.Kubernetes.AllowPrivilegedContainers)
	} else {
		c.AllowPrivilegedContainers = types.BoolValue(DefaultAllowPrivileged)
	}

	c.Status = types.StringValue(string(*cl.Status.Aggregated))
	c.NodePools = []NodePool{}

	for _, np := range cl.Nodepools {
		maimna := types.StringNull()
		if np.Machine.Image.Name != nil {
			maimna = types.StringValue(*np.Machine.Image.Name)
		}
		ms := types.Int64Null()
		if np.MaxSurge != nil {
			ms = types.Int64Value(int64(*np.MaxSurge))
		}
		mu := types.Int64Null()
		if np.MaxUnavailable != nil {
			mu = types.Int64Value(int64(*np.MaxUnavailable))
		}
		vt := types.StringNull()
		if np.Volume.Type != nil {
			vt = types.StringValue(*np.Volume.Type)
		}
		crin := types.StringNull()
		if np.CRI.Name != nil {
			crin = types.StringValue(string(*np.CRI.Name))
		}
		n := NodePool{
			Name:             types.StringValue(np.Name),
			MachineType:      types.StringValue(np.Machine.Type),
			OSName:           maimna,
			OSVersion:        types.StringValue(np.Machine.Image.Version),
			Minimum:          types.Int64Value(int64(np.Minimum)),
			Maximum:          types.Int64Value(int64(np.Maximum)),
			MaxSurge:         ms,
			MaxUnavailable:   mu,
			VolumeType:       vt,
			VolumeSizeGB:     types.Int64Value(int64(np.Volume.Size)),
			Labels:           types.MapNull(types.StringType),
			Taints:           nil,
			ContainerRuntime: crin,
			Zones:            types.ListNull(types.StringType),
		}

		if np.Labels != nil {
			elems := map[string]attr.Value{}
			for k, v := range *np.Labels {
				elems[k] = types.StringValue(v)
			}
			n.Labels = types.MapValueMust(types.StringType, elems)
		}

		if np.Taints != nil {
			for _, v := range *np.Taints {
				if n.Taints == nil {
					n.Taints = []Taint{}
				}
				taintval := types.StringNull()
				if v.Value != nil {
					taintval = types.StringValue(*v.Value)
				}
				n.Taints = append(n.Taints, Taint{
					Effect: types.StringValue(string(v.Effect)),
					Key:    types.StringValue(v.Key),
					Value:  taintval,
				})
			}
		}

		elems := []attr.Value{}
		for _, v := range np.AvailabilityZones {
			elems = append(elems, types.StringValue(v))
		}

		n.Zones = types.ListValueMust(types.StringType, elems)
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

	ekvu := types.BoolNull()
	if cl.Maintenance.AutoUpdate.KubernetesVersion != nil {
		ekvu = types.BoolValue(*cl.Maintenance.AutoUpdate.KubernetesVersion)
	}
	emvu := types.BoolNull()
	if cl.Maintenance.AutoUpdate.KubernetesVersion != nil {
		emvu = types.BoolValue(*cl.Maintenance.AutoUpdate.MachineImageVersion)
	}

	c.Maintenance = &Maintenance{
		EnableKubernetesVersionUpdates:   ekvu,
		EnableMachineImageVersionUpdates: emvu,
		Start:                            types.StringValue(cl.Maintenance.TimeWindow.Start),
		End:                              types.StringValue(cl.Maintenance.TimeWindow.End),
	}
}

func (c *Cluster) transformExtensions(cl cluster.Cluster) {
	if c.Extensions == nil || cl.Extensions == nil {
		return
	}

	if cl.Extensions.Argus != nil {
		c.Extensions.Argus = &ArgusExtension{
			Enabled:         types.BoolValue(cl.Extensions.Argus.Enabled),
			ArgusInstanceID: types.StringValue(cl.Extensions.Argus.ArgusInstanceID),
		}
	}

	if cl.Extensions.Acl != nil {
		cidr := []attr.Value{}
		for _, v := range cl.Extensions.Acl.AllowedCidrs {
			cidr = append(cidr, types.StringValue(v))
		}
		c.Extensions.ACL = &ACL{
			Enabled:      types.BoolValue(cl.Extensions.Acl.Enabled),
			AllowedCIDRs: types.ListValueMust(types.StringType, cidr),
		}
	}
}
