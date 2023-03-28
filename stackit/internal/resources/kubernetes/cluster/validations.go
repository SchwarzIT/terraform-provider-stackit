package cluster

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/cluster"
	provideroptions "github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/provider-options"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
)

func (r Resource) validate(
	ctx context.Context,
	projectID string,
	clusterName string,
	clusterConfig cluster.Kubernetes,
	nodePools *[]cluster.Nodepool,
	maintenance *cluster.Maintenance,
	hibernation *cluster.Hibernation,
	extensions *cluster.Extension,
) error {

	// General validation
	if err := validate.ProjectID(projectID); err != nil {
		return err
	}

	// Validate against real options
	c := r.client
	opts, err := c.Kubernetes.ProviderOptions.List(ctx)

	if agg := validate.Response(opts, err, "JSON200.KubernetesVersions"); agg != nil {
		// if options cannot be fetched, skip validation
		return nil
	}

	if err := validateKubernetesVersion(clusterConfig.Version, *opts.JSON200.KubernetesVersions); err != nil {
		return err
	}

	for i, np := range *nodePools {
		imageName := ""
		if np.Machine.Image.Name != nil {
			imageName = *np.Machine.Image.Name
		}
		versionOption, err := validateMachineImage(imageName, np.Machine.Image.Version, opts.JSON200.MachineImages)
		if err != nil {
			return err
		}
		if np.Machine.Image.Version == "" {
			(*nodePools)[i].Machine.Image.Version = versionOption
		}
		if err := validateMachineType(np.Machine.Type, opts.JSON200.MachineTypes); err != nil {
			return err
		}
		volType := ""
		if np.Volume.Type != nil {
			volType = *np.Volume.Type
		}
		if err := validateVolumeType(volType, opts.JSON200.VolumeTypes); err != nil {
			return err
		}
		if err := validateZones(np.AvailabilityZones, opts.JSON200.AvailabilityZones); err != nil {
			return err
		}
	}

	// General cluster validations
	if err := cluster.Validate(clusterName, clusterConfig, *nodePools, maintenance, hibernation, extensions); err != nil {
		return err
	}

	return nil
}

func validateKubernetesVersion(version string, versionOptions []provideroptions.KubernetesVersion) error {
	found := false
	accepted := ""
	for _, v := range versionOptions {
		if v.Version == nil {
			continue
		}
		if *v.Version == version {
			found = true
			break
		}
		ed := ""
		if v.ExpirationDate != nil {
			ed = *v.ExpirationDate
		}
		s := ""
		if v.State != nil {
			s = *v.State
		}
		accepted = fmt.Sprintf("%s- %s (state: %s, expires: %s)\n", accepted, *v.Version, s, ed)
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

func validateMachineImage(image, version string, imageOptions *[]provideroptions.MachineImage) (versionOption string, err error) {
	if imageOptions == nil {
		return "", errors.New("received empty machine image list")
	}
	foundImage := false
	foundVersion := false
	acceptedImages := ""
	acceptedVersions := ""
	supportedVersion := ""
	for _, v := range *imageOptions {
		if v.Name == nil || v.Versions == nil {
			continue
		}
		if *v.Name == image {
			foundImage = true
			for _, v2 := range *v.Versions {
				if v2.State != nil && strings.EqualFold(*v2.State, "supported") {
					if v2.Version == nil {
						continue
					}
					if supportedVersion == "" {
						supportedVersion = *v2.Version
					}
					if *v2.Version == version {
						foundVersion = true
						break
					}
					ed := ""
					if v2.ExpirationDate != nil {
						ed = *v2.ExpirationDate
					}
					acceptedVersions = fmt.Sprintf("%s- %s (state: %s, expires: %s)\n", acceptedVersions, *v2.Version, *v2.State, ed)
				}
			}

		}
		acceptedImages = fmt.Sprintf("%s- %s (versions: %v)\n", acceptedImages, *v.Name, *v.Versions)
	}
	if !foundImage {
		return "", fmt.Errorf(
			"incorrect machine image '%s'\naccepted options are:\n%v",
			image,
			imageOptions,
		)
	}
	if !foundVersion {
		if version != "" {
			return "", fmt.Errorf(
				"incorrect version '%s'\naccepted options are:\n%s",
				version,
				acceptedVersions,
			)
		}
	}
	return supportedVersion, nil
}

func validateMachineType(machine string, machineTypes *[]provideroptions.MachineType) error {
	if machineTypes == nil {
		return errors.New("received nil machine type list")
	}
	found := false
	accepted := ""
	for _, v := range *machineTypes {
		if v.Name == nil {
			continue
		}
		if *v.Name == machine {
			found = true
			break
		}
		accepted = fmt.Sprintf("%s- %s (cpu: %d, mem: %d)\n", accepted, *v.Name, v.CPU, v.Memory)
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

func validateVolumeType(volume string, volumeTypes *[]provideroptions.VolumeType) error {
	if volumeTypes == nil {
		return errors.New("received nil volune type list")
	}
	found := false
	accepted := ""
	for _, v := range *volumeTypes {
		if v.Name == nil {
			continue
		}
		if *v.Name == volume {
			found = true
			break
		}
		accepted = fmt.Sprintf("%s- %s\n", accepted, *v.Name)
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

func validateZones(zones []string, zoneOptions *[]provideroptions.AvailabilityZone) error {
	if zoneOptions == nil {
		return errors.New("received nil avaiability zones")
	}
	var found bool
	accepted := ""
	for _, v := range *zoneOptions {
		accepted = fmt.Sprintf("%s- %s\n", accepted, *v.Name)
	}
	if len(zones) == 0 {
		return fmt.Errorf(
			"please specify a list of zones\naccepted options are:\n%s",
			accepted,
		)
	}

	for _, v := range zones {
		found = false
		for _, v2 := range *zoneOptions {
			if v2.Name == nil {
				continue
			}
			if v == *v2.Name {
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
