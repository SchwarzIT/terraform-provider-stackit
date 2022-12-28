package postgresinstance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/postgres-flex/instances"
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
	res, err := r.client.PostgresFlex.Options.GetVersions(ctx, projectID)
	if err != nil {
		return err
	}

	opts := ""
	for _, v := range res.Versions {
		opts = opts + "\n- " + v
		if strings.ToLower(v) == strings.ToLower(version) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", version, opts)
}

func (r Resource) validateMachineType(ctx context.Context, projectID, flavorID string) error {
	res, err := r.client.PostgresFlex.Options.GetFlavors(ctx, projectID)
	if err != nil {
		return err
	}

	opts := ""
	for _, v := range res.Flavors {
		opts = fmt.Sprintf("%s\n - ID: %s (CPU: %d, Mem: %d)", opts, v.ID, v.CPU, v.Memory)
		if strings.ToLower(v.ID) == strings.ToLower(flavorID) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find machine type '%s'. Available options are:%s\n", flavorID, opts)
}

func (r Resource) validateStorage(ctx context.Context, projectID, machineType string, storage Storage) error {
	res, err := r.client.PostgresFlex.Options.GetStorageClasses(ctx, projectID, machineType)
	if err != nil {
		return err
	}

	size := storage.Size.ValueInt64()
	if int64(res.StorageRange.Max) < size || int64(res.StorageRange.Min) > size {
		return fmt.Errorf("storage size %d is not in the allowed range: %d..%d", size, res.StorageRange.Min, res.StorageRange.Max)
	}

	opts := ""
	for _, v := range res.StorageClasses {
		opts = opts + "\n- " + v
		if strings.ToLower(v) == strings.ToLower(storage.Class.ValueString()) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", storage.Class.ValueString(), opts)
}

func applyClientResponse(pi *Instance, i instances.Instance) error {
	elems := []attr.Value{}
	for _, v := range i.ACL.Items {
		elems = append(elems, types.StringValue(v))
	}
	pi.ACL = types.ListValueMust(types.StringType, elems)
	pi.BackupSchedule = types.StringValue(i.BackupSchedule)
	pi.MachineType = types.StringValue(i.Flavor.ID)
	pi.Name = types.StringValue(i.Name)
	pi.Replicas = types.Int64Value(int64(i.Replicas))
	storage, diags := types.ObjectValue(
		map[string]attr.Type{
			"class": types.StringType,
			"size":  types.Int64Type,
		},
		map[string]attr.Value{
			"class": types.StringValue(i.Storage.Class),
			"size":  types.Int64Value(int64(i.Storage.Size)),
		})
	if diags.HasError() {
		return errors.New("failed setting storage object")
	}
	pi.Storage = storage
	pi.Version = types.StringValue(i.Version)
	return nil
}
