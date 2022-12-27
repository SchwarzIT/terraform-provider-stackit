package instance

import (
	"context"
	"fmt"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gorm.io/gorm/schema"
)

// Instance is the schema model
type Instance struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProjectID          types.String `tfsdk:"project_id"`
	Plan               types.String `tfsdk:"plan"`
	PlanID             types.String `tfsdk:"plan_id"`
	Version            types.String `tfsdk:"version"`
	ACL                types.List   `tfsdk:"acl"`
	DashboardURL       types.String `tfsdk:"dashboard_url"`
	CFGUID             types.String `tfsdk:"cf_guid"`
	CFSpaceGUID        types.String `tfsdk:"cf_space_guid"`
	CFOrganizationGUID types.String `tfsdk:"cf_organization_guid"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages RabbitMQ instances

~> **Note:** RabbitMQ API (Part of DSA APIs) currently has issues reflecting updates & configuration correctly. Therefore, this resource is not ready for production usage.
		`,
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": {
				Description: "Specifies the instance name. Changing this value requires the resource to be recreated. Changing this value requires the resource to be recreated.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": {
				Description: "The project ID.",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan": {
				Description: fmt.Sprintf("The RabbitMQ Plan. Default is `%s`.\nOptions are: `stackit-messaging-cluster-big`, `stackit-messaging-cluster-medium`, `stackit-messaging-cluster-small`, `stackit-messaging-single-medium`, `stackit-messaging-single-small`, `stackit-rabbitmq-cluster-medium`, `stackit-rabbitmq-single-medium`, `stackit-rabbitmq-cluster-big`,`stackit-rabbitmq-cluster-small`, `stackit-rabbitmq-single-small`", default_plan),
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					modifiers.StringDefault(default_plan),
				},
			},
			"plan_id": {
				Description: "The selected plan ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"version": {
				Description: fmt.Sprintf("RabbitMQ version. Default is %s.\nOptions: `3.10`, `3.8`, `3.7`. Changing this value requires the resource to be recreated.", default_version),
				Type:        types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					modifiers.StringDefault(default_version),
				},
			},
			"acl": {
				Description: "Access Control rules to whitelist IP addresses",
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dashboard_url": {
				Description: "Dashboard URL",
				Type:        types.StringType,
				Computed:    true,
			},
			"cf_guid": {
				Description: "Cloud Foundry GUID",
				Type:        types.StringType,
				Computed:    true,
			},
			"cf_space_guid": {
				Description: "Cloud Foundry Space GUID",
				Type:        types.StringType,
				Computed:    true,
			},
			"cf_organization_guid": {
				Description: "Cloud Foundry Organization GUID",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}
}
