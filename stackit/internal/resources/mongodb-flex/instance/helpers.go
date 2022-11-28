package instance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/mongodb-flex/instances"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	default_version               = "6.0"
	default_replicas        int64 = 1
	default_username              = "stackit"
	default_backup_schedule       = "0 2 * * *"
	default_storage_class         = "premium-perf2-mongodb"
	default_storage_size    int64 = 10
)

func (i *Instance) setDefaults() {
	if i.Version.IsNull() || i.Version.IsUnknown() {
		i.Version = types.String{Value: default_version}
	}

	if i.Replicas.IsNull() || i.Replicas.IsUnknown() {
		i.Replicas = types.Int64{Value: default_replicas}
	}

	if i.BackupSchedule.IsNull() || i.BackupSchedule.IsUnknown() {
		i.BackupSchedule = types.String{Value: default_backup_schedule}
	}
}

func (r Resource) validate(ctx context.Context, data Instance) error {
	if err := r.validateVersion(ctx, data.ProjectID.Value, data.Version.Value); err != nil {
		return err
	}
	if err := r.validateMachineType(ctx, data.ProjectID.Value, data.MachineType.Value); err != nil {
		return err
	}

	if data.Storage.IsNull() || data.Storage.IsUnknown() {
		return nil
	}

	storage := Storage{}
	diag := data.Storage.As(ctx, &storage, types.ObjectAsOptions{})
	if diag.HasError() {
		return errors.New("failed setting storage from object")
	}

	if err := r.validateStorage(ctx, data.ProjectID.Value, data.MachineType.Value, storage); err != nil {
		return err
	}
	return nil
}

func (r Resource) validateVersion(ctx context.Context, projectID, version string) error {
	res, err := r.client.MongoDBFlex.Options.GetVersions(ctx, projectID)
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
	res, err := r.client.MongoDBFlex.Options.GetFlavors(ctx, projectID)
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
	res, err := r.client.MongoDBFlex.Options.GetStorageClasses(ctx, projectID, machineType)
	if err != nil {
		return err
	}

	size := storage.Size.Value
	if int64(res.StorageRange.Max) < size || int64(res.StorageRange.Min) > size {
		return fmt.Errorf("storage size %d is not in the allowed range: %d..%d", size, res.StorageRange.Min, res.StorageRange.Max)
	}

	opts := ""
	for _, v := range res.StorageClasses {
		opts = opts + "\n- " + v
		if strings.ToLower(v) == strings.ToLower(storage.Class.Value) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", storage.Class.Value, opts)
}

func applyClientResponse(pi *Instance, i instances.Instance) error {
	pi.ACL = types.List{ElemType: types.StringType}
	for _, v := range i.ACL.Items {
		pi.ACL.Elems = append(pi.ACL.Elems, types.String{Value: v})
	}
	pi.BackupSchedule = types.String{Value: i.BackupSchedule}
	pi.MachineType = types.String{Value: i.Flavor.ID}

	pi.Name = types.String{Value: i.Name}
	pi.Replicas = types.Int64{Value: int64(i.Replicas)}

	storage, diags := types.ObjectValue(
		map[string]attr.Type{
			"class": types.StringType,
			"size":  types.Int64Type,
		},
		map[string]attr.Value{
			"class": types.String{Value: i.Storage.Class},
			"size":  types.Int64{Value: int64(i.Storage.Size)},
		})
	if diags.HasError() {
		return errors.New("failed setting storage object")
	}
	pi.Storage = storage
	if len(i.Version) > 3 {
		i.Version = i.Version[0:3]
	}
	pi.Version = types.String{Value: i.Version}
	return nil
}
