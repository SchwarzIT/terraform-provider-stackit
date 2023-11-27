package loadbalancer

import (
	"context"

	openapiTypes "github.com/SchwarzIT/community-stackit-go-client/pkg/helpers/types"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/load-balancer/1.3.0/instances"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func resToBool(f *bool) basetypes.BoolValue {
	if f == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*f)
}

func resToInt64(f *int) basetypes.Int64Value {
	if f == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*f))
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

func (i *Instance) parse(ctx context.Context, lb instances.LoadBalancer, diags *diag.Diagnostics) {
	i.ID = resToStr(lb.Name)
	i.Name = resToStr(lb.Name)
	i.ExternalAddress = resToStr(lb.ExternalAddress)
	i.PrivateAddress = resToStr(lb.PrivateAddress)
	i.parseOptions(ctx, lb, diags)
	i.parseNetworks(ctx, lb, diags)
	i.parseListeners(ctx, lb, diags)
	i.parseTargetPools(ctx, lb, diags)
}

func (i *Instance) parseOptions(ctx context.Context, lb instances.LoadBalancer, diags *diag.Diagnostics) {
	if lb.Options == nil {
		return
	}
	// Private Network only
	i.PrivateNetworkOnly = resToBool(lb.Options.PrivateNetworkOnly)

	// ACL
	if lb.Options.AccessControl == nil ||
		lb.Options.AccessControl.AllowedSourceRanges == nil {
		return
	}
	ranges := *lb.Options.AccessControl.AllowedSourceRanges
	acl := []attr.Value{}
	for _, r := range ranges {
		acl = append(acl, types.StringValue(r))
	}
	val, d := types.SetValueFrom(ctx, types.StringType, acl)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	i.ACL = val
}

func (i *Instance) parseNetworks(ctx context.Context, lb instances.LoadBalancer, diags *diag.Diagnostics) {
	if lb.Networks == nil {
		return
	}
	networks := []attr.Value{}
	for _, n := range *lb.Networks {
		attrs := map[string]attr.Value{
			"network_id": types.StringNull(),
			"role":       types.StringNull(),
		}
		if n.NetworkID != nil {
			attrs["network_id"] = types.StringValue(n.NetworkID.String())
		}
		if n.Role != nil {
			attrs["role"] = types.StringValue(string(*n.Role))
		}
		networks = append(networks, types.ObjectValueMust(networkType, attrs))
	}
	val, d := types.SetValueFrom(
		ctx,
		types.ObjectType{AttrTypes: networkType},
		networks,
	)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	i.Networks = val
}

func (i *Instance) parseListeners(ctx context.Context, lb instances.LoadBalancer, diags *diag.Diagnostics) {
	if lb.Listeners == nil {
		return
	}
	listeners := []attr.Value{}
	for _, l := range *lb.Listeners {
		attrs := map[string]attr.Value{
			"display_name": resToStr(l.DisplayName),
			"port":         resToInt64(l.Port),
			"protocol":     types.StringNull(),
			"target_pool":  resToStr(l.TargetPool),
		}
		if l.Protocol != nil {
			attrs["protocol"] = types.StringValue(string(*l.Protocol))
		}
		listeners = append(listeners, types.ObjectValueMust(listenerType, attrs))
	}
	val, d := types.SetValueFrom(
		ctx,
		types.ObjectType{AttrTypes: listenerType},
		listeners,
	)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	i.Listeners = val
}

func (i *Instance) parseTargetPools(ctx context.Context, lb instances.LoadBalancer, diags *diag.Diagnostics) {
	if lb.TargetPools == nil {
		return
	}
	targetPools := []attr.Value{}
	for _, tp := range *lb.TargetPools {
		attrs := map[string]attr.Value{
			"name":         resToStr(tp.Name),
			"target_port":  resToInt64(tp.TargetPort),
			"targets":      types.SetNull(targetsType),
			"health_check": types.ObjectNull(healthCheckType),
		}
		if tp.Targets != nil {
			targets := []attr.Value{}
			for _, t := range *tp.Targets {
				targets = append(targets, types.ObjectValueMust(targetType, map[string]attr.Value{
					"display_name": resToStr(t.DisplayName),
					"ip_address":   resToStr(t.Ip),
				}))
			}
			attrs["targets"] = types.SetValueMust(types.ObjectType{AttrTypes: targetType}, targets)
		}
		if tp.ActiveHealthCheck != nil {
			attrs["health_check"] = types.ObjectValueMust(healthCheckType, map[string]attr.Value{
				"healthy_threshold":   resToInt64(tp.ActiveHealthCheck.HealthyThreshold),
				"interval":            resToStr(tp.ActiveHealthCheck.Interval),
				"interval_jitter":     resToStr(tp.ActiveHealthCheck.IntervalJitter),
				"timeout":             resToStr(tp.ActiveHealthCheck.Timeout),
				"unhealthy_threshold": resToInt64(tp.ActiveHealthCheck.UnhealthyThreshold),
			})
		}
		targetPools = append(targetPools, types.ObjectValueMust(targetPoolType, attrs))
	}
	v, d := types.SetValueFrom(
		ctx,
		types.ObjectType{AttrTypes: targetPoolType},
		targetPools,
	)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	i.TargetPools = v
}
