package loadbalancer

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProjectID          types.String `tfsdk:"project_id"`
	ExternalIP         types.String `tfsdk:"external_ip"`
	Listeners          types.Set    `tfsdk:"listeners"`
	Networks           types.Set    `tfsdk:"networks"`
	TargetPools        types.Set    `tfsdk:"target_pools"`
	ACL                types.Set    `tfsdk:"acl"`
	PrivateNetworkOnly types.Bool   `tfsdk:"private_network_only"`
}

type Listener struct {
	DisplayName types.String `tfsdk:"display_name"`
	Port        types.Number `tfsdk:"port"`
	Protocol    types.String `tfsdk:"protocol"`
	TargetPool  types.String `tfsdk:"target_pool"`
}

type Network struct {
	NetworkID types.String `tfsdk:"network_id"`
	Role      types.String `tfsdk:"role"`
}

type TargetPool struct {
	Name         types.String `tfsdk:"name"`
	TargetPort   types.Number `tfsdk:"target_port"`
	Targets      types.Set    `tfsdk:"targets"`
	HealthChecks types.Set    `tfsdk:"health_check"`
}

type Target struct {
	DisplayName types.Number `tfsdk:"display_name"`
	IPAddress   types.String `tfsdk:"ip_address"`
}

type HealthCheck struct {
	HealthyThreshold   types.Number `tfsdk:"healthy_threshold"`
	Interval           types.String `tfsdk:"interval"`
	IntervalJitter     types.String `tfsdk:"interval_jitter"`
	Timeout            types.String `tfsdk:"timeout"`
	UnhealthyThreshold types.Number `tfsdk:"unhealthy_threshold"`
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
			"external_ip": schema.StringAttribute{
				Description: "The external IP address of the instance.",
				Optional:    true,
			},
			"listeners": schema.SetNestedAttribute{
				Description: "The load balancers listeners.",
				Optional:    true,
				// Required:    true,
				// Validators: []validator.Set{
				// 	setvalidator.SizeAtLeast(1),
				// },
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.NumberAttribute{
							Description: "The port the load balancer listens on.",
							Required:    true,
						},
						"port": schema.NumberAttribute{
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
				Optional:    true,
				// Required:    true,
				// Validators: []validator.Set{
				// 	setvalidator.SizeAtLeast(1),
				// },
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_id": schema.StringAttribute{
							Description: "The network UUID.",
							Required:    true,
						},
						"role": schema.StringAttribute{
							Description: "The network role.",
							Required:    true,
						},
					},
				},
			},
			"target_pools": schema.SetNestedAttribute{
				Description: "The load balancers target pools.",
				Optional:    true,
				// Required:    true,
				// Validators: []validator.Set{
				// 	setvalidator.SizeAtLeast(1),
				// },
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The target pool name.",
							Required:    true,
						},
						"target_port": schema.NumberAttribute{
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
									"display_name": schema.NumberAttribute{
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
						"health_check": schema.SetNestedAttribute{
							Description: "The target pool health checks.",
							Required:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"healthy_threshold": schema.NumberAttribute{
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
									"unhealthy_threshold": schema.NumberAttribute{
										Description: "The unhealthy threshold.",
										Required:    true,
									},
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
			},
			"private_network_only": schema.BoolAttribute{
				Description: "Whether the load balancer is only accessible via private networks.",
				Optional:    true,
			},
		},
	}
}
