package loadbalancer

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages Secrets Manager instances\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID. Changing this value requires the resource to be recreated.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"external_address": schema.StringAttribute{
				Description: "The external address of the instance.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listeners": schema.SetNestedAttribute{
				Description: "The load balancers listeners.",
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Description: "The port the load balancer listens on.",
							Required:    true,
						},
						"port": schema.Int64Attribute{
							Description: "The port the load balancer listens on [ 1 .. 65535 ].",
							Required:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "The protocol the load balancer listens on. Options: `PROTOCOL_TCP`, `PROTOCOL_UDP`, `PROTOCOL_TCP_PROXY`",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("PROTOCOL_TCP", "PROTOCOL_UDP", "PROTOCOL_TCP_PROXY"),
							},
						},
						"target_pool": schema.StringAttribute{
							Description: "The target pool name.",
							Optional:    true,
						},
					},
				},
			},
			"networks": schema.SetNestedAttribute{
				Description: "The load balancers networks.",
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_id": schema.StringAttribute{
							Description: "The network UUID.",
							Required:    true,
							Validators: []validator.String{
								validate.UUID(),
							},
						},
						"role": schema.StringAttribute{
							Description: "The network role. only `ROLE_LISTENERS_AND_TARGETS` is supported.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("ROLE_LISTENERS_AND_TARGETS"),
							Validators: []validator.String{
								stringvalidator.OneOf("ROLE_LISTENERS_AND_TARGETS"),
							},
						},
					},
				},
			},
			"target_pools": schema.SetNestedAttribute{
				Description: "The load balancers target pools.",
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The target pool name.",
							Required:    true,
						},
						"target_port": schema.Int64Attribute{
							Description: "The target port.",
							Required:    true,
						},
						"targets": schema.SetNestedAttribute{
							Description: "The target pool targets.",
							Required:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"display_name": schema.StringAttribute{
										Description: "The target display name.",
										Required:    true,
									},
									"ip_address": schema.StringAttribute{
										Description: "The target IP address.",
										Required:    true,
									},
								},
							},
						},
						"health_check": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"healthy_threshold": schema.Int64Attribute{
									Description: "The healthy threshold.",
									Required:    true,
								},
								"interval": schema.StringAttribute{
									Description: "The interval.",
									Required:    true,
								},
								"interval_jitter": schema.StringAttribute{
									Description: "The interval jitter.",
									Required:    true,
								},
								"timeout": schema.StringAttribute{
									Description: "The timeout.",
									Required:    true,
								},
								"unhealthy_threshold": schema.Int64Attribute{
									Description: "The unhealthy threshold.",
									Required:    true,
								},
							},
						},
					},
				},
			},
			"acl": schema.SetAttribute{
				Description: "The load balancers ACLs.",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"private_network_only": schema.BoolAttribute{
				Description: "Whether the load balancer is only accessible via private networks.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"private_address": schema.StringAttribute{
				Description: "The private address of the load balancer.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
