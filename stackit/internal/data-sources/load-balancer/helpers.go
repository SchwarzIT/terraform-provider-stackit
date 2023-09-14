package loadbalancer

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/load-balancer/1beta.0.0/instances"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

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
		lb.Options = &instances.LoadBalancerOptions{}
	}
	// Private Network only
	i.PrivateNetworkOnly = resToBool(lb.Options.PrivateNetworkOnly)

	// ACL
	if lb.Options.AccessControl == nil ||
		lb.Options.AccessControl.AllowedSourceRanges == nil {
		i.ACL = types.SetNull(types.StringType)
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
		i.Networks = types.SetNull(types.ObjectType{AttrTypes: networkType})
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
		i.Listeners = types.SetNull(types.ObjectType{AttrTypes: listenerType})
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
		i.TargetPools = types.SetNull(types.ObjectType{AttrTypes: targetPoolType})
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
