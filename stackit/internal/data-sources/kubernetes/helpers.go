package kubernetes

import (
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Transform transforms clusters.Cluster structure to Cluster
func transform(c *kubernetes.Cluster, cl clusters.Cluster) {
	c.ID = types.String{Value: cl.Name}
	c.KubernetesVersion = types.String{Value: cl.Kubernetes.Version}
	c.KubernetesVersionUsed = types.String{Value: cl.Kubernetes.Version}
	c.AllowPrivilegedContainers = types.Bool{Value: cl.Kubernetes.AllowPrivilegedContainers}
	c.Status = types.String{Value: cl.Status.Aggregated}

	transformNodepools(c, cl)
	transformMaintenance(c, cl)
	transformHibernations(c, cl)
	transformExtensions(c, cl)
}

func transformNodepools(c *kubernetes.Cluster, cl clusters.Cluster) {
	c.NodePools = []kubernetes.NodePool{}
	for _, np := range cl.Nodepools {
		n := kubernetes.NodePool{
			Name:             types.String{Value: np.Name},
			MachineType:      types.String{Value: np.Machine.Type},
			OSName:           types.String{Value: np.Machine.Image.Name},
			OSVersion:        types.String{Value: np.Machine.Image.Version},
			Minimum:          types.Int64{Value: int64(np.Minimum)},
			Maximum:          types.Int64{Value: int64(np.Maximum)},
			MaxSurge:         types.Int64{Value: int64(np.MaxSurge)},
			MaxUnavailable:   types.Int64{Value: int64(np.MaxUnavailable)},
			VolumeType:       types.String{Value: np.Volume.Type},
			VolumeSizeGB:     types.Int64{Value: int64(np.Volume.Size)},
			Labels:           types.Map{ElemType: types.StringType, Null: true},
			Taints:           nil,
			ContainerRuntime: types.String{Value: np.CRI.Name},
			Zones:            types.List{ElemType: types.StringType, Null: true},
		}
		for k, v := range np.Labels {
			if n.Labels.Null {
				n.Labels.Null = false
				n.Labels.Elems = make(map[string]attr.Value, len(np.Labels))
			}
			n.Labels.Elems[k] = types.String{Value: v}
		}
		for _, v := range np.Taints {
			if n.Taints == nil {
				n.Taints = []kubernetes.Taint{}
			}
			n.Taints = append(n.Taints, kubernetes.Taint{
				Effect: types.String{Value: v.Effect},
				Key:    types.String{Value: v.Key},
				Value:  types.String{Value: v.Value},
			})
		}
		for _, v := range np.AvailabilityZones {
			if n.Zones.Null {
				n.Zones.Null = false
			}
			n.Zones.Elems = append(n.Zones.Elems, types.String{Value: v})
		}
		c.NodePools = append(c.NodePools, n)
	}
}

func transformHibernations(c *kubernetes.Cluster, cl clusters.Cluster) {
	c.Hibernations = []kubernetes.Hibernation{}

	if cl.Hibernation == nil {
		return
	}

	for _, h := range cl.Hibernation.Schedules {
		c.Hibernations = append(c.Hibernations, kubernetes.Hibernation{
			Start:    types.String{Value: h.Start},
			End:      types.String{Value: h.End},
			Timezone: types.String{Value: h.Timezone},
		})
	}
}

func transformMaintenance(c *kubernetes.Cluster, cl clusters.Cluster) {
	c.Maintenance = &kubernetes.Maintenance{}

	if cl.Maintenance == nil {
		return
	}

	c.Maintenance = &kubernetes.Maintenance{
		EnableKubernetesVersionUpdates:   types.Bool{Value: cl.Maintenance.AutoUpdate.KubernetesVersion},
		EnableMachineImageVersionUpdates: types.Bool{Value: cl.Maintenance.AutoUpdate.MachineImageVersion},
		Start:                            types.String{Value: cl.Maintenance.TimeWindow.Start},
		End:                              types.String{Value: cl.Maintenance.TimeWindow.End},
	}
}

func transformExtensions(c *kubernetes.Cluster, cl clusters.Cluster) {
	c.Extensions = &kubernetes.Extensions{
		Argus: &kubernetes.ArgusExtension{
			Enabled:         types.Bool{},
			ArgusInstanceID: types.String{},
		},
	}

	if cl.Extensions == nil {
		return
	}

	if cl.Extensions.Argus != nil {
		c.Extensions.Argus = &kubernetes.ArgusExtension{
			Enabled:         types.Bool{Value: cl.Extensions.Argus.Enabled},
			ArgusInstanceID: types.String{Value: cl.Extensions.Argus.ArgusInstanceID},
		}
	}
}
