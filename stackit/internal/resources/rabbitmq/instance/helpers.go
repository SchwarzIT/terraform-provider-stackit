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
	default_version = "3.7"
	default_plan    = "stackit-rabbitmq-single-small"
)

func (r Resource) validate(ctx context.Context, data *Instance) error {
	if !data.ACL.IsUnknown() && len(data.ACL.Elems) == 0 {
		return errors.New("at least 1 ip address must be specified for `acl`")
	}

	res, err := r.client.DataServices.RabbitMQ.Options.GetOfferings(ctx, data.ProjectID.ValueString())
	if err != nil {
		return errors.Wrap(err, "failed to fetch offerings")
	}

	if err := r.validateVersion(ctx, res.Offerings, data.Version.ValueString()); err != nil {
		return err
	}

	planID, err := r.validatePlan(ctx, res.Offerings, data.Version.ValueString(), data.Plan.ValueString())
	if err != nil {
		return err
	}
	data.PlanID = types.StringValue(planID)
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
			pi.ACL.Elems = append(pi.ACL.Elems, types.StringValue(v))
		}
	} else {
		pi.ACL.Null = true
	}
	pi.Name = types.StringValue(i.Name)
	pi.PlanID = types.StringValue(i.PlanID)
	pi.DashboardURL = types.StringValue(i.DashboardURL)
	pi.CFGUID = types.StringValue(i.CFGUID)
	pi.CFSpaceGUID = types.StringValue(i.CFSpaceGUID)
	pi.CFOrganizationGUID = types.StringValue(i.CFOrganizationGUID)
	return nil
}

func (r Resource) getPlanAndVersion(ctx context.Context, projectID, instanceID string) (plan, version string, err error) {
	dsa := r.client.DataServices.RabbitMQ
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
