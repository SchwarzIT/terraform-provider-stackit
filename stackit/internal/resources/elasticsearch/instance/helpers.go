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
	default_version = "7"
	default_plan    = "0 2 * * *"
)

func (r Resource) validate(ctx context.Context, data *Instance) error {
	res, err := r.client.DataServices.ElasticSearch.Options.GetOfferings(ctx, data.ProjectID.Value)
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

func applyClientResponse(pi *Instance, i instances.Instance) error {
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
	return nil
}
