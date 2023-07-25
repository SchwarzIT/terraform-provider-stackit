package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/secrets-manager/v1.1.0/instances"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	uuidProjectID, err := uuid.Parse(plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Project ID",
			fmt.Sprintf("Couldn't parse project_id\n%s", err.Error()),
		)
		return
	}

	res, err := c.SecretsManager.Instances.Create(ctx, uuidProjectID, instances.CreateJSONRequestBody{
		Name: plan.Name.ValueString(),
	})
	if agg := validate.Response(res, err, "JSON201"); agg != nil {
		if res == nil || res.StatusCode() != http.StatusOK {
			if res != nil {
				common.Dump(&resp.Diagnostics, res.Body)
			}
			resp.Diagnostics.AddError("failed to create instance", agg.Error())
			return
		}
		// handle wrong status code response from API
		res.JSON201 = &instances.Instance{}
		if err := json.Unmarshal(res.Body, res.JSON201); err != nil {
			resp.Diagnostics.AddError("failed to parse response", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(res.JSON201.ID)
	plan.Frontend = types.StringValue(res.JSON201.ApiUrl + "/ui")
	plan.API = types.StringValue(res.JSON201.ApiUrl)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.client
	var state Instance

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := c.SecretsManager.Instances.Get(ctx, uuid.MustParse(state.ProjectID.ValueString()), uuid.MustParse(state.ID.ValueString()))
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		if validate.StatusEquals(res, http.StatusNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to get instance", agg.Error())
		return
	}

	state.Name = types.StringValue(res.JSON200.Name)
	state.Frontend = types.StringValue(res.JSON200.ApiUrl + "/ui")
	state.API = types.StringValue(res.JSON200.ApiUrl)

	// update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Instance
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.client
	_, err := c.SecretsManager.Instances.Delete(ctx, uuid.MustParse(state.ProjectID.ValueString()), uuid.MustParse(state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("failed to delete instance", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `project_id,id` where `id` is the instance ID.\nInstead got: %q", req.ID),
		)
		return
	}

	// validate project id
	if err := clientValidate.ProjectID(idParts[0]); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate kubernetes_project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}

}
