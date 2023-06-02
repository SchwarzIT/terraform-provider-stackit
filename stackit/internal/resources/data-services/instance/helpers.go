package instance

import (
	"context"
	"fmt"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0/instances"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/data-services/v1.0/offerings"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkg/errors"
)

func (r Resource) getDefaultVersion() string {
	switch r.service {
	case ElasticSearch:
		return "7"
	case LogMe:
		return "LogMe"
	case MariaDB:
		return "10.4"
	case Opensearch:
		return "2"
	case Postgres:
		return "11"
	case Redis:
		return "6"
	case RabbitMQ:
		return "3.7"
	}
	return ""
}

func (r Resource) getDefaultPlan() string {
	switch r.service {
	case ElasticSearch:
		return "stackit-elasticsearch-single-small"
	case LogMe:
		return "stackit-logme-single-small-non-ssl"
	case MariaDB:
		return "stackit-mariadb-single-small"
	case Opensearch:
		return "stackit-opensearch-1.4.10-single"
	case Postgres:
		return "stackit-postgresql-single-small"
	case Redis:
		return "stackit-redis-single-small"
	case RabbitMQ:
		return "stackit-rabbitmq-single-small"
	}
	return ""
}

func (r Resource) validate(ctx context.Context, data *Instance) error {
	if !data.ACL.IsUnknown() && len(data.ACL.Elements()) == 0 {
		return errors.New("at least 1 ip address must be specified for `acl`")
	}

	res, err := r.client.Offerings.List(ctx, data.ProjectID.ValueString())
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		return agg
	}

	if err := r.validateVersion(ctx, res.JSON200.Offerings, data.Version.ValueString()); err != nil {
		return err
	}

	planID, err := r.validatePlan(ctx, res.JSON200.Offerings, data.Version.ValueString(), data.Plan.ValueString())
	if err != nil {
		return err
	}
	data.PlanID = types.StringValue(planID)
	return nil
}

func (r Resource) validateVersion(ctx context.Context, offers []offerings.Offering, version string) error {
	opts := []string{}
	for _, offer := range offers {
		if offer.Version == version {
			return nil
		}
		opts = append(opts, fmt.Sprintf("- %s (%s)", offer.Version, offer.Name))
	}
	return fmt.Errorf("couldn't find version '%s'. Available options are:\n%s\n", version, strings.Join(opts, "\n"))
}

func (r Resource) validatePlan(ctx context.Context, offers []offerings.Offering, version, planName string) (planID string, err error) {
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
	return "", fmt.Errorf("couldn't find plan name '%s' for version '%s'. Available options are:\n%s\n", planName, version, strings.Join(opts, "\n"))
}

func (r Resource) applyClientResponse(ctx context.Context, pi *Instance, i *instances.Instance) error {
	elems := []attr.Value{}
	if acl, ok := i.Parameters["sgw_acl"]; ok {
		aclString, ok := acl.(string)
		if !ok {
			return errors.New("couldn't parse acl interface as string")
		}
		items := strings.Split(aclString, ",")
		for _, v := range items {
			elems = append(elems, types.StringValue(v))
		}
	}
	pi.ACL = types.ListValueMust(types.StringType, elems)
	pi.Name = types.StringValue(i.Name)
	pi.PlanID = types.StringValue(i.PlanID)
	pi.DashboardURL = types.StringValue(i.DashboardUrl)
	pi.CFGUID = types.StringValue(i.CFGUID)
	pi.CFSpaceGUID = types.StringValue(i.CFSpaceGUID)
	pi.CFOrganizationGUID = types.StringValue("")
	if i.OrganizationGUID != nil {
		pi.CFOrganizationGUID = types.StringValue(*i.OrganizationGUID)
	}
	return nil
}

func (r Resource) getPlanAndVersion(ctx context.Context, projectID, instanceID string) (plan, version string, err error) {
	i, err := r.client.Instances.Get(ctx, projectID, instanceID)
	if agg := validate.Response(i, err, "JSON200"); agg != nil {
		return "", "", agg
	}

	res, err := r.client.Offerings.List(ctx, projectID)
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		return "", "", agg
	}

	for _, offer := range res.JSON200.Offerings {
		for _, p := range offer.Plans {
			if p.ID != i.JSON200.PlanID {
				continue
			}
			return p.Name, offer.Version, nil
		}
	}

	return "", "", errors.Wrapf(err, "couldn't find plan ID %s", i.JSON200.PlanID)
}
