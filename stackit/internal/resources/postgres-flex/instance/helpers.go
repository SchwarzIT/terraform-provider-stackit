package postgresinstance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/postgres-flex/instances"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	default_username              = "stackit"
	default_backup_schedule       = "0 2 * * *"
	default_storage_class         = "premium-perf6-stackit"
	default_storage_size    int64 = 20
)

func (r Resource) validate(ctx context.Context, data PostgresInstance) error {
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

	if err := r.validateStorageClass(ctx, data.ProjectID.Value, data.MachineType.Value, storage.Class.Value); err != nil {
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
	for _, v := range res {
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

func (r Resource) validateStorageClass(ctx context.Context, projectID, machineType, storageClass string) error {
	res, err := r.client.PostgresFlex.Options.GetStorageClasses(ctx, projectID, machineType)
	if err != nil {
		return err
	}

	opts := ""
	for _, v := range res.StorageClasses {
		opts = opts + "\n- " + v
		if strings.ToLower(v) == strings.ToLower(storageClass) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", storageClass, opts)
}

func applyClientResponse(pi *PostgresInstance, i instances.Instance) error {
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
	pi.Version = types.String{Value: i.Version}
	return nil
}
