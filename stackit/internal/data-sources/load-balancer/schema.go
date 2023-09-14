package loadbalancer

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProjectID          types.String `tfsdk:"project_id"`
	ExternalAddress    types.String `tfsdk:"external_address"`
	Listeners          types.Set    `tfsdk:"listeners"`
	Networks           types.Set    `tfsdk:"networks"`
	TargetPools        types.Set    `tfsdk:"target_pools"`
	ACL                types.Set    `tfsdk:"acl"`
	PrivateNetworkOnly types.Bool   `tfsdk:"private_network_only"`
	PrivateAddress     types.String `tfsdk:"private_address"`
}

type Listener struct {
	DisplayName types.String `tfsdk:"display_name"`
	Port        types.Int64  `tfsdk:"port"`
	Protocol    types.String `tfsdk:"protocol"`
	TargetPool  types.String `tfsdk:"target_pool"`
}

var listenerType = map[string]attr.Type{
	"display_name": types.StringType,
	"port":         types.Int64Type,
	"protocol":     types.StringType,
	"target_pool":  types.StringType,
}

type Network struct {
	NetworkID types.String `tfsdk:"network_id"`
	Role      types.String `tfsdk:"role"`
}

var networkType = map[string]attr.Type{
	"network_id": types.StringType,
	"role":       types.StringType,
}

type TargetPool struct {
	Name        types.String `tfsdk:"name"`
	TargetPort  types.Int64  `tfsdk:"target_port"`
	Targets     types.Set    `tfsdk:"targets"`
	HealthCheck types.Object `tfsdk:"health_check"`
}

var targetPoolType = map[string]attr.Type{
	"name":        types.StringType,
	"target_port": types.Int64Type,
	"targets":     targetsType,
	"health_check": types.ObjectType{
		AttrTypes: healthCheckType,
	},
}

type Target struct {
	DisplayName types.String `tfsdk:"display_name"`
	IPAddress   types.String `tfsdk:"ip_address"`
}

var targetsType = types.SetType{
	ElemType: types.ObjectType{
		AttrTypes: targetType,
	},
}

var targetType = map[string]attr.Type{
	"display_name": types.StringType,
	"ip_address":   types.StringType,
}

type HealthCheck struct {
	HealthyThreshold   types.Int64  `tfsdk:"healthy_threshold"`
	Interval           types.String `tfsdk:"interval"`
	IntervalJitter     types.String `tfsdk:"interval_jitter"`
	Timeout            types.String `tfsdk:"timeout"`
	UnhealthyThreshold types.Int64  `tfsdk:"unhealthy_threshold"`
}

var healthCheckType = map[string]attr.Type{
	"healthy_threshold":   types.Int64Type,
	"interval":            types.StringType,
	"interval_jitter":     types.StringType,
	"timeout":             types.StringType,
	"unhealthy_threshold": types.Int64Type,
}

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for Load Balancer instances\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID. Changing this value requires the resource to be recreated.",
				Required:    true,
			},
			"external_address": schema.StringAttribute{
				Description: "The external address of the instance.",
				Computed:    true,
			},
			"listeners": schema.SetNestedAttribute{
				Description: "The load balancers listeners.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Description: "The port the load balancer listens on.",
							Computed:    true,
						},
						"port": schema.Int64Attribute{
							Description: "The port the load balancer listens on [ 1 .. 65535 ].",
							Computed:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "The protocol the load balancer listens on. Options: `PROTOCOL_TCP`, `PROTOCOL_UDP`, `PROTOCOL_TCP_PROXY`",
							Computed:    true,
						},
						"target_pool": schema.StringAttribute{
							Description: "The target pool name.",
							Computed:    true,
						},
					},
				},
			},
			"networks": schema.SetNestedAttribute{
				Description: "The load balancers networks.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_id": schema.StringAttribute{
							Description: "The network UUID.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "The network role. only `ROLE_LISTENERS_AND_TARGETS` is supported.",
							Computed:    true,
						},
					},
				},
			},
			"target_pools": schema.SetNestedAttribute{
				Description: "The load balancers target pools.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The target pool name.",
							Computed:    true,
						},
						"target_port": schema.Int64Attribute{
							Description: "The target port.",
							Computed:    true,
						},
						"targets": schema.SetNestedAttribute{
							Description: "The target pool targets.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"display_name": schema.StringAttribute{
										Description: "The target display name.",
										Computed:    true,
									},
									"ip_address": schema.StringAttribute{
										Description: "The target IP address.",
										Computed:    true,
									},
								},
							},
						},
						"health_check": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"healthy_threshold": schema.Int64Attribute{
									Description: "The healthy threshold.",
									Computed:    true,
								},
								"interval": schema.StringAttribute{
									Description: "The interval.",
									Computed:    true,
								},
								"interval_jitter": schema.StringAttribute{
									Description: "The interval jitter.",
									Computed:    true,
								},
								"timeout": schema.StringAttribute{
									Description: "The timeout.",
									Computed:    true,
								},
								"unhealthy_threshold": schema.Int64Attribute{
									Description: "The unhealthy threshold.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"acl": schema.SetAttribute{
				Description: "The load balancers ACLs.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"private_network_only": schema.BoolAttribute{
				Description: "Whether the load balancer is only accessible via private networks.",
				Computed:    true,
			},
			"private_address": schema.StringAttribute{
				Description: "The private address of the load balancer.",
				Computed:    true,
			},
		},
	}
}
