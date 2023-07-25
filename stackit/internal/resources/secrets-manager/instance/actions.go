package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/secrets-manager/v1.1.0/acls"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/secrets-manager/v1.1.0/instances"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/utils/strings/slices"
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
	uuidProjectID := uuid.MustParse(plan.ProjectID.ValueString())
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

	if !plan.ACL.IsUnknown() {
		r.manageACLs(ctx, &plan, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	r.readACLs(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r Resource) manageACLs(ctx context.Context, plan *Instance, diags *diag.Diagnostics) {
	want, idsToRemove := []string{}, []string{}
	diags.Append(plan.ACL.ElementsAs(ctx, &want, true)...)
	if diags.HasError() {
		return
	}
	r.readACLs(ctx, plan, diags)
	if diags.HasError() {
		return
	}
	c := r.client
	res, err := c.SecretsManager.Acls.List(ctx, uuid.MustParse(plan.ProjectID.ValueString()), uuid.MustParse(plan.ID.ValueString()))
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		diags.AddError("failed to get instance ACLs", agg.Error())
		return
	}
	for _, el := range res.JSON200.Acls {
		if !slices.Contains(want, el.Cidr) {
			idsToRemove = append(idsToRemove, el.ID)
			continue
		}
		// remove from want
		if i := slices.Index(want, el.Cidr); i > -1 {
			want = append(want[:i], want[i+1:]...)
		}
	}
	// remove
	for _, id := range idsToRemove {
		_, err := c.SecretsManager.Acls.Delete(ctx, uuid.MustParse(plan.ProjectID.ValueString()), uuid.MustParse(plan.ID.ValueString()), uuid.MustParse(id))
		if err != nil {
			diags.AddError("failed to delete instance ACL", err.Error())
			return
		}
	}
	// add
	for _, cidr := range want {
		_, err := c.SecretsManager.Acls.Create(ctx, uuid.MustParse(plan.ProjectID.ValueString()), uuid.MustParse(plan.ID.ValueString()), acls.AclCreate{
			Cidr: cidr,
		})
		if err != nil {
			diags.AddError("failed to create instance ACL", err.Error())
			return
		}
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

	r.readACLs(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) readACLs(ctx context.Context, config *Instance, diags *diag.Diagnostics) {
	c := r.client
	res, err := c.SecretsManager.Acls.List(ctx, uuid.MustParse(config.ProjectID.ValueString()), uuid.MustParse(config.ID.ValueString()))
	if agg := validate.Response(res, err, "JSON200"); agg != nil {
		diags.AddError("failed to get instance ACLs", agg.Error())
		return
	}
	els := []attr.Value{}
	for _, el := range res.JSON200.Acls {
		els = append(els, types.StringValue(el.Cidr))
	}
	config.ACL = types.SetValueMust(types.StringType, els)
	return
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Instance
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state Instance
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ACL.Equal(state.ACL) {
		return
	}

	r.manageACLs(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readACLs(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
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
