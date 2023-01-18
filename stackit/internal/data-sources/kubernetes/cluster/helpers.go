package cluster

import (
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/generated/cluster"
	kubernetesCluster "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes/cluster"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Transform transforms cluster.Cluster structure to Cluster
func transform(c *kubernetesCluster.Cluster, cl *cluster.Cluster) {
	if cl.Name != nil {
		c.ID = types.StringValue(*cl.Name)
	}
	c.KubernetesVersion = types.StringValue(cl.Kubernetes.Version)
	c.KubernetesVersionUsed = types.StringValue(cl.Kubernetes.Version)
	if cl.Kubernetes.AllowPrivilegedContainers != nil {
		c.AllowPrivilegedContainers = types.BoolValue(*cl.Kubernetes.AllowPrivilegedContainers)
	}
	if cl.Status.Aggregated != nil {
		c.Status = types.StringValue(string(*cl.Status.Aggregated))
	}
	transformNodepools(c, cl)
	transformMaintenance(c, cl)
	transformHibernations(c, cl)
	transformExtensions(c, cl)
}

func transformNodepools(c *kubernetesCluster.Cluster, cl *cluster.Cluster) {
	c.NodePools = []kubernetesCluster.NodePool{}
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
		n := kubernetesCluster.NodePool{
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
					n.Taints = []kubernetesCluster.Taint{}
				}
				taintval := types.StringNull()
				if v.Value != nil {
					taintval = types.StringValue(*v.Value)
				}
				n.Taints = append(n.Taints, kubernetesCluster.Taint{
					Effect: types.StringValue(string(v.Effect)),
					Key:    types.StringValue(v.Key),
					Value:  taintval,
				})
			}
		}
		zones := []attr.Value{}
		for _, v := range np.AvailabilityZones {
			zones = append(zones, types.StringValue(v))
		}
		n.Zones = types.ListValueMust(types.StringType, zones)
		c.NodePools = append(c.NodePools, n)
	}
}

func transformHibernations(c *kubernetesCluster.Cluster, cl *cluster.Cluster) {
	if cl.Hibernation == nil {
		return
	}

	c.Hibernations = []kubernetesCluster.Hibernation{}
	for _, h := range cl.Hibernation.Schedules {
		c.Hibernations = append(c.Hibernations, kubernetesCluster.Hibernation{
			Start:    types.StringValue(h.Start),
			End:      types.StringValue(h.End),
			Timezone: types.StringValue(*h.Timezone),
		})
	}
}

func transformMaintenance(c *kubernetesCluster.Cluster, cl *cluster.Cluster) {
	if cl.Maintenance == nil {
		return
	}

	eKvu := types.BoolNull()
	if cl.Maintenance.AutoUpdate.KubernetesVersion != nil {
		eKvu = types.BoolValue(*cl.Maintenance.AutoUpdate.KubernetesVersion)
	}

	eMiv := types.BoolNull()
	if cl.Maintenance.AutoUpdate.MachineImageVersion != nil {
		eMiv = types.BoolValue(*cl.Maintenance.AutoUpdate.MachineImageVersion)
	}
	c.Maintenance = &kubernetesCluster.Maintenance{
		EnableKubernetesVersionUpdates:   eKvu,
		EnableMachineImageVersionUpdates: eMiv,
		Start:                            types.StringValue(cl.Maintenance.TimeWindow.Start),
		End:                              types.StringValue(cl.Maintenance.TimeWindow.End),
	}
}

func transformExtensions(c *kubernetesCluster.Cluster, cl *cluster.Cluster) {
	if c.Extensions == nil || cl.Extensions == nil {
		return
	}
	c.Extensions = &kubernetesCluster.Extensions{
		Argus: &kubernetesCluster.ArgusExtension{
			Enabled:         types.BoolValue(false),
			ArgusInstanceID: types.StringNull(),
		},
		ACL: types.ObjectValueMust(map[string]attr.Type{
			"enabled":       types.BoolType,
			"allowed_cidrs": types.ListType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"enabled":       types.BoolValue(false),
			"allowed_cidrs": types.ListNull(types.StringType),
		}),
	}
	if cl.Extensions.Argus != nil {
		c.Extensions.Argus = &kubernetesCluster.ArgusExtension{
			Enabled:         types.BoolValue(cl.Extensions.Argus.Enabled),
			ArgusInstanceID: types.StringValue(cl.Extensions.Argus.ArgusInstanceID),
		}
	}

	if cl.Extensions.Acl != nil {
		cidr := []attr.Value{}
		for _, v := range cl.Extensions.Acl.AllowedCidrs {
			cidr = append(cidr, types.StringValue(v))
		}
		cidrList, diags := types.ListValue(types.StringType, cidr)
		if diags.HasError() {
			cidrList = types.ListNull(types.StringType)
		}
		c.Extensions.ACL = types.ObjectValueMust(map[string]attr.Type{
			"enabled":       types.BoolType,
			"allowed_cidrs": types.ListType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"enabled":       types.BoolValue(cl.Extensions.Acl.Enabled),
			"allowed_cidrs": cidrList,
		})
	}
}
