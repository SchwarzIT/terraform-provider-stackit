package credential

import (
	"context"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/data-services/credentials"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r Resource) validate(ctx context.Context, data *Credential) error {

	return nil
}

func (r Resource) applyClientResponse(ctx context.Context, c *Credential, cgr credentials.GetResponse) error {

	c.ID = types.String{Value: cgr.ID}
	c.URI = types.String{Value: cgr.URI}

	// c.ProjectID = types.String{Value: cgr.}
	// c.InstanceID = types.String{Value: cgr.}

	c.SyslogDrainURL = types.String{Value: cgr.Raw.SyslogDrainURL}
	c.RouteServiceURL = types.String{Value: cgr.Raw.RouteServiceURL}

	c.VolumeMounts = types.List{ElemType: types.MapType{
		ElemType: types.StringType,
	}}

	// if len(cgr.Raw.VolumeMounts) > 0 {
	// 	c.VolumeMounts.Elems = make([]types.MapType, len(cgr.Raw.VolumeMounts))
	// }
	// for k, v := range cgr.Raw.VolumeMounts {
	//
	// 	c.VolumeMounts.Elems[k] = types.Map{ElemType: types.StringType}
	//
	// 	if len(cgr.Raw.VolumeMounts[k]) > 0 {
	// 		c.VolumeMounts.Elems[k] = make(map[string]attr.Value, len(cgr.Raw.VolumeMounts[k]))
	// 	}
	// 	for k, v := range cgr.Raw.VolumeMounts {
	//
	// 		c.VolumeMounts.Elems[k] = types.String{Value: v}
	// 	}
	// }

	c.Host = types.String{Value: cgr.Raw.Credential.Host}
	c.Port = types.Int64{Value: int64(cgr.Raw.Credential.Port)}

	c.Hosts = types.List{ElemType: types.StringType}
	if len(cgr.Raw.Credential.Hosts) > 0 {
		c.Hosts.Elems = make([]attr.Value, len(cgr.Raw.Credential.Hosts))
	}
	for k, v := range cgr.Raw.Credential.Hosts {
		c.Hosts.Elems[k] = types.String{Value: v}
	}

	c.Username = types.String{Value: cgr.Raw.Credential.Name}
	c.Password = types.String{Value: cgr.Raw.Credential.Password}

	c.Protocols = types.Map{ElemType: types.StringType}
	if len(cgr.Raw.Credential.Protocols) > 0 {
		c.Protocols.Elems = make(map[string]attr.Value, len(cgr.Raw.Credential.Protocols))
	}
	for k, v := range cgr.Raw.Credential.Protocols {
		c.Protocols.Elems[k] = types.String{Value: v}
	}
	return nil
}
