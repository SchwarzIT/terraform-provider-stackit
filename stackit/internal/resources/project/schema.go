package project

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gorm.io/gorm/schema"
)

// Project is the schema model
type Project struct {
	ID                types.String `tfsdk:"id"`
	ContainerID       types.String `tfsdk:"container_id"`
	ParentContainerID types.String `tfsdk:"parent_container_id"`
	Name              types.String `tfsdk:"name"`
	BillingRef        types.String `tfsdk:"billing_ref"`
	OwnerEmail        types.String `tfsdk:"owner_email"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages projects",
		Attributes: map[string]schema.Attribute{
			"id": {
				Description: "the project ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"container_id": {
				Description: "the project container ID",
				Type:        types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"parent_container_id": {
				Description: "the container ID in which the project will be created",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: planmodifier.Strings{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": {
				Description: "the project name",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.ProjectName(),
				},
			},

			"billing_ref": {
				Description: "billing reference for cost transparency",
				Type:        types.StringType,
				Required:    true,
				Validators: []validator.String{
					validate.BillingRef(),
				},
			},

			"owner_email": {
				Description: "Email address of owner of the project. This value is only considered during creation. changing it afterwards will have no effect.",
				Type:        types.StringType,
				Required:    true,
			},
		},
	}
}
