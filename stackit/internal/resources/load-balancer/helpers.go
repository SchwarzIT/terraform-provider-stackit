package loadbalancer

import (
	"context"

	openapiTypes "github.com/SchwarzIT/community-stackit-go-client/pkg/helpers/types"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/load-balancer/1beta.0.0/instances"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func valptr[K string | int | bool](v K) *K {
	return &v
}

func strPtrOrNil(f basetypes.StringValue) *string {
	if f.IsNull() || f.IsUnknown() {
		return nil
	}
	return valptr(f.ValueString())
}

func boolPtrOrNil(f basetypes.BoolValue) *bool {
	if f.IsNull() || f.IsUnknown() {
		return nil
	}
	return valptr(f.ValueBool())
}

func intPtrOrNil(f basetypes.Int64Value) *int {
	if f.IsNull() || f.IsUnknown() {
		return nil
	}
	return valptr(int(f.ValueInt64()))
}

func resToStr(f *string) basetypes.StringValue {
	if f == nil {
		return types.StringNull()
	}
	return types.StringValue(*f)
}

func uuidPtrOrNil(f basetypes.StringValue) *openapiTypes.UUID {
	if f.IsNull() || f.IsUnknown() {
		return nil
	}
	u := uuid.MustParse(f.ValueString())
	utu := openapiTypes.UUID(u)
	return &utu
}

func prepareData(lb Instance) instances.LoadBalancer {

	ilb := instances.LoadBalancer{
		Name:            strPtrOrNil(lb.Name),
		ExternalAddress: strPtrOrNil(lb.ExternalAddress),
		Listeners:       prepareListeners(lb),
		Networks:        prepareNetworks(lb),
		Options:         prepareOptions(lb),
	}

	return ilb
}

func prepareListeners(lb Instance) *[]instances.Listener {
	var listeners []instances.Listener
	for _, l := range lb.Listeners {
		listeners = append(listeners, instances.Listener{
			DisplayName: strPtrOrNil(l.DisplayName),
			Port:        intPtrOrNil(l.Port),
			Protocol:    (*instances.ListenerProtocol)(strPtrOrNil(l.Protocol)),
			TargetPool:  strPtrOrNil(l.TargetPool),
		})
	}
	return &listeners
}

func prepareNetworks(lb Instance) *[]instances.Network {
	var networks []instances.Network
	for _, n := range lb.Networks {
		networks = append(networks, instances.Network{
			NetworkID: uuidPtrOrNil(n.NetworkID),
			Role:      (*instances.NetworkRole)(strPtrOrNil(n.Role)),
		})
	}
	return &networks
}

func prepareTargetPools(lb Instance) *[]instances.TargetPool {
	var targetPools []instances.TargetPool
	for _, tp := range lb.TargetPools {
		targetPools = append(targetPools, instances.TargetPool{
			Name:              strPtrOrNil(tp.Name),
			TargetPort:        intPtrOrNil(tp.TargetPort),
			Targets:           prepareTargets(tp),
			ActiveHealthCheck: prepareHealthCheck(tp),
		})
	}
	return &targetPools
}

func prepareTargets(tp TargetPool) *[]instances.Target {
	var targets []instances.Target
	for _, t := range tp.Targets {
		targets = append(targets, instances.Target{
			DisplayName: strPtrOrNil(t.DisplayName),
			Ip:          strPtrOrNil(t.IPAddress),
		})
	}
	return &targets
}

func prepareHealthCheck(tp TargetPool) *instances.ActiveHealthCheck {
	var hc = tp.HealthCheck
	var healthCheck = instances.ActiveHealthCheck{
		HealthyThreshold:   intPtrOrNil(hc.HealthyThreshold),
		Interval:           strPtrOrNil(hc.Interval),
		IntervalJitter:     strPtrOrNil(hc.IntervalJitter),
		Timeout:            strPtrOrNil(hc.Timeout),
		UnhealthyThreshold: intPtrOrNil(hc.UnhealthyThreshold),
	}
	if deep.Equal(healthCheck, instances.ActiveHealthCheck{}) == nil {
		return nil
	}
	return &healthCheck
}

func prepareOptions(lb Instance) *instances.LoadBalancerOptions {
	opts := instances.LoadBalancerOptions{
		PrivateNetworkOnly: boolPtrOrNil(lb.PrivateNetworkOnly),
		AccessControl:      prepareACL(lb),
	}
	if deep.Equal(opts, instances.LoadBalancerOptions{}) == nil {
		return nil
	}
	return &opts
}

func prepareACL(lb Instance) *instances.LoadbalancerOptionAccessControl {
	if lb.ACL.IsNull() || lb.ACL.IsUnknown() {
		return nil
	}
	var sl []string
	_ = lb.ACL.ElementsAs(context.Background(), &sl, false)
	acl := instances.LoadbalancerOptionAccessControl{
		AllowedSourceRanges: &sl,
	}
	return &acl
}
