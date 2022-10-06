package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/options"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
)

func (r Resource) validate(
	ctx context.Context,
	projectID string,
	clusterName string,
	clusterConfig clusters.Kubernetes,
	nodePools []clusters.NodePool,
	maintenance *clusters.Maintenance,
	hibernation *clusters.Hibernation,
	extensions *clusters.Extensions,
) error {

	// General validation
	if err := validate.ProjectID(projectID); err != nil {
		return err
	}
	if err := clusters.ValidateCluster(clusterName, clusterConfig, nodePools, maintenance, hibernation, extensions); err != nil {
		return err
	}

	// Validate against real options
	c := r.Provider.Client()
	opts, err := c.Kubernetes.Options.List(ctx)
	if err != nil {
		// if options cannot be fetched, skip validation
		return nil
	}

	clusterConfig.Version = maxVersionOption(clusterConfig.Version, opts.KubernetesVersions)

	if err := validateKubernetesVersion(clusterConfig.Version, opts.KubernetesVersions); err != nil {
		return err
	}

	for _, np := range nodePools {
		if err := validateMachineImage(np.Machine.Image.Name, np.Machine.Image.Version, opts.MachineImages); err != nil {
			return err
		}
		if err := validateMachineType(np.Machine.Type, opts.MachineTypes); err != nil {
			return err
		}
		if err := validateVolumeType(np.Volume.Type, opts.VolumeTypes); err != nil {
			return err
		}
		if err := validateZones(np.AvailabilityZones, opts.AvailabilityZones); err != nil {
			return err
		}
	}

	return nil
}

// maxVersionOption returns the maximal version that matches the given version.
// If the given version only contains major and minor version, the latest patch version is returned.
func maxVersionOption(version string, versionOptions []options.KubernetesVersion) string {
	if strings.Count(version, ".") == 2 {
		return version
	}
	ret := version
	for _, v := range versionOptions {
		if len(v.Version) > len(version) && strings.HasPrefix(v.Version, version+".") && v.Version > ret {
			ret = v.Version
		}
	}
	return ret
}

func validateKubernetesVersion(version string, versionOptions []options.KubernetesVersion) error {
	found := false
	accepted := ""
	for _, v := range versionOptions {
		if strings.EqualFold(v.State, consts.SKE_VERSION_STATE_SUPPORTED) {
			if v.Version == version {
				found = true
				break
			}
			accepted = fmt.Sprintf("%s- %s (state: %s, expires: %s)\n", accepted, v.Version, v.State, v.ExpirationDate)
		}
	}
	if !found {
		return fmt.Errorf(
			"incorrect kubernetes version '%s'\naccepted options are:\n%s",
			version,
			accepted,
		)
	}
	return nil
}

func validateMachineImage(image, version string, imageOptions []options.MachineImage) error {
	foundImage := false
	foundVersion := false
	acceptedImages := ""
	acceptedVersions := ""
	for _, v := range imageOptions {
		if v.Name == image {
			foundImage = true
			for _, v2 := range v.Versions {
				if strings.EqualFold(v2.State, consts.SKE_VERSION_STATE_SUPPORTED) {
					if v2.Version == version {
						foundVersion = true
						break
					}
					acceptedVersions = fmt.Sprintf("%s- %s (state: %s, expires: %s)\n", acceptedVersions, v2.Version, v2.State, v2.ExpirationDate)
				}
			}

		}
		acceptedImages = fmt.Sprintf("%s- %s (versions: %v)\n", acceptedImages, v.Name, v.Versions)
	}
	if !foundImage {
		return fmt.Errorf(
			"incorrect machine image '%s'\naccepted options are:\n%v",
			image,
			imageOptions,
		)
	}
	if !foundVersion {
		return fmt.Errorf(
			"incorrect version '%s'\naccepted options are:\n%s",
			version,
			acceptedVersions,
		)
	}
	return nil
}

func validateMachineType(machine string, machineTypes []options.MachineType) error {
	found := false
	accepted := ""
	for _, v := range machineTypes {
		if v.Name == machine {
			found = true
			break
		}
		accepted = fmt.Sprintf("%s- %s (cpu: %d, mem: %d)\n", accepted, v.Name, v.CPU, v.Memory)
	}
	if !found {
		return fmt.Errorf(
			"incorrect machine '%s'\naccepted options are:\n%s",
			machine,
			accepted,
		)
	}
	return nil
}

func validateVolumeType(volume string, volumeTypes []options.VolumeType) error {
	found := false
	accepted := ""
	for _, v := range volumeTypes {
		if v.Name == volume {
			found = true
			break
		}
		accepted = fmt.Sprintf("%s- %s\n", accepted, v.Name)
	}
	if !found {
		return fmt.Errorf(
			"incorrect volume type '%s'\naccepted options are:\n%s",
			volume,
			accepted,
		)
	}
	return nil
}

func validateZones(zones []string, zoneOptions []options.AvailabilityZone) error {
	var found bool
	accepted := ""
	for _, v := range zoneOptions {
		accepted = fmt.Sprintf("%s- %s\n", accepted, v.Name)
	}
	if len(zones) == 0 {
		return fmt.Errorf(
			"please specify a list of zones\naccepted options are:\n%s",
			accepted,
		)
	}

	for _, v := range zones {
		found = false
		for _, v2 := range zoneOptions {
			if v == v2.Name {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf(
				"incorrect zone '%s'\naccepted options are:\n%s",
				v,
				accepted,
			)
		}
	}
	return nil
}
