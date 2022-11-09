package instance

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Instance is the schema model
type Instance struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ProjectID types.String `tfsdk:"project_id"`
	Plan      types.String `tfsdk:"plan"`
	PlanID    types.String `tfsdk:"plan_id"`
	Version   types.String `tfsdk:"version"`
	ACL       types.List   `tfsdk:"acl"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `Manages ElasticSearch instances`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Required:    true,
			},
			"project_id": {
				Description: "The project ID the cluster runs in. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"plan": {
				Description: "The ElasticSearch Plan. Options are: `stackit-elasticsearch-single-small`, `stackit-elasticsearch-cluster-small`, `stackit-elasticsearch-single-medium`, `stackit-elasticsearch-cluster-medium`, `stackit-elasticsearch-cluster-big`",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.StringDefault(default_plan),
				},
			},
			"plan_id": {
				Description: "The selected plan ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"version": {
				Description: "Elasticsearch version. Options: `5`, `6`, `7`. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
					modifiers.StringDefault(default_version),
				},
			},
			"acl": {
				Description: "Access Control rules to whitelist IP addresses",
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
			},
		},
	}, nil
}
