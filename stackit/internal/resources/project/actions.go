package project

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/projects"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Create - lifecycle function
func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.Provider.IsConfigured() {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on another resource.",
		)
		return
	}

	var plan Project
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// handle creation
	plan = r.createProject(ctx, resp, plan)
	if resp.Diagnostics.HasError() {
		return
	}

	p := Project{
		ID:                  types.String{Value: plan.ID.Value},
		Name:                types.String{Value: plan.Name.Value},
		BillingRef:          types.String{Value: plan.BillingRef.Value},
		Owner:               types.String{Value: plan.Owner.Value},
		EnableKubernetes:    types.Bool{Null: true},
		EnableObjectStorage: types.Bool{Null: true},
	}

	if !plan.EnableKubernetes.IsNull() {
		p.EnableKubernetes = types.Bool{Value: plan.EnableKubernetes.Value}
		r.createKubernetesProject(ctx, &resp.Diagnostics, plan.ID.Value)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.EnableObjectStorage.IsNull() {
		p.EnableObjectStorage = types.Bool{Value: plan.EnableObjectStorage.Value}
		r.createObjectStorageProject(ctx, &resp.Diagnostics, plan.ID.Value)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// update state
	diags = resp.State.Set(ctx, p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) createProject(ctx context.Context, resp *resource.CreateResponse, plan Project) Project {
	owners := projects.ProjectRole{
		Name: consts.ROLE_PROJECT_OWNER,
		Users: []projects.ProjectRoleMember{
			{ID: plan.Owner.Value},
		},
		ServiceAccounts: []projects.ProjectRoleMember{ // service account is added automatically
			{ID: r.Provider.ServiceAccountID()},
		},
	}

	c := r.Provider.Client()
	project, process, err := c.Projects.Create(ctx, plan.Name.Value, plan.BillingRef.Value, owners)
	if err != nil {
		resp.Diagnostics.AddError("failed to create project", err.Error())
		return plan
	}

	if _, err := process.Wait(); err != nil {
		resp.Diagnostics.AddError("failed to verify project is active", err.Error())
		return plan
	}

	plan.ID = types.String{Value: project.ID}
	return plan
}

func (r Resource) createKubernetesProject(ctx context.Context, d *diag.Diagnostics, projectID string) {
	c := r.Provider.Client()
	_, process, err := c.Kubernetes.Projects.Create(ctx, projectID)
	if err != nil {
		d.AddError("failed to verify kubernetes is enabled for project", err.Error())
		return
	}

	if _, err := process.Wait(); err != nil {
		d.AddError("failed to validate kubernetes is enabled for project", err.Error())
		return
	}
}

func (r Resource) createObjectStorageProject(ctx context.Context, d *diag.Diagnostics, projectID string) {
	c := r.Provider.Client()
	created := false

	if err := helper.RetryContext(ctx, common.DURATION_10M, func() *helper.RetryError {
		if !created {
			if _, err := c.ObjectStorage.Projects.Create(ctx, projectID); err != nil {
				return helper.RetryableError(err)
			}
			created = true
		}

		// get object storage project
		_, err := c.ObjectStorage.Projects.Get(ctx, projectID)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		d.AddError("failed to verify object storage is enabled for project", err.Error())
		return
	}
}

// Read - lifecycle function
func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	c := r.Provider.Client()
	var project projects.Project
	var p Project

	diags := req.State.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.RetryContext(ctx, common.DURATION_1M, func() *helper.RetryError {
		var err error
		project, err = c.Projects.Get(ctx, p.ID.Value)
		if err != nil {
			return helper.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("failed to read project", err.Error())
		return
	}

	if !p.EnableKubernetes.IsNull() {
		kubernetesEnabled := false
		if res, err := c.Kubernetes.Projects.Get(ctx, p.ID.Value); err == nil && res.State == consts.SKE_PROJECT_STATUS_CREATED {
			kubernetesEnabled = true
		}
		p.EnableKubernetes = types.Bool{Value: kubernetesEnabled}
	}

	if !p.EnableObjectStorage.IsNull() {
		obejctStorageEnabled := false
		if _, err := c.ObjectStorage.Projects.Get(ctx, p.ID.Value); err == nil {
			obejctStorageEnabled = true
		}
		p.EnableObjectStorage = types.Bool{Value: obejctStorageEnabled}
	}

	p.ID = types.String{Value: project.ID}
	p.Name = types.String{Value: project.Name}
	p.BillingRef = types.String{Value: project.BillingReference}

	diags = resp.State.Set(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - lifecycle function
func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Project
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state Project
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsUnknown() {
		plan.ID = state.ID
	}

	r.updateProject(ctx, plan, state, resp)
	r.updateKubernetesProject(ctx, plan, state, resp)
	r.updateObjectStorageProject(ctx, plan, state, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	// update state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r Resource) updateProject(ctx context.Context, plan, state Project, resp *resource.UpdateResponse) {
	if plan.Name.Equal(state.Name) && plan.BillingRef.Equal(state.BillingRef) {
		return
	}
	c := r.Provider.Client()
	err := c.Projects.Update(ctx, plan.ID.Value, plan.Name.Value, plan.BillingRef.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to update project", err.Error())
		return
	}
}

func (r Resource) updateKubernetesProject(ctx context.Context, plan, state Project, resp *resource.UpdateResponse) {
	if plan.EnableKubernetes.Equal(state.EnableKubernetes) {
		return
	}

	if plan.EnableKubernetes.IsNull() || !plan.EnableKubernetes.Value {
		r.deleteKubernetesProject(ctx, &resp.Diagnostics, plan.ID.Value)
		return
	}

	r.createKubernetesProject(ctx, &resp.Diagnostics, plan.ID.Value)
}

func (r Resource) updateObjectStorageProject(ctx context.Context, plan, state Project, resp *resource.UpdateResponse) {
	if plan.EnableObjectStorage.Equal(state.EnableObjectStorage) {
		return
	}

	if plan.EnableObjectStorage.IsNull() || !plan.EnableObjectStorage.Value {
		r.deleteObjectStorageProject(ctx, &resp.Diagnostics, plan.ID.Value)
		return
	}

	r.createObjectStorageProject(ctx, &resp.Diagnostics, plan.ID.Value)
}

// Delete - lifecycle function
func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Project
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := r.Provider.Client()

	if state.EnableKubernetes.Value {
		_, _ = c.Kubernetes.Projects.Delete(ctx, state.ID.Value)
	}

	if state.EnableObjectStorage.Value {
		_ = c.ObjectStorage.Projects.Delete(ctx, state.ID.Value)
	}

	process, err := c.Projects.Delete(ctx, state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete project", err.Error())
		return
	}

	if _, err := process.Wait(); err != nil {
		resp.Diagnostics.AddError("failed to verify project deletion", err.Error())
	}

	resp.State.RemoveResource(ctx)
}

func (r Resource) deleteKubernetesProject(ctx context.Context, d *diag.Diagnostics, projectID string) {
	c := r.Provider.Client()
	res, err := c.Kubernetes.Clusters.List(ctx, projectID)
	if err == nil && len(res.Items) > 0 {
		d.AddWarning("Kubernetes disabling considerations", `We detected active Kubernetes clusters in your project
Therefore, in order to prevent them from automatically being deleted, we ignored your setting "enable_kubernetes=false".
If you wish for the change to be applied, please delete all existing clusters first & re-run the plan.`)
		return
	}

	process, err := c.Kubernetes.Projects.Delete(ctx, projectID)
	if err != nil {
		d.AddError("error disabling kubernetes", err.Error())
		return
	}

	if _, err := process.Wait(); err != nil {
		d.AddError("kubernetes disabling validation failed", err.Error())
		return
	}
}

func (r Resource) deleteObjectStorageProject(ctx context.Context, d *diag.Diagnostics, projectID string) {
	c := r.Provider.Client()
	res, err := c.ObjectStorage.Buckets.List(ctx, projectID)
	if err == nil && len(res.Buckets) > 0 {
		d.AddWarning("Object Storage disabling considerations", `We detected active buckets in your project
Therefor, in order to prevent them from automatically being deleted, we ignored your setting "enable_object_storage=false".
If you wish for the change to be applied, please delete all existing buckets first & re-run the plan.`)
		return
	}
	if err := c.ObjectStorage.Projects.Delete(ctx, projectID); err != nil {
		if strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
			return
		}
		d.AddError("failed to disable object storage", err.Error())
		return
	}
}

// ImportState handles terraform import
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// validate project id
	if err := clientValidate.ProjectID(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Couldn't validate project_id.\n%s", err.Error()),
		)
		return
	}

	// set main attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

}
