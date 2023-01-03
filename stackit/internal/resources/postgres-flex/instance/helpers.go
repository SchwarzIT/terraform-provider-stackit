package postgresinstance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/generated/instance"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/generated/versions"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	default_version               = "14"
	default_replicas        int64 = 1
	default_username              = "stackit"
	default_backup_schedule       = "0 2 * * *"
	default_storage_class         = "premium-perf6-stackit"
	default_storage_size    int64 = 20
)

func (i *Instance) setDefaults() {
	if i.Version.IsNull() || i.Version.IsUnknown() {
		i.Version = types.StringValue(default_version)
	}

	if i.Replicas.IsNull() || i.Replicas.IsUnknown() {
		i.Replicas = types.Int64Value(default_replicas)
	}

	if i.BackupSchedule.IsNull() || i.BackupSchedule.IsUnknown() {
		i.BackupSchedule = types.StringValue(default_backup_schedule)
	}
}

func (r Resource) validate(ctx context.Context, data Instance) error {
	if err := r.validateVersion(ctx, data.ProjectID.ValueString(), data.Version.ValueString()); err != nil {
		return err
	}
	if err := r.validateMachineType(ctx, data.ProjectID.ValueString(), data.MachineType.ValueString()); err != nil {
		return err
	}

	if data.Storage.IsNull() || data.Storage.IsUnknown() {
		return nil
	}

	storage := Storage{}
	diag := data.Storage.As(ctx, &storage, basetypes.ObjectAsOptions{})
	if diag.HasError() {
		return errors.New("failed setting storage from object")
	}

	if err := r.validateStorage(ctx, data.ProjectID.ValueString(), data.MachineType.ValueString(), storage); err != nil {
		return err
	}
	return nil
}

func (r Resource) validateVersion(ctx context.Context, projectID, version string) error {
	res, err := r.client.PostgresFlex.Versions.GetVersionsWithResponse(ctx, projectID, &versions.GetVersionsParams{})
	if err != nil {
		return err
	}
	if res.HasError != nil {
		return res.HasError
	}
	if res.JSON200 == nil {
		return errors.New("received an empty response for versions")
	}

	opts := ""
	for _, v := range *res.JSON200.Versions {
		opts = opts + "\n- " + v
		if strings.ToLower(v) == strings.ToLower(version) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", version, opts)
}

func (r Resource) validateMachineType(ctx context.Context, projectID, flavorID string) error {
	res, err := r.client.PostgresFlex.Flavors.GetFlavorsWithResponse(ctx, projectID)
	if err != nil {
		return err
	}
	if res.HasError != nil {
		return res.HasError
	}
	if res.JSON200 == nil {
		return errors.New("received an empty response for versions")
	}

	opts := ""
	for _, v := range *res.JSON200.Flavors {
		if v.ID == nil || v.Cpu == nil || v.Memory == nil {
			continue
		}
		opts = fmt.Sprintf("%s\n - ID: %s (CPU: %d, Mem: %d)", opts, *v.ID, *v.Cpu, *v.Memory)
		if strings.ToLower(*v.ID) == strings.ToLower(flavorID) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find machine type '%s'. Available options are:%s\n", flavorID, opts)
}

func (r Resource) validateStorage(ctx context.Context, projectID, machineType string, storage Storage) error {
	res, err := r.client.PostgresFlex.Storage.GetFlavorWithResponse(ctx, projectID, machineType)
	if err != nil {
		return err
	}
	if res.HasError != nil {
		return res.HasError
	}
	if res.JSON200 == nil {
		return errors.New("received an empty response for versions")
	}

	size := storage.Size.ValueInt64()
	if res.JSON200.StorageRange != nil && res.JSON200.StorageRange.Max != nil && res.JSON200.StorageRange.Min != nil {
		if int64(*res.JSON200.StorageRange.Max) < size || int64(*res.JSON200.StorageRange.Min) > size {
			return fmt.Errorf("storage size %d is not in the allowed range: %d..%d", size, *res.JSON200.StorageRange.Min, *res.JSON200.StorageRange.Max)
		}
	}

	opts := ""
	if res.JSON200.StorageClasses != nil {
		for _, v := range *res.JSON200.StorageClasses {
			opts = opts + "\n- " + v
			if strings.ToLower(v) == strings.ToLower(storage.Class.ValueString()) {
				return nil
			}
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", storage.Class.ValueString(), opts)
}

func applyClientResponse(pi *Instance, i *instance.InstanceSingleInstance) error {
	elems := []attr.Value{}
	if i.ACL != nil && i.ACL.Items != nil {
		for _, v := range *i.ACL.Items {
			elems = append(elems, types.StringValue(v))
		}
	}
	pi.ACL = types.ListValueMust(types.StringType, elems)

	pi.BackupSchedule = types.StringNull()
	if i.BackupSchedule != nil {
		pi.BackupSchedule = types.StringValue(*i.BackupSchedule)
	}
	pi.MachineType = types.StringNull()
	if i.Flavor != nil && i.Flavor.ID != nil {
		pi.MachineType = types.StringValue(*i.Flavor.ID)
	}
	pi.Name = types.StringNull()
	if i.Name != nil {
		pi.Name = types.StringValue(*i.Name)
	}
	pi.Replicas = types.Int64Null()
	if i.Replicas != nil {
		pi.Replicas = types.Int64Value(int64(*i.Replicas))
	}
	if i.Storage != nil {
		class := types.StringNull()
		if i.Storage.Class != nil {
			class = types.StringValue(*i.Storage.Class)
		}
		size := types.Int64Null()
		if i.Storage.Class != nil {
			size = types.Int64Value(int64(*i.Storage.Size))
		}
		storage, diags := types.ObjectValue(
			map[string]attr.Type{
				"class": types.StringType,
				"size":  types.Int64Type,
			},
			map[string]attr.Value{
				"class": class,
				"size":  size,
			})
		if diags.HasError() {
			return errors.New("failed setting storage object")
		}
		pi.Storage = storage
	}
	pi.Version = types.StringNull()
	if i.Version != nil {
		pi.Version = types.StringValue(*i.Version)
	}
	return nil
}
