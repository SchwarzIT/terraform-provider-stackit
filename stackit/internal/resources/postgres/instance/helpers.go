package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/instances"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/options"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkg/errors"
)

const (
	default_version = "11"
	default_plan    = "stackit-postgresql-single-small"
)

func (r Resource) validate(ctx context.Context, data *Instance) error {
	if data.Version.IsNull() || data.Version.IsUnknown() {
		data.Version = types.String{Value: default_version}
	}
	if data.Plan.IsNull() || data.Plan.IsUnknown() {
		data.Plan = types.String{Value: default_plan}
	}
	if !data.ACL.IsUnknown() && len(data.ACL.Elems) == 0 {
		return errors.New("at least 1 ip address must be specified for `acl`")
	}

	res, err := r.client.DataServices.PostgresDB.Options.GetOfferings(ctx, data.ProjectID.Value)
	if err != nil {
		return errors.Wrap(err, "failed to fetch offerings")
	}

	if err := r.validateVersion(ctx, res.Offerings, data.Version.Value); err != nil {
		return err
	}

	planID, err := r.validatePlan(ctx, res.Offerings, data.Version.Value, data.Plan.Value)
	if err != nil {
		return err
	}

	data.PlanID = types.String{Value: planID}
	return nil
}

func (r Resource) validateVersion(ctx context.Context, offers []options.Offer, version string) error {
	opts := []string{}
	for _, offer := range offers {
		if offer.Version == version {
			return nil
		}
		opts = append(opts, fmt.Sprintf("- %s (%s)", offer.Version, offer.Name))
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:\n%s\n", version, strings.Join(opts, "\n"))
}

func (r Resource) validatePlan(ctx context.Context, offers []options.Offer, version, planName string) (planID string, err error) {
	opts := []string{}
	offerIndex := 0
	for i, offer := range offers {
		if offer.Version == version {
			offerIndex = i
			break
		}
	}
	for _, plan := range offers[offerIndex].Plans {
		if plan.Name == planName {
			return plan.ID, nil
		}
		opts = append(opts, fmt.Sprintf("- %s (%s)", plan.Name, plan.Description))
	}
	return "", fmt.Errorf("couldn't find plan name '%s'. Available options are:\n%s\n", version, strings.Join(opts, "\n"))
}

func (r Resource) applyClientResponse(ctx context.Context, pi *Instance, i instances.Instance) error {
	pi.ACL = types.List{ElemType: types.StringType}
	if aclString, ok := i.Parameters["sgw_acl"]; ok {
		items := strings.Split(aclString, ",")
		for _, v := range items {
			pi.ACL.Elems = append(pi.ACL.Elems, types.String{Value: v})
		}
	} else {
		pi.ACL.Null = true
	}

	pi.Name = types.String{Value: i.Name}
	pi.PlanID = types.String{Value: i.PlanID}
	pi.DashboardURL = types.String{Value: i.DashboardURL}
	pi.CFGUID = types.String{Value: i.CFGUID}
	pi.CFSpaceGUID = types.String{Value: i.CFSpaceGUID}
	pi.CFOrganizationGUID = types.String{Value: i.CFOrganizationGUID}
	return nil
}

func (r Resource) getPlanAndVersion(ctx context.Context, projectID, instanceID string) (plan, version string, err error) {
	dsa := r.client.DataServices.PostgresDB
	i, err := dsa.Instances.Get(ctx, projectID, instanceID)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to fetch instance")
	}

	res, err := dsa.Options.GetOfferings(ctx, projectID)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to fetch offerings")
	}

	for _, offer := range res.Offerings {
		for _, p := range offer.Plans {
			if p.ID != i.PlanID {
				continue
			}
			return p.Name, offer.Version, nil
		}
	}

	return "", "", errors.Wrapf(err, "couldn't find plan ID %s", i.PlanID)
}
