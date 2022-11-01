package postgresinstance

import (
	"context"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/postgres/instances"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r Resource) validate(ctx context.Context, data PostgresInstance) error {
	if err := r.validateVersion(ctx, data.ProjectID.Value, data.Version.Value); err != nil {
		return err
	}
	if err := r.validateMachineType(ctx, data.ProjectID.Value, data.MachineType.Value); err != nil {
		return err
	}

	if data.Storage.Class.IsNull() || data.Storage.Class.IsUnknown() {
		return nil
	}
	if err := r.validateStorageClass(ctx, data.ProjectID.Value, data.MachineType.Value, data.Storage.Class.Value); err != nil {
		return err
	}
	return nil
}

func (r Resource) validateVersion(ctx context.Context, projectID, version string) error {
	res, err := r.client.Incubator.Postgres.Options.GetVersions(ctx, projectID)
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
	res, err := r.client.Incubator.Postgres.Options.GetFlavors(ctx, projectID)
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
	res, err := r.client.Incubator.Postgres.Options.GetStorageClasses(ctx, projectID, machineType)
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

func (pi *PostgresInstance) ApplyClientResponse(i instances.Instance) {
	pi.ACL = types.List{ElemType: types.StringType}
	els := []attr.Value{}
	for _, v := range i.ACL.Items {
		els = append(els, types.String{Value: v})
	}
	pi.ACL.Elems = els
	pi.BackupSchedule = types.String{Value: i.BackupSchedule}
	pi.MachineType = types.String{Value: i.Flavor.ID}
	pi.Name = types.String{Value: i.Name}
	pi.Replicas = types.Int64{Value: int64(i.Replicas)}
	pi.Storage = Storage{
		Class: types.String{Value: i.Storage.Class},
		Size:  types.Int64{Value: int64(i.Storage.Size)},
	}
	pi.Version = types.String{Value: i.Version}

	if len(i.Users) == 0 {
		pi.Users = nil
		return
	}

	pi.Users = []User{}
	for _, user := range i.Users {
		pi.Users = append(pi.Users, User{
			ID:       types.String{Value: user.ID},
			Username: types.String{Value: user.Username},
			Password: types.String{Value: user.Password},
			Hostname: types.String{Value: user.Hostname},
			Database: types.String{Value: user.Database},
			Port:     types.Int64{Value: int64(user.Port)},
			URI:      types.String{Value: user.URI},
			Roles:    user.Roles,
		})
	}
}
