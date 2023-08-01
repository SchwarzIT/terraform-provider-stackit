package postgresinstance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/instance"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/postgres-flex/v1.0/versions"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	DefaultVersion              = "14"
	DefaultReplicas       int64 = 1
	DefaultBackupSchedule       = "0 2 * * *"
	DefaultStorageClass         = "premium-perf6-stackit"
	DefaultStorageSize    int64 = 5
	DefaultMachineType          = "2.4"
)

func (i *Instance) setDefaults() {
	if i.Version.IsNull() || i.Version.IsUnknown() {
		i.Version = types.StringValue(DefaultVersion)
	}

	if i.Replicas.IsNull() || i.Replicas.IsUnknown() {
		i.Replicas = types.Int64Value(DefaultReplicas)
	}

	if i.BackupSchedule.IsNull() || i.BackupSchedule.IsUnknown() {
		i.BackupSchedule = types.StringValue(DefaultBackupSchedule)
	}
}

func (r Resource) validate(ctx context.Context, diags *diag.Diagnostics, data Instance) error {
	if err := r.validateVersion(ctx, diags, data.ProjectID.ValueString(), data.Version.ValueString()); err != nil {
		return err
	}
	if err := r.validateMachineType(ctx, diags, data.ProjectID.ValueString(), data.MachineType.ValueString()); err != nil {
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

	if err := r.validateStorage(ctx, diags, data.ProjectID.ValueString(), data.MachineType.ValueString(), storage); err != nil {
		return err
	}
	return nil
}

func (r Resource) validateVersion(ctx context.Context, diags *diag.Diagnostics, projectID, version string) error {
	res, err := r.client.PostgresFlex.Versions.List(ctx, projectID, &versions.ListParams{})
	if agg := common.Validate(diags, res, err, "JSON200.Versions"); agg != nil {
		return agg
	}

	opts := ""
	for _, v := range *res.JSON200.Versions {
		opts = opts + "\n- " + v
		if strings.EqualFold(v, version) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", version, opts)
}

func (r Resource) validateMachineType(ctx context.Context, diags *diag.Diagnostics, projectID, flavorID string) error {
	res, err := r.client.PostgresFlex.Flavors.List(ctx, projectID)
	if agg := common.Validate(diags, res, err, "JSON200.Flavors"); agg != nil {
		return agg
	}

	opts := ""
	for _, v := range *res.JSON200.Flavors {
		if v.ID == nil || v.Cpu == nil || v.Memory == nil {
			continue
		}
		opts = fmt.Sprintf("%s\n - ID: %s (CPU: %d, Mem: %d)", opts, *v.ID, *v.Cpu, *v.Memory)
		if strings.EqualFold(*v.ID, flavorID) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find machine type '%s'. Available options are:%s\n", flavorID, opts)
}

func (r Resource) validateStorage(ctx context.Context, diags *diag.Diagnostics, projectID, machineType string, storage Storage) error {
	res, err := r.client.PostgresFlex.Storage.GetStorageOptions(ctx, projectID, machineType)
	if agg := common.Validate(diags, res, err, "JSON200.StorageClasses"); agg != nil {
		return agg
	}

	size := storage.Size.ValueInt64()
	if res.JSON200.StorageRange != nil && res.JSON200.StorageRange.Max != nil && res.JSON200.StorageRange.Min != nil {
		if int64(*res.JSON200.StorageRange.Max) < size || int64(*res.JSON200.StorageRange.Min) > size {
			return fmt.Errorf("storage size %d is not in the allowed range: %d..%d", size, *res.JSON200.StorageRange.Min, *res.JSON200.StorageRange.Max)
		}
	}

	opts := ""
	for _, v := range *res.JSON200.StorageClasses {
		opts = opts + "\n- " + v
		if strings.EqualFold(v, storage.Class.ValueString()) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", storage.Class.ValueString(), opts)
}

func applyClientResponse(pi *Instance, i *instance.InstanceSingleInstance) error {
	elems := []attr.Value{}
	if i.ACL != nil && i.ACL.Items != nil {
		for _, v := range *i.ACL.Items {
			// only include correctly formatted CIDR range
			// this is to overcome a current bug in the API
			if strings.Contains(v, "/") {
				elems = append(elems, types.StringValue(v))
			}
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
