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
		TargetPools:     prepareTargetPools(lb),
		Listeners:       prepareListeners(lb),
		Networks:        prepareNetworks(lb),
		Options:         prepareOptions(lb),
	}

	return ilb
}

func prepareListeners(lb Instance) *[]instances.Listener {
	var listeners []instances.Listener
	if lb.Listeners.IsNull() || lb.Listeners.IsUnknown() {
		return nil
	}
	var ls []Listener
	_ = lb.Listeners.ElementsAs(context.Background(), &ls, false)
	for _, l := range ls {
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
	if lb.Networks.IsNull() || lb.Networks.IsUnknown() {
		return nil
	}
	var ns []Network
	_ = lb.Networks.ElementsAs(context.Background(), &ns, false)
	for _, n := range ns {
		networks = append(networks, instances.Network{
			NetworkID: uuidPtrOrNil(n.NetworkID),
			Role:      (*instances.NetworkRole)(strPtrOrNil(n.Role)),
		})
	}
	return &networks
}

func prepareTargetPools(lb Instance) *[]instances.TargetPool {
	var targetPools []instances.TargetPool
	if lb.TargetPools.IsNull() || lb.TargetPools.IsUnknown() {
		return nil
	}
	var tp []TargetPool
	_ = lb.TargetPools.ElementsAs(context.Background(), &tp, false)
	for _, tp := range tp {
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
	if tp.Targets.IsNull() || tp.Targets.IsUnknown() {
		return nil
	}
	var ts []Target
	_ = tp.Targets.ElementsAs(context.Background(), &ts, false)
	for _, t := range ts {
		targets = append(targets, instances.Target{
			DisplayName: strPtrOrNil(t.DisplayName),
			Ip:          strPtrOrNil(t.IPAddress),
		})
	}
	return &targets
}

func prepareHealthCheck(tp TargetPool) *instances.ActiveHealthCheck {
	var hc HealthCheck
	if tp.HealthCheck.IsNull() || tp.HealthCheck.IsUnknown() {
		return nil
	}
	_ = tp.HealthCheck.As(context.Background(), &hc, basetypes.ObjectAsOptions{})
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
