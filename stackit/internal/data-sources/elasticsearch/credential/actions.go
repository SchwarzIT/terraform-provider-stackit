package credential

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Read - lifecycle function
func (r DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	c := r.client
	var config Credential
	diags := req.Config.Get(ctx, &config)
	es := c.DataServices.ElasticSearch

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := es.Credentials.List(ctx, config.ProjectID.Value, config.InstanceID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to list credentials", err.Error())
		return
	}

	found := -1
	existing := ""
	for i, instance := range list.CredentialsList {
		if instance.ID == config.ID.Value {
			found = i
			break
		}
		if existing == "" {
			existing = "\navailable credentials in the project are:"
		}
		existing = fmt.Sprintf("%s\n- %s", existing, instance.ID)
	}

	if found == -1 {
		resp.State.RemoveResource(ctx)
		diags.AddError("couldn't find credentials", "credentials could not be found."+existing)
		resp.Diagnostics.Append(diags...)
		return
	}

	// set found instance
	intance := list.CredentialsList[found]

	res, err := es.Options.GetOfferings(ctx, config.ProjectID.Value)
	if err != nil {
		resp.Diagnostics.AddError("failed to get offerings", err.Error())
		return
	}

	for _, offer := range res.Offerings {
		for _, p := range offer.Plans {
			if p.ID != intance.ID {
				continue
			}
			// config.Plan = types.String{Value: p.Name}
			// config.Version = types.String{Value: offer.Version}
		}
	}

	// config.ID = types.String{Value: intance.InstanceID}
	// config.PlanID = types.String{Value: intance.PlanID}
	// config.DashboardURL = types.String{Value: intance.DashboardURL}
	// config.CFGUID = types.String{Value: intance.CFGUID}
	// config.CFSpaceGUID = types.String{Value: intance.CFSpaceGUID}
	// config.CFOrganizationGUID = types.String{Value: intance.CFOrganizationGUID}
	//
	// config.ACL = types.List{ElemType: types.StringType}
	// if aclString, ok := intance.Parameters["sgw_acl"]; ok {
	// 	items := strings.Split(aclString, ",")
	// 	for _, v := range items {
	// 		config.ACL.Elems = append(config.ACL.Elems, types.String{Value: v})
	// 	}
	// } else {
	// 	config.ACL.Null = true
	// }

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
