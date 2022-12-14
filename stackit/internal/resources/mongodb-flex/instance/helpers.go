package instance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/mongodb-flex/v1.0/generated/instance"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	default_version               = "6.0"
	default_replicas        int64 = 1
	default_username              = "stackit"
	default_backup_schedule       = "0 0/24 * * *"
	default_storage_class         = "premium-perf2-mongodb"
	default_storage_size    int64 = 10
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
	res, err := r.client.MongoDBFlex.Versions.GetVersionsWithResponse(ctx, projectID)
	if err != nil {
		return err
	}
	if res.HasError != nil {
		return res.HasError
	}
	if res.JSON200 == nil || res.JSON200.Versions == nil {

		return fmt.Errorf("JSON200 == nil or Versions == nil.\nRequest URL: %s\nStatus: %s\nBody: %s", res.HTTPResponse.Request.URL, res.Status(), string(res.Body))
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
	res, err := r.client.MongoDBFlex.Flavors.GetFlavorsWithResponse(ctx, projectID)
	if err != nil {
		return err
	}
	if res.HasError != nil {
		return res.HasError
	}
	if res.JSON200 == nil || res.JSON200.Flavors == nil {
		return errors.New("JSON200 == nil or Flavors == nil")
	}

	opts := ""
	for _, v := range *res.JSON200.Flavors {
		if v.ID == nil {
			continue
		}
		opts = fmt.Sprintf("%s\n - ID: %s (CPU: %d, Mem: %d)", opts, *v.ID, v.CPU, v.Memory)
		if strings.ToLower(*v.ID) == strings.ToLower(flavorID) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find machine type '%s'. Available options are:%s\n", flavorID, opts)
}

func (r Resource) validateStorage(ctx context.Context, projectID, machineType string, storage Storage) error {
	res, err := r.client.MongoDBFlex.Storage.GetStoragesFlavorWithResponse(ctx, projectID, machineType)
	if err != nil {
		return err
	}
	if res.HasError != nil {
		return res.HasError
	}
	if res.JSON200 == nil || res.JSON200.StorageRange == nil {
		return errors.New("JSON200 == nil or Flavors == nil")
	}

	size := storage.Size.ValueInt64()
	if res.JSON200.StorageRange.Max != nil && res.JSON200.StorageRange.Min != nil {
		if int64(*res.JSON200.StorageRange.Max) < size || int64(*res.JSON200.StorageRange.Min) > size {
			return fmt.Errorf("storage size %d is not in the allowed range: %d..%d", size, *res.JSON200.StorageRange.Min, *res.JSON200.StorageRange.Max)
		}
	}

	opts := ""
	for _, v := range *res.JSON200.StorageClasses {
		opts = opts + "\n- " + v
		if strings.ToLower(v) == strings.ToLower(storage.Class.ValueString()) {
			return nil
		}
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:%s\n", storage.Class.ValueString(), opts)
}

func applyClientResponse(pi *Instance, i *instance.InstancesSingleInstance) error {
	elems := []attr.Value{}
	if i == nil {
		return errors.New("instance response is empty")
	}
	if i.ACL != nil && i.ACL.Items != nil {
		for _, v := range *i.ACL.Items {
			elems = append(elems, types.StringValue(v))
		}
	}
	pi.ACL = types.ListValueMust(types.StringType, elems)
	pi.BackupSchedule = nullOrValStr(i.BackupSchedule)
	pi.MachineType = nullOrValStr(i.Flavor.ID)
	pi.Name = nullOrValStr(i.Name)
	pi.Replicas = nullOrValInt64(i.Replicas)
	storage, diags := types.ObjectValue(
		map[string]attr.Type{
			"class": types.StringType,
			"size":  types.Int64Type,
		},
		map[string]attr.Value{
			"class": nullOrValStr(i.Storage.Class),
			"size":  nullOrValInt64(i.Storage.Size),
		})
	if diags.HasError() {
		return errors.New("failed setting storage object")
	}
	pi.Storage = storage
	pi.Version = nullOrValStr(i.Version)
	if !pi.Version.IsNull() && len(pi.Version.ValueString()) > 3 {
		v := pi.Version.ValueString()
		pi.Version = types.StringValue(v[0:3])
	}
	return nil
}

func nullOrValStr(v *string) basetypes.StringValue {
	a := types.StringNull()
	if v != nil {
		a = types.StringValue(*v)
	}
	return a
}

func nullOrValInt64(v *int) basetypes.Int64Value {
	a := types.Int64Null()
	if v != nil {
		a = types.Int64Value(int64(*v))
	}
	return a
}
